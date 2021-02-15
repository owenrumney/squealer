package formatters

import (
	"bytes"
	"fmt"
	"github.com/owenrumney/go-sarif/sarif"

	"github.com/owenrumney/squealer/internal/app/squealer/match"
)

type SarifFormatter struct {
}

func (s SarifFormatter) PrintTransgressions(transgressions []match.Transgression, redacted bool) (string, error) {
	report, err := sarif.New(sarif.Version210)
	if err != nil {
		return "", err
	}

	run := report.AddRun("squealer", "https://github.com/owenrumney/squealer")

	for _, t := range transgressions {
		var content = t.LineContent
		if redacted {
			content = t.RedactedContent
		}
		rule := run.AddRule(t.Hash).
			WithDescription("There should be no sensitive data stored in the repository").
			WithHelp("Add exclude rules to the config for squealer to ignore. Exclude rules take the format filename:hash")

		result := run.AddResult(rule.Id).
			WithMessage(fmt.Sprintf("found transgression [%s], secret hashs [%s]", content, t.Hash)).
			WithLevel("error").
			WithLocationDetails(t.Filename, 1, 1)

		run.AddResultDetails(rule, result, t.Filename)
	}

	var buf bytes.Buffer
	if err = report.PrettyWrite(&buf); err != nil {
		return "", err
	}
	return buf.String(), nil
}
