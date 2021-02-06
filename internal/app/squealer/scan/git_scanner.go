package scan

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/owenrumney/squealer/internal/app/squealer/match"
	"github.com/owenrumney/squealer/internal/app/squealer/mertics"
	"io"
	"math"
	"os"
	"runtime"
	"strings"
	"sync"
)

type GitScanner struct {
	mc               match.MatcherController
	metrics          *mertics.Metrics
	workingDirectory string
	ignorePaths      []string
	fromRef          plumbing.Hash
	ignoreExtensions []string
}

func NewGitScanner(sc ScannerConfig) (*GitScanner, error) {
	if _, err := os.Stat(sc.basepath); err != nil {
		return nil, err
	}
	metrics := mertics.NewMetrics()
	mc := match.NewMatcherController(sc.cfg, metrics, sc.redacted)

	scanner := &GitScanner{
		mc:               *mc,
		metrics:          metrics,
		workingDirectory: sc.basepath,
		ignorePaths:      sc.cfg.IgnorePrefixes,
		ignoreExtensions: sc.cfg.IgnoreExtensions,
	}
	if len(sc.fromRef) > 0 {
		scanner.fromRef = plumbing.NewHash(sc.fromRef)
	}

	return scanner, nil
}

func (s *GitScanner) Scan() error {
	client, err := git.PlainOpen(s.workingDirectory)
	if err != nil {
		return err
	}

	commits, err := s.getRelevantCommitIter(client)
	if err != nil {
		return err
	}

	s.metrics.StartTimer()
	commit, err := commits.Next()
	for err == nil && commit != nil {
		func(c *object.Commit) {
			if err := s.processCommit(c); err != nil {
				fmt.Println(err.Error())
			}
		}(commit)
		commit, err = commits.Next()
	}
	s.metrics.StopTimer()

	if err != nil && err != io.EOF {
		fmt.Printf("error is not null %s\n", err.Error())
	}
	return nil
}

func (s *GitScanner) getRelevantCommitIter(client *git.Repository) (object.CommitIter, error) {
	if s.fromRef == plumbing.ZeroHash {
		headRef, _ := client.Head()
		if headRef != nil {
			s.fromRef = headRef.Hash()
		}
	}

	var commits object.CommitIter
	var err error

	if s.fromRef != plumbing.ZeroHash {
		commits, err = client.Log(&git.LogOptions{From: s.fromRef})
		if err != nil {
			return nil, err
		}
	} else {
		commits, err = client.CommitObjects()
		if err != nil {
			return nil, err
		}
	}
	return commits, err
}

func (s *GitScanner) processCommit(commit *object.Commit) error {
	files, err := commit.Files()
	if err != nil {
		return err
	}

	var ch = make(chan *object.File, 50)
	var wg sync.WaitGroup

	processes := int(math.Max(float64(runtime.NumCPU()/2-1), 1))
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

	s.metrics.IncrementCommitsProcessed()
	return nil
}

func (s *GitScanner) processFile(file *object.File) error {
	s.metrics.IncrementFilesProcessed()
	if isBin, err := file.IsBinary(); err != nil || isBin {
		return nil
	}

	for _, ignorePath := range s.ignorePaths {
		if strings.HasPrefix(file.Name, ignorePath) {
			return nil
		}
	}

	for _, ext := range s.ignoreExtensions {
		if strings.HasSuffix(file.Name, ext) {
			return nil
		}
	}
	return s.mc.Evaluate(file)
}

func (s *GitScanner) Shutdown(printMetrics bool) {
	if !printMetrics {
		duration, _ := s.metrics.Duration()
		fmt.Printf(`
Processing:
  duration:     %4.2fs
  commits:      %d
  commit files: %d

transgressionMap:
  identified:   %d
  ignored:      %d
  reported:     %d

`,
			duration,
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

func (s *GitScanner) GetMetrics() *mertics.Metrics {
	return s.metrics
}
