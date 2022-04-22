package match

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTransgressionUpdate(t *testing.T) {
	t1 := createTestTransgression("Joe Bloggs", "joe@bloggs.com", "2001-01-01", "abcd")
	t2 := createTestTransgression("Thom Thumb", "joe@bloggs.com", "2001-12-01", "1234")

	assert.Equal(t, "Joe Bloggs", t1.Committer)
	assert.Equal(t, "joe@bloggs.com", t1.CommitterEmail)
	assert.Equal(t, "2001-01-01", t1.Committed)
	assert.Equal(t, "abcd", t1.CommitHash)

	t1.update(t2)

	assert.Equal(t, "Thom Thumb", t1.Committer)
	assert.Equal(t, "joe@bloggs.com", t1.CommitterEmail)
	assert.Equal(t, "2001-12-01", t1.Committed)
	assert.Equal(t, "1234", t1.CommitHash)

}

func createTestTransgression(committer, committerEmail, committed, commitHash string) Transgression {
	return Transgression{
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
