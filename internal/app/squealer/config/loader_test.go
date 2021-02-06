package config

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
)

func TestEmptyConfigPathReturnsDefault(t *testing.T) {
	config, err := LoadConfig("")

	assert.NoError(t, err)
	assert.Equal(t, DefaultConfig(), config)
}

func TestJsonConfigLoaded(t *testing.T) {
	tempFile, err := ioutil.TempFile(os.TempDir(), "*.json")
	assert.NoError(t, err)
	ioutil.WriteFile(tempFile.Name(), []byte(jsonConfig), 0777)
	defer func() { _ = os.Remove(tempFile.Name()) }()

	config, err := LoadConfig(tempFile.Name())

	assert.NoError(t, err)
	assert.Equal(t, 7, len(config.Rules))
	assert.Equal(t, 1, len(config.Exceptions))
	assert.Equal(t, 2, len(config.IgnorePrefixes))
	assert.Equal(t, 7, len(config.IgnoreExtensions))
}

func TestYamlConfigLoaded(t *testing.T) {
	tempFile, err := ioutil.TempFile(os.TempDir(), "*.yaml")
	assert.NoError(t, err)
	ioutil.WriteFile(tempFile.Name(), []byte(yamlConfig), 0777)
	defer func() { _ = os.Remove(tempFile.Name()) }()

	config, err := LoadConfig(tempFile.Name())

	assert.NoError(t, err)
	assert.Equal(t, 7, len(config.Rules))
	assert.Equal(t, 1, len(config.Exceptions))
	assert.Equal(t, 2, len(config.IgnorePrefixes))
	assert.Equal(t, 7, len(config.IgnoreExtensions))
}

var yamlConfig = `rules:
- rule: (A3T[A-Z0-9]|AKIA|AGPA|AIDA|AROA|AIPA|ANPA|ANVA|ASIA)[A-Z0-9]{16}
  description: Check for AWS Access Key Id
- rule: (?i)aws(.{0,20})?(?-i)['\"][0-9a-zA-Z\/+]{40}['\"]
  description: Check for AWS Secret Access Key
- rule: amzn\.mws\.[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}
  description: Check for AWS MWS Key
- rule: (?i)github[_\-\.]?token[\s:,="\]']+?(?-i)[0-9a-zA-Z]{35,40}
  description: Check for Github Token 
- rule: https://hooks.slack.com/services/T[a-zA-Z0-9_]{8}/B[a-zA-Z0-9_]{8}/[a-zA-Z0-9_]{24}
  description: Check for Slack webhook
- rule: xox[baprs]-([0-9a-zA-Z]{10,48})?
  description: Check for Slack token
- rule: '-----BEGIN ((EC|PGP|DSA|RSA|OPENSSH) )?PRIVATE KEY( BLOCK)?-----'
  description: Check for Private Asymetric Key
ignore_prefixes:
- vendor
- node_modules
ignore_extensions:
- .zip
- .png
- .jpg
- .pdf
- .xls
- .doc
- .docx
exceptions:
- exception: release/update.go:D2IDetI6aidl58GE6dv5uAaWmXM=
  reason: This is a webhook that we got rid of - can be ignored in this file`

var jsonConfig = `
{
  "rules": [
    {
      "rule": "(A3T[A-Z0-9]|AKIA|AGPA|AIDA|AROA|AIPA|ANPA|ANVA|ASIA)[A-Z0-9]{16}",
      "description": "Check for AWS Access Key Id"
    },
    {
      "rule": "(?i)aws(.{0,20})?(?-i)['\\\"][0-9a-zA-Z\\/+]{40}['\\\"]",
      "description": "Check for AWS Secret Access Key"
    },
    {
      "rule": "amzn\\.mws\\.[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}",
      "description": "Check for AWS MWS Key"
    },
    {
      "rule": "(?i)github[_\\-\\.]?token[\\s:,=\"\\]']+?(?-i)[0-9a-zA-Z]{35,40}",
      "description": "Check for Github Token"
    },
    {
      "rule": "https://hooks.slack.com/services/T[a-zA-Z0-9_]{8}/B[a-zA-Z0-9_]{8}/[a-zA-Z0-9_]{24}",
      "description": "Check for Slack webhook"
    },
    {
      "rule": "xox[baprs]-([0-9a-zA-Z]{10,48})?",
      "description": "Check for Slack token"
    },
    {
      "rule": "-----BEGIN ((EC|PGP|DSA|RSA|OPENSSH) )?PRIVATE KEY( BLOCK)?-----",
      "description": "Check for Private Asymetric Key"
    }
  ],
  "ignore_prefixes": [
    "vendor",
    "node_modules"
  ],
  "ignore_extensions": [
    ".zip",
    ".png",
    ".jpg",
    ".pdf",
    ".xls",
    ".doc",
    ".docx"
  ],
  "exceptions": [
    {
      "exception": "release/update.go:D2IDetI6aidl58GE6dv5uAaWmXM=",
      "reason": "This is a webhook that we got rid of - can be ignored in this file"
    }
  ]
}
`
