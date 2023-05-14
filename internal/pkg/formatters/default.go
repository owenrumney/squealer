package formatters

import (
	"strings"

	"github.com/liamg/tml"
	"github.com/owenrumney/squealer/internal/pkg/match"
)

type DefaultFormatter struct {
}

func (d DefaultFormatter) PrintTransgressions(transgressions []match.Transgression, redacted bool) (string, error) {
	builder := strings.Builder{}

	for _, t := range transgressions {
		var content = t.LineContent
		if redacted {
			content = t.RedactedContent
		}
		builder.Write([]byte(tml.Sprintf(`
<blue>Match Description:</blue> <yellow>│</yellow> %s
<blue>Content:</blue>           <yellow>|</yellow> <red>%s</red>
<blue>Filename:</blue>          <yellow>|</yellow> %s
<blue>Line No:</blue>           <yellow>|</yellow> %d
<blue>Secret Hash:</blue>       <yellow>|</yellow> %s
<blue>Commit:</blue>            <yellow>|</yellow> %s
<blue>Committer:</blue>         <yellow>|</yellow> %s (%s)
<blue>Committed:</blue>         <yellow>|</yellow> %s
<blue>Exclude rule:</blue>      <yellow>|</yellow> %s
`, t.MatchDescription, content, t.Filename, t.LineNo, t.Hash, t.CommitHash, t.Committer, t.CommitterEmail, t.Committed, t.ExcludeRule)))
	}
	return builder.String(), nil
}
