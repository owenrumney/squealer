package formatters

import (
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/owenrumney/squealer/internal/app/squealer/match"
)

func TestSarifFormmaterOutput(t *testing.T) {
	trans := []match.Transgression{createTestTransgression("Joe Bloggs", "joe@bloggs.com", "2001-01-01", "abcd123456efg")}

	plainText, _ := SarifFormatter{}.PrintTransgressions(trans, false)
	assert.Equal(t, `{
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
          }
        }
      ],
      "results": [
        {
          "level": "error",
          "message": {
            "text": "found transgression [password=Password1234], secret hashs [sdjn34rf32fds]"
          },
          "ruleId": "sdjn34rf32fds",
          "ruleIndex": 0,
          "locations": [
            {
              "physicalLocation": {
                "artifactLocation": {
                  "uri": "/config.yml",
                  "index": 0
                },
                "region": {
                  "startLine": 1,
                  "startColumn": 1
                }
              }
            }
          ]
        }
      ]
    }
  ]
}`, plainText)

	redacted, _ := SarifFormatter{}.PrintTransgressions(trans, true)
	assert.Equal(t, `{
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
          }
        }
      ],
      "results": [
        {
          "level": "error",
          "message": {
            "text": "found transgression [password=REDACTED], secret hashs [sdjn34rf32fds]"
          },
          "ruleId": "sdjn34rf32fds",
          "ruleIndex": 0,
          "locations": [
            {
              "physicalLocation": {
                "artifactLocation": {
                  "uri": "/config.yml",
                  "index": 0
                },
                "region": {
                  "startLine": 1,
                  "startColumn": 1
                }
              }
            }
          ]
        }
      ]
    }
  ]
}`, redacted)
}
