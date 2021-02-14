package formatters

import (
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/owenrumney/squealer/internal/app/squealer/match"
)

func TestDefaultFormatterOutput(t *testing.T) {
	trans := []match.Transgression{createTestTransgression("Joe Bloggs", "joe@bloggs.com", "2001-01-01", "abcd123456efg")}

	plainText, _ := DefaultFormatter{}.PrintTransgressions(trans, false)
	assert.Equal(t, `
content:      | password=Password1234
Filename:     | /config.yml
secret Hash:  | sdjn34rf32fds
commit:       | abcd123456efg
Committer:    | Joe Bloggs (joe@bloggs.com)
Committed:    | 2001-01-01
exclude rule: | 
`, plainText)

	redacted, _ := DefaultFormatter{}.PrintTransgressions(trans, true)
	assert.Equal(t, `
content:      | password=REDACTED
Filename:     | /config.yml
secret Hash:  | sdjn34rf32fds
commit:       | abcd123456efg
Committer:    | Joe Bloggs (joe@bloggs.com)
Committed:    | 2001-01-01
exclude rule: | 
`, redacted)

}
