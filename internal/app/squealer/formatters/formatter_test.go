package formatters

import (
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/owenrumney/squealer/internal/app/squealer/match"
)

func TestGetFormatter(t *testing.T) {
	assert.IsType(t, GetFormatter("json"), &JsonFormatter{})
	assert.IsType(t, GetFormatter("sarif"), &SarifFormatter{})
	assert.IsType(t, GetFormatter("default"), &DefaultFormatter{})
	assert.IsType(t, GetFormatter("text"), &DefaultFormatter{})
}

func createTestTransgression(committer, committerEmail, committed, commitHash string) match.Transgression {
	return match.Transgression{
		LineContent:     "password=Password1234",
		Filename:        "/config.yml",
		Hash:            "sdjn34rf32fds",
		Match:           "Password1234",
		RedactedContent: "password=REDACTED",
		CommitterEmail:  committerEmail,
		Committer:       committer,
		CommitHash:      commitHash,
		Committed:       committed,
	}
}
