package formatters

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/owenrumney/squealer/internal/pkg/match"
)

func TestDefaultFormatterOutput(t *testing.T) {
	trans := []match.Transgression{createTestTransgression("Joe Bloggs", "joe@bloggs.com", "2001-01-01", "abcd123456efg")}

	plainText, _ := DefaultFormatter{}.PrintTransgressions(trans, false)
	assert.Equal(t, `
Match Description: │ Some Description
Content:           │ password=Password1234
Filename:          │ /config.yml
Line No:           │ 10
Secret Hash:       │ sdjn34rf32fds
Commit:            │ abcd123456efg
Committer:         │ Joe Bloggs (joe@bloggs.com)
Committed:         │ 2001-01-01
Exclude rule:      │ 
`, plainText)

	redacted, _ := DefaultFormatter{}.PrintTransgressions(trans, true)
	assert.Equal(t, `
Match Description: │ Some Description
Content:           │ password=REDACTED
Filename:          │ /config.yml
Line No:           │ 10
Secret Hash:       │ sdjn34rf32fds
Commit:            │ abcd123456efg
Committer:         │ Joe Bloggs (joe@bloggs.com)
Committed:         │ 2001-01-01
Exclude rule:      │ 
`, redacted)

}
