package scan

import (
	"fmt"
	"github.com/owenrumney/squealer/internal/app/squealer/config"
	"github.com/owenrumney/squealer/internal/app/squealer/mertics"
	"os"
	"strings"
)

type ScannerConfig struct {
	cfg      *config.Config
	basepath string
	redacted bool
	fromRef  string
	noGit    bool
}

func NewScannerConfig(basepath string, redacted, noGit bool, cfg *config.Config) ScannerConfig {
	return ScannerConfig{
		cfg:      cfg,
		basepath: basepath,
		redacted: redacted,
		noGit:    noGit,
	}
}

type Scanner interface {
	Scan() error
	GetMetrics() *mertics.Metrics
}

func NewScanner(sc ScannerConfig) (Scanner, error) {
	if sc.noGit || notGit(sc.basepath) {
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
