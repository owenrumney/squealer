package scan

import (
	"github.com/owenrumney/squealer/internal/app/squealer/match"
	"github.com/owenrumney/squealer/internal/app/squealer/mertics"
	"os"
)

type directoryScanner struct {
	mc               match.MatcherController
	metrics          *mertics.Metrics
	workingDirectory string
	ignorePaths      []string
	ignoreExtensions []string
}

func newDirectoryScanner(sc ScannerConfig) (*directoryScanner, error) {
	if _, err := os.Stat(sc.basepath); err != nil {
		return nil, err
	}
	metrics := mertics.NewMetrics()
	mc := match.NewMatcherController(sc.cfg, metrics, sc.redacted)
	scanner := &directoryScanner{
		mc:               *mc,
		metrics:          metrics,
		workingDirectory: sc.basepath,
		ignorePaths:      sc.cfg.IgnorePrefixes,
		ignoreExtensions: sc.cfg.IgnoreExtensions,
	}
	return scanner, nil
}

func (d directoryScanner) Scan() error {
	panic("implement me")
}

func (d directoryScanner) GetMetrics() *mertics.Metrics {
	panic("implement me")
}
