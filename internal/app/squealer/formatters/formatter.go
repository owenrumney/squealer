package formatters

import (
	"github.com/owenrumney/squealer/internal/app/squealer/match"
)

type Formatter interface {
	PrintTransgressions([]match.Transgression, bool) (string, error)
}

func GetFormatter(format string) Formatter {
	switch format {
	case "sarif":
		return &SarifFormatter{}
	case "json":
		return &JsonFormatter{}
	default:
		return &DefaultFormatter{}
	}
}
