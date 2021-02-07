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
hash:         | sdjn34rf32fds
exclude rule: | /config.yml:sdjn34rf32fds
	`, tr.String())
}

func TestTransgressionOutputRedacted(t *testing.T) {
	tr := createTestTransgression()

	assert.Equal(t, `
content:      | password=REDACTED
filename:     | /config.yml
hash:         | sdjn34rf32fds
exclude rule: | /config.yml:sdjn34rf32fds
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
