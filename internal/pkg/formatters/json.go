package formatters

import (
	"encoding/json"

	"github.com/owenrumney/squealer/internal/pkg/match"
)

type JsonFormatter struct {
}

type transgressionsBlock struct {
	Transgressions []transgressionBlock `json:"transgressions"`
}

type committer struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type transgressionBlock struct {
	Content     string    `json:"content"`
	Filename    string    `json:"filename"`
	LineNo      int       `json:"line_number"`
	Hash        string    `json:"secret_hash"`
	Match       string    `json:"match_string"`
	Committer   committer `json:"committer"`
	CommitHash  string    `json:"commit_hash"`
	Committed   string    `json:"committed"`
	ExcludeRule string    `json:"exclude_rule"`
}

func (j JsonFormatter) PrintTransgressions(transgressions []match.Transgression, redacted bool) (string, error) {
	var tb []transgressionBlock

	for _, t := range transgressions {
		var content = t.LineContent
		if redacted {
			content = t.RedactedContent
		}

		tb = append(tb, transgressionBlock{
			Content:  content,
			Filename: t.Filename,
			LineNo:   t.LineNo,
			Hash:     t.Hash,
			Match:    t.Match,
			Committer: committer{
				Name:  t.Committer,
				Email: t.CommitterEmail,
			},
			CommitHash:  t.CommitHash,
			Committed:   t.Committed,
			ExcludeRule: t.ExcludeRule,
		})
	}

	if tb == nil {
		tb = []transgressionBlock{}
	}

	tBlock := &transgressionsBlock{
		Transgressions: tb,
	}

	outBytes, err := json.MarshalIndent(tBlock, "", "  ")
	if err != nil {
		return "", err
	}
	return string(outBytes), err
}
