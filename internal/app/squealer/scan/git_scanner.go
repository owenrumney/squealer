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
	"sync"
)

type gitScanner struct {
	mc               match.MatcherController
	metrics          *mertics.Metrics
	workingDirectory string
	ignorePaths      []string
	fromHash         plumbing.Hash
	toHash           plumbing.Hash
	ignoreExtensions []string
	headSet          bool
}

func (s *gitScanner) GetType() ScannerType {
	return GitScanner
}

func newGitScanner(sc ScannerConfig) (*gitScanner, error) {
	if _, err := os.Stat(sc.Basepath); err != nil {
		return nil, err
	}
	metrics := mertics.NewMetrics()
	mc := match.NewMatcherController(sc.Cfg, metrics, sc.Redacted)

	scanner := &gitScanner{
		mc:               *mc,
		metrics:          metrics,
		workingDirectory: sc.Basepath,
		ignorePaths:      sc.Cfg.IgnorePrefixes,
		ignoreExtensions: sc.Cfg.IgnoreExtensions,
	}
	if len(sc.FromHash) > 0 {
		scanner.fromHash = plumbing.NewHash(sc.FromHash)
	}

	if len(sc.ToHash) > 0 {
		scanner.toHash = plumbing.NewHash(sc.ToHash)
		scanner.headSet = true
	}

	return scanner, nil
}

func (s *gitScanner) Scan() error {
	client, err := git.PlainOpen(s.workingDirectory)
	if err != nil {
		return err
	}

	commits, err := s.getRelevantCommitIter(client)
	if err != nil {
		return err
	}

	s.metrics.StartTimer()
	defer s.metrics.StopTimer()
	commit, err := commits.Next()
	for err == nil && commit != nil {
		func(c *object.Commit) {
			if err := s.processCommit(c); err != nil {
				fmt.Println(err.Error())
			}
		}(commit)
		if commit.Hash.String() == s.fromHash.String() {
			fmt.Println("commit hash reached - stopping")
			// reached the starting commit - stop here
			return nil
		}
		commit, err = commits.Next()
	}

	if err != nil && err != io.EOF {
		fmt.Printf("error is not null %s\n", err.Error())
	}
	return nil
}

func (s *gitScanner) getRelevantCommitIter(client *git.Repository) (object.CommitIter, error) {
	var headRef plumbing.Hash
	if s.headSet {
		headRef = s.toHash

	} else {
		ref, _ := client.Head()
		if ref != nil {
			headRef = ref.Hash()
		}
		headRef = plumbing.ZeroHash
	}

	var commits object.CommitIter
	var err error

	if headRef != plumbing.ZeroHash {
		commits, err = client.Log(&git.LogOptions{
			From:  headRef,
			Order: git.LogOrderCommitterTime,
		})
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

func (s *gitScanner) processCommit(commit *object.Commit) error {

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

	s.metrics.IncrementCommitsProcessed()
	close(ch)
	wg.Wait()
	return nil
}

func (s *gitScanner) processFile(file *object.File) error {

	if isBin, err := file.IsBinary(); err != nil || isBin {
		return nil
	}
	if shouldIgnore(file.Name, s.ignorePaths, s.ignoreExtensions) {
		return nil
	}
	content, err := file.Contents()
	if err != nil {
		return err
	}

	err = s.mc.Evaluate(file.Name, content)
	s.metrics.IncrementFilesProcessed()
	return err
}

func (s *gitScanner) GetMetrics() *mertics.Metrics {
	return s.metrics
}
