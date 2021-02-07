package main

import (
	"fmt"
	"github.com/owenrumney/squealer/internal/app/squealer/config"
	"github.com/owenrumney/squealer/internal/app/squealer/mertics"
	"github.com/owenrumney/squealer/internal/app/squealer/scan"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestPrintMetrics(t *testing.T) {
	m := mertics.Metrics{
		CommitsProcessed:       1,
		FilesProcessed:         2,
		TransgressionsFound:    3,
		TransgressionsIgnored:  4,
		TransgressionsReported: 5,
	}
	m.StartTimer()
	m.StopTimer()

	output := printMetrics(&m)
	assert.Equal(t, `
Processing:
  duration:     0.00s
  commits:      1
  commit files: 2

transgressionMap:
  identified:   3
  ignored:      4
  reported:     5

`, output)
}

func TestNewScannerIsGitScanner(t *testing.T) {
	gopath := os.Getenv("GOPATH")
	scanner := getScanner(config.DefaultConfig(), fmt.Sprintf("%s/src/github.com/owenrumney/squealer/", gopath))
	assert.Equal(t, scan.GitScanner, scanner.GetType())
}

func TestNewScannerIsDirectoryScanner(t *testing.T) {
	gopath := os.Getenv("GOPATH")
	scanner := getScanner(config.DefaultConfig(), fmt.Sprintf("%s/src/github.com/owenrumney/squealer/test_resources", gopath))
	assert.Equal(t, scan.DirectoryScanner, scanner.GetType())
}
