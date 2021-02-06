package scan

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"io"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/owenrumney/squealer/internal/app/squealer/match"
	"github.com/owenrumney/squealer/internal/app/squealer/mertics"
)

type GitScanner struct {
	mc               match.MatcherController
	metrics          *mertics.Metrics
	workingDirectory string
	ignorePaths      []string
}

func NewGitScanner(sc ScannerConfig) (*GitScanner, error) {
	if _, err := os.Stat(sc.basepath); err != nil {
		return nil, err
	}
	metrics := mertics.NewMetrics()
	mc := match.NewMatcherController(sc.cfg, metrics, sc.redacted)

	return &GitScanner{
		mc:               *mc,
		metrics:          metrics,
		workingDirectory: sc.basepath,
		ignorePaths: []string{
			"vendor",
			"internal/swagger-models",
			"internal/swagger-clients",
		},
	}, nil
}

func (s *GitScanner) Scan() error {
	client, err := git.PlainOpen(s.workingDirectory)
	if err != nil {
		return err
	}

	var commits object.CommitIter
	ref, _ := client.Head()
	if ref != nil {
		commits, err = client.Log(&git.LogOptions{From: ref.Hash()})
		if err != nil {
			return err
		}
	} else {
		commits, err = client.CommitObjects()
		if err != nil {
			return err
		}
	}

	start := time.Now()
	commit, err := commits.Next()
	for err == nil && commit != nil {
		func(c *object.Commit) {
			if err := s.processCommit(c); err != nil {
				fmt.Println(err.Error())
			}
		}(commit)
		commit, err = commits.Next()
	}
	stop := time.Now()

	if err != nil && err != io.EOF {
		fmt.Printf("error is not null %s\n", err.Error())
	}

	fmt.Printf("Process took %v\n", stop.Sub(start).Seconds())
	return nil
}

func (s *GitScanner) processCommit(commit *object.Commit) error {
	files, err := commit.Files()
	if err != nil {
		return err
	}

	var ch = make(chan *object.File, 50)
	var wg sync.WaitGroup

	processes := runtime.NumCPU()/2 - 1
	wg.Add(processes)
	for i := 0; i < processes; i++ {
		go func() {
			for {
				f, ok := <-ch
				if !ok {
					wg.Done()
					return
				}
				err := s.processFile(f)
				if err != nil {
					fmt.Println(err.Error())
				}
			}
		}()
	}

	file, err := files.Next()
	for err == nil && file != nil {
		ch <- file
		file, err = files.Next()
	}

	close(ch)
	wg.Wait()

	s.metrics.UpdateCommitsProcessed()
	return nil
}

func (s *GitScanner) processFile(file *object.File) error {
	s.metrics.UpdateFilesProcessed()
	if isBin, err := file.IsBinary(); err != nil || isBin {
		return nil
	}

	for _, ignorePath := range s.ignorePaths {
		if strings.HasPrefix(file.Name, ignorePath) {
			return nil
		}
	}

	if strings.HasSuffix(file.Name, ".zip") {
		return nil
	}

	if err := s.mc.Evaluate(file); err == nil {
		return err
	}
	return nil
}

func (s *GitScanner) Exit(printMetrics bool) {
	if !printMetrics {
		fmt.Printf(`
Processing:
  commits:      %d
  commit files: %d

Transgressions:
  identified:   %d
  ignored:      %d
  reported:     %d

`,
			s.metrics.CommitsProcessed,
			s.metrics.FilesProcessed,
			s.metrics.TransgressionsFound,
			s.metrics.TransgressionsIgnored,
			s.metrics.TransgressionsReported)
	}
	if s.metrics.TransgressionsReported > 0 {
		os.Exit(1)
	}
	os.Exit(0)
}