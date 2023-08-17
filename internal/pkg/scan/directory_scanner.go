package scan

import (
	"os"
	"path/filepath"

	"github.com/owenrumney/squealer/internal/pkg/match"
	"github.com/owenrumney/squealer/internal/pkg/metrics"
)

type directoryScanner struct {
	mc               match.MatcherController
	metrics          *metrics.Metrics
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
	metrics := metrics.NewMetrics()
	mc := match.NewMatcherController(sc.Cfg, metrics, sc.Redacted)
	scanner := &directoryScanner{
		mc:               *mc,
		metrics:          metrics,
		workingDirectory: sc.Basepath,
		ignorePaths:      sc.Cfg.IgnorePaths,
		ignoreExtensions: sc.Cfg.IgnoreExtensions,
	}
	return scanner, nil
}

func (d directoryScanner) Scan() ([]match.Transgression, error) {
	if err := filepath.Walk(d.workingDirectory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || shouldIgnore(path, d.ignorePaths, d.ignoreExtensions) {
			return nil
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		d.metrics.IncrementFilesProcessed()
		return d.mc.Evaluate(path, string(content), nil)
	}); err != nil {
		return nil, err
	}
	return d.mc.Transgressions(), nil
}

func (d directoryScanner) GetMetrics() *metrics.Metrics {
	return d.metrics
}
