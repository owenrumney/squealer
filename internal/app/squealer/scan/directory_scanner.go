package scan

import (
	"github.com/owenrumney/squealer/internal/app/squealer/match"
	"github.com/owenrumney/squealer/internal/app/squealer/mertics"
	"io/ioutil"
	"os"
	"path/filepath"
)

type directoryScanner struct {
	mc               match.MatcherController
	metrics          *mertics.Metrics
	workingDirectory string
	ignorePaths      []string
	ignoreExtensions []string
}

func (d directoryScanner) GetType() ScannerType {
	return DirectoryScanner
}

func newDirectoryScanner(sc ScannerConfig) (*directoryScanner, error) {
	if _, err := os.Stat(sc.Basepath); err != nil {
		return nil, err
	}
	metrics := mertics.NewMetrics()
	mc := match.NewMatcherController(sc.Cfg, metrics, sc.Redacted)
	scanner := &directoryScanner{
		mc:               *mc,
		metrics:          metrics,
		workingDirectory: sc.Basepath,
		ignorePaths:      sc.Cfg.IgnorePrefixes,
		ignoreExtensions: sc.Cfg.IgnoreExtensions,
	}
	return scanner, nil
}

func (d directoryScanner) Scan() error {
	return filepath.Walk(d.workingDirectory, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() || shouldIgnore(path, d.ignorePaths, d.ignoreExtensions) {
			return nil
		}
		content, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		return d.mc.Evaluate(path, string(content))
	})
}

func (d directoryScanner) GetMetrics() *mertics.Metrics {
	return d.metrics
}
