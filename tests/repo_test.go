package tests

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/owenrumney/squealer/internal/pkg/scan"
	"github.com/owenrumney/squealer/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestRepoEndToEnd(t *testing.T) {
	scanner, err := scan.NewScanner(scan.ScannerConfig{
		Cfg:      config.DefaultConfig(),
		Basepath: gitTestPath,
		Redacted: true,
	})

	require.NoError(t, err)
	_, err = scanner.Scan()
	require.NoError(t, err)

	metrics := scanner.GetMetrics()
	assert.Equal(t, int32(3), metrics.CommitsProcessed)
	assert.Equal(t, int32(7), metrics.TransgressionsFound)
	assert.Equal(t, int32(0), metrics.TransgressionsIgnored)
	assert.Equal(t, int32(4), metrics.TransgressionsReported)
}

func TestDirEndToEnd(t *testing.T) {
	scanner, err := scan.NewScanner(scan.ScannerConfig{
		Cfg:      config.DefaultConfig(),
		Basepath: dirTestPath,
		Redacted: true,
	})

	require.NoError(t, err)
	_, err = scanner.Scan()
	require.NoError(t, err)

	metrics := scanner.GetMetrics()
	assert.Equal(t, int32(0), metrics.CommitsProcessed)
	assert.Equal(t, int32(6), metrics.TransgressionsFound)
	assert.Equal(t, int32(0), metrics.TransgressionsIgnored)
	assert.Equal(t, int32(3), metrics.TransgressionsReported)
}

func TestRepoEndToEndWithEverything(t *testing.T) {
	scanner, err := scan.NewScanner(scan.ScannerConfig{
		Cfg:        config.DefaultConfig(),
		Basepath:   gitTestPath,
		Redacted:   true,
		Everything: true,
	})

	require.NoError(t, err)
	_, err = scanner.Scan()
	require.NoError(t, err)

	metrics := scanner.GetMetrics()
	assert.Equal(t, int32(3), metrics.CommitsProcessed)
	assert.Equal(t, int32(7), metrics.TransgressionsFound)
	assert.Equal(t, int32(0), metrics.TransgressionsIgnored)
	assert.Equal(t, int32(4), metrics.TransgressionsReported)
}
