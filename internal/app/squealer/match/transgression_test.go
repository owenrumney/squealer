package match

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTransgressionOutputString(t *testing.T) {
	tr := createTestTransgression()

	assert.Equal(t, `
content:      | password=Password1234
filename:     | /config.yml
secret hash:  | sdjn34rf32fds
commit:       | 
committer:    |  ()
committed:    |
exclude rule: | 
	`, tr.String())
}

func TestTransgressionOutputRedacted(t *testing.T) {
	tr := createTestTransgression()

	assert.Equal(t, `
content:      | password=REDACTED
filename:     | /config.yml
secret hash:  | sdjn34rf32fds
commit:       | 
committer:    |  ()
committed:    |
exclude rule: | 
	`, tr.Redacted())
}

func createTestTransgression() Transgression {
	return Transgression{
		lineContent: "password=Password1234",
		filename:    "/config.yml",
		hash:        "sdjn34rf32fds",
		match:       "Password1234",
		redacted:    "password=REDACTED",
	}
}
