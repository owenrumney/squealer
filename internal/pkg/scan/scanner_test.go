package scan

import (
	"fmt"

	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/owenrumney/squealer/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestNewScannerIsGitScanner(t *testing.T) {
	tempdir := t.TempDir()
	dir := fmt.Sprintf("%s/.git", tempdir)
	err := os.MkdirAll(dir, 0600)
	require.NoError(t, err)
	sc := ScannerConfig{
		Cfg:      config.DefaultConfig(),
		Basepath: tempdir,
	}
	scanner, err := NewScanner(sc)
	require.NoError(t, err)
	assert.IsType(t, &gitScanner{}, scanner)
}

func TestNewScannerIsDirectoryScanner(t *testing.T) {
	sc := ScannerConfig{
		Cfg:      config.DefaultConfig(),
		Basepath: "../../../test_resources",
	}
	scanner, err := NewScanner(sc)

	require.NoError(t, err)
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
