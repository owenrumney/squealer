package squealer

import (
	"github.com/owenrumney/squealer/internal/app/squealer/match"
	"github.com/owenrumney/squealer/pkg/config"
	"github.com/owenrumney/squealer/pkg/result"
)

type stringScanner struct {
	mc match.MatcherController
}

func NewStringScanner() *stringScanner {
	return NewStringScannerWithConfig(config.DefaultConfig())
}

func NewStringScannerWithConfig(conf *config.Config) *stringScanner {
	mc := match.NewMatcherController(conf, nil, true)

	return &stringScanner{
		mc: *mc,
	}
}

func (s stringScanner) Scan(content string) result.StringScanResult {
	return s.mc.EvaluateString(content)
}
