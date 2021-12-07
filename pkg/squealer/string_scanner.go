package squealer

import (
	"github.com/owenrum/squealer/internal/app/squealer/match"
	"github.com/owenrum/squealer/pkg/config"
	"github.com/owenrum/squealer/pkg/result"
)

type StringScanner struct {
	mc match.MatcherController
}

func NewStringScanner() *StringScanner {
	return NewStringScannerWithConfig(config.DefaultConfig())
}

func NewStringScannerWithConfig(conf *config.Config) *StringScanner {
	mc := match.NewMatcherController(conf, nil, true)

	return &StringScanner{
		mc: *mc,
	}
}

func (s StringScanner) Scan(content string) result.StringScanResult {
	return s.mc.EvaluateString(content)
}
