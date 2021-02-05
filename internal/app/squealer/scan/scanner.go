package scan

import "github.com/owenrumney/squealer/internal/app/squealer/config"

type Scanner interface {
	Scan() error
}

type ScannerConfig struct {
	cfg *config.Config
	basepath string
	redacted bool
}

func NewScannerConfig(basepath string, redacted bool, cfg *config.Config) ScannerConfig {
	return ScannerConfig{
		cfg:      cfg,
		basepath: basepath,
		redacted: redacted,
	}
}