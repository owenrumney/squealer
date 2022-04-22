package formatters

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/owenrumney/squealer/internal/pkg/match"
)

func TestJsonFormatterOutput(t *testing.T) {
	trans := []match.Transgression{createTestTransgression("Joe Bloggs", "joe@bloggs.com", "2001-01-01", "abcd123456efg")}

	plainText, _ := JsonFormatter{}.PrintTransgressions(trans, false)
	assert.Equal(t, `{
  "transgressions": [
    {
      "content": "password=Password1234",
      "filename": "/config.yml",
      "line_number": 10,
      "secret_hash": "sdjn34rf32fds",
      "match_string": "Password1234",
      "committer": {
        "name": "Joe Bloggs",
        "email": "joe@bloggs.com"
      },
      "commit_hash": "abcd123456efg",
      "committed": "2001-01-01",
      "exclude_rule": ""
    }
  ]
}`, plainText)

	redacted, _ := JsonFormatter{}.PrintTransgressions(trans, true)
	assert.Equal(t, `{
  "transgressions": [
    {
      "content": "password=REDACTED",
      "filename": "/config.yml",
      "line_number": 10,
      "secret_hash": "sdjn34rf32fds",
      "match_string": "Password1234",
      "committer": {
        "name": "Joe Bloggs",
        "email": "joe@bloggs.com"
      },
      "commit_hash": "abcd123456efg",
      "committed": "2001-01-01",
      "exclude_rule": ""
    }
  ]
}`, redacted)
}
