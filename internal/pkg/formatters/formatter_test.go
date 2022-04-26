package formatters

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/owenrumney/squealer/internal/pkg/match"
)

func TestGetFormatter(t *testing.T) {
	assert.IsType(t, GetFormatter("json"), &JsonFormatter{})
	assert.IsType(t, GetFormatter("sarif"), &SarifFormatter{})
	assert.IsType(t, GetFormatter("default"), &DefaultFormatter{})
	assert.IsType(t, GetFormatter("text"), &DefaultFormatter{})
}

func createTestTransgression(committer, committerEmail, committed, commitHash string) match.Transgression {
	return match.Transgression{
		LineNo:           10,
		LineContent:      "password=Password1234",
		Filename:         "/config.yml",
		Hash:             "sdjn34rf32fds",
		Match:            "Password1234",
		MatchDescription: "Some Description",
		RedactedContent:  "password=REDACTED",
		CommitterEmail:   committerEmail,
		Committer:        committer,
		CommitHash:       commitHash,
		Committed:        committed,
	}
}
