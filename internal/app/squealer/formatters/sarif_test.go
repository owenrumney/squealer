package formatters

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/owenrumney/squealer/internal/app/squealer/match"
)

func TestSarifFormmaterOutput(t *testing.T) {
	trans := []match.Transgression{createTestTransgression("Joe Bloggs", "joe@bloggs.com", "2001-01-01", "abcd123456efg")}

	plainText, _ := SarifFormatter{}.PrintTransgressions(trans, false)
	expected := `{
  "version": "2.1.0",
  "$schema": "http://json.schemastore.org/sarif-2.1.0-rtm.4",
  "runs": [
    {
      "tool": {
        "driver": {
          "name": "squealer",
          "informationUri": "https://github.com/owenrumney/squealer",
          "rules": [
            {
              "id": "sdjn34rf32fds",
              "shortDescription": {
                "text": "There should be no sensitive data stored in the repository"
              },
              "help": {
                "text": "Add exclude rules to the config for squealer to ignore. Exclude rules take the format filename:hash"
              }
            }
          ]
        }
      },
      "artifacts": [
        {
          "location": {
            "uri": "/config.yml"
          },
          "length": 0
        }
      ],
      "results": [
        {
          "ruleId": "sdjn34rf32fds",
          "ruleIndex": 0,
          "level": "error",
          "message": {
            "text": "found transgression [password=Password1234], secret hashs [sdjn34rf32fds]"
          },
          "locations": [
            {
              "physicalLocation": {
                "artifactLocation": {
                  "uri": "/config.yml",
                  "index": 0
                },
                "region": {
                  "startLine": 10,
                  "startColumn": 1
                }
              }
            }
          ]
        }
      ]
    }
  ]
}`
	assert.Equal(t, expected, plainText)

	redacted, _ := SarifFormatter{}.PrintTransgressions(trans, true)
	expected = `{
  "version": "2.1.0",
  "$schema": "http://json.schemastore.org/sarif-2.1.0-rtm.4",
  "runs": [
    {
      "tool": {
        "driver": {
          "name": "squealer",
          "informationUri": "https://github.com/owenrumney/squealer",
          "rules": [
            {
              "id": "sdjn34rf32fds",
              "shortDescription": {
                "text": "There should be no sensitive data stored in the repository"
              },
              "help": {
                "text": "Add exclude rules to the config for squealer to ignore. Exclude rules take the format filename:hash"
              }
            }
          ]
        }
      },
      "artifacts": [
        {
          "location": {
            "uri": "/config.yml"
          },
          "length": 0
        }
      ],
      "results": [
        {
          "ruleId": "sdjn34rf32fds",
          "ruleIndex": 0,
          "level": "error",
          "message": {
            "text": "found transgression [password=REDACTED], secret hashs [sdjn34rf32fds]"
          },
          "locations": [
            {
              "physicalLocation": {
                "artifactLocation": {
                  "uri": "/config.yml",
                  "index": 0
                },
                "region": {
                  "startLine": 10,
                  "startColumn": 1
                }
              }
            }
          ]
        }
      ]
    }
  ]
}`
	assert.Equal(t, expected, redacted)
}
