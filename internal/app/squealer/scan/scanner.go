package scan

import (
	"fmt"
	"github.com/owenrumney/squealer/internal/app/squealer/config"
	"github.com/owenrumney/squealer/internal/app/squealer/mertics"
	"os"
	"strings"
)

type ScannerType string

const (
	GitScanner       ScannerType = "gitScanner"
	DirectoryScanner ScannerType = "dirScanner"
)

type ScannerConfig struct {
	Cfg      *config.Config
	Basepath string
	Redacted bool
	NoGit    bool
	FromHash string
	ToHash   string
}

type Scanner interface {
	Scan() error
	GetMetrics() *mertics.Metrics
	GetType() ScannerType
}

func NewScanner(sc ScannerConfig) (Scanner, error) {
	if sc.NoGit || notGit(sc.Basepath) {
		return newDirectoryScanner(sc)
	}
	return newGitScanner(sc)
}

func notGit(basepath string) bool {
	if stat, err := os.Stat(basepath); err == nil && stat != nil {
		if _, ok := os.Stat(fmt.Sprintf("%s/.git", basepath)); ok == nil {
			return false
		}
	}
	return true
}

func shouldIgnore(filename string, ignorePaths []string, ignoreExtensions []string) bool {
	for _, ignorePath := range ignorePaths {
		if strings.HasPrefix(filename, ignorePath) {
			return true
		}
	}
	for _, ext := range ignoreExtensions {
		if strings.HasSuffix(filename, ext) {
			return true
		}
	}
	return false
}
