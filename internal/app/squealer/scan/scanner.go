package scan

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"regexp"
	"strings"

	"github.com/owenrumney/squealer/internal/app/squealer/config"
	"github.com/owenrumney/squealer/internal/app/squealer/match"
	"github.com/owenrumney/squealer/internal/app/squealer/mertics"
)

type ScannerType string

const (
	GitScanner       ScannerType = "gitScanner"
	DirectoryScanner ScannerType = "dirScanner"
)

type ScannerConfig struct {
	Cfg        *config.Config
	Basepath   string
	Redacted   bool
	NoGit      bool
	FromHash   string
	ToHash     string
	Everything bool
}

type Scanner interface {
	Scan() ([]match.Transgression, error)
	GetMetrics() *mertics.Metrics
	GetType() ScannerType
}

func NewScanner(sc ScannerConfig) (Scanner, error) {
	if sc.NoGit || notGit(sc.Basepath) {
		log.Infof("Using a directory scanner to process %s\n", sc.Basepath)
		return newDirectoryScanner(sc)
	}
	log.Infof("Using a git scanner to process %s\n", sc.Basepath)
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
		if match, err := regexp.MatchString(fmt.Sprintf(`\b%s\b`, ignorePath), filename); err == nil && match {
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
