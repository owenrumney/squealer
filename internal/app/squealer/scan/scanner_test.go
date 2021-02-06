package scan

import (
	"fmt"
	"github.com/owenrumney/squealer/internal/app/squealer/config"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestNewScannerIsGitScanner(t *testing.T) {
	gopath := os.Getenv("GOPATH")
	sc := ScannerConfig{
		cfg:      config.DefaultConfig(),
		basepath: fmt.Sprintf("%s/src/github.com/owenrumney/squealer/", gopath),
	}
	scanner, err := NewScanner(sc)
	fmt.Println(sc.basepath)
	assert.NoError(t, err)
	assert.IsType(t, &gitScanner{}, scanner)
}

func TestNewScannerIsDirectoryScanner(t *testing.T) {
	sc := ScannerConfig{
		cfg:      config.DefaultConfig(),
		basepath: "../../../../test_resources",
	}
	scanner, err := NewScanner(sc)

	assert.NoError(t, err)
	assert.IsType(t, &directoryScanner{}, scanner)
}
