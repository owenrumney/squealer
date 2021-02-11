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
		Cfg:      config.DefaultConfig(),
		Basepath: fmt.Sprintf("%s/src/github.com/owenrumney/squealer/", gopath),
	}
	scanner, err := NewScanner(sc)
	assert.NoError(t, err)
	assert.IsType(t, &gitScanner{}, scanner)
}

func TestNewScannerIsDirectoryScanner(t *testing.T) {
	sc := ScannerConfig{
		Cfg:      config.DefaultConfig(),
		Basepath: "../../../../test_resources",
	}
	scanner, err := NewScanner(sc)

	assert.NoError(t, err)
	assert.IsType(t, &directoryScanner{}, scanner)
}

func TestShouldIgnore(t *testing.T) {
	ignorePaths := []string{"vendor", "npm_modules"}
	ignoreExtensions := []string{"zip"}
	assert.True(t, shouldIgnore("/src/scan/vendor/github.com", ignorePaths, ignoreExtensions))
	assert.True(t, shouldIgnore("/src/scan/npm_modules/github.com", ignorePaths, ignoreExtensions))
	assert.True(t, shouldIgnore("vendor/github.com", ignorePaths, ignoreExtensions))
	assert.True(t, shouldIgnore("npm_modules/github.com", ignorePaths, ignoreExtensions))
	assert.True(t, shouldIgnore("pingu.zip", ignorePaths, ignoreExtensions))
	assert.False(t, shouldIgnore("src/scan", ignorePaths, ignoreExtensions))
	assert.False(t, shouldIgnore("test", ignorePaths, ignoreExtensions))
	assert.False(t, shouldIgnore("govendor", ignorePaths, ignoreExtensions))
	assert.False(t, shouldIgnore("pingu.honk", ignorePaths, ignoreExtensions))
}
