package formatters

import (
	"bytes"
	"fmt"

	"github.com/owenrumney/go-sarif/sarif"

	"github.com/owenrumney/squealer/internal/pkg/match"
)

type SarifFormatter struct {
}

func (s SarifFormatter) PrintTransgressions(transgressions []match.Transgression, redacted bool) (string, error) {
	if len(transgressions) == 0 {
		return "", nil
	}

	report, err := sarif.New(sarif.Version210)
	if err != nil {
		return "", err
	}

	run := sarif.NewRun("squealer", "https://github.com/owenrumney/squealer")

	for _, t := range transgressions {
		var content = t.LineContent
		if redacted {
			content = t.RedactedContent
		}
		rule := run.AddRule(t.Hash).
			WithDescription("There should be no sensitive data stored in the repository").
			WithHelp("Add exclude rules to the config for squealer to ignore. Exclude rules take the format filename:hash")

		run.AddDistinctArtifact(t.Filename)

		result := run.AddResult(rule.ID).
			WithMessage(sarif.NewTextMessage(fmt.Sprintf("found transgression [%s], secret hashs [%s]", content, t.Hash))).
			WithLevel("error").
			WithLocation(
				sarif.NewLocationWithPhysicalLocation(
					sarif.NewPhysicalLocation().
						WithArtifactLocation(
							sarif.NewSimpleArtifactLocation(t.Filename),
						).WithRegion(
						sarif.NewRegion().
							WithStartLine(t.LineNo).
							WithStartColumn(1),
					),
				),
			)

		pb := sarif.NewPropertyBag()
		pb.Add("commit", t.CommitHash)
		pb.Add("committed", t.Committed)
		pb.Add("committer", t.Committer)
		result.AttachPropertyBag(pb)

	}

	report.AddRun(run)

	var buf bytes.Buffer
	if err = report.PrettyWrite(&buf); err != nil {
		return "", err
	}
	return buf.String(), nil
}
