package tests

import (
	"testing"

	"github.com/owenrum/squealer/internal/app/squealer/scan"
	"github.com/owenrum/squealer/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestRepoEndToEnd(t *testing.T) {
	scanner, err := scan.NewScanner(scan.ScannerConfig{
		Cfg:      config.DefaultConfig(),
		Basepath: gitTestPath,
		Redacted: true,
	})

	assert.NoError(t, err)
	_, err = scanner.Scan()
	assert.NoError(t, err)

	metrics := scanner.GetMetrics()
	assert.Equal(t, int32(3), metrics.CommitsProcessed)
	assert.Equal(t, int32(4), metrics.TransgressionsFound)
	assert.Equal(t, int32(0), metrics.TransgressionsIgnored)
	assert.Equal(t, int32(4), metrics.TransgressionsReported)
}

func TestDirEndToEnd(t *testing.T) {
	scanner, err := scan.NewScanner(scan.ScannerConfig{
		Cfg:      config.DefaultConfig(),
		Basepath: dirTestPath,
		Redacted: true,
	})

	assert.NoError(t, err)
	_, err = scanner.Scan()
	assert.NoError(t, err)

	metrics := scanner.GetMetrics()
	assert.Equal(t, int32(0), metrics.CommitsProcessed)
	assert.Equal(t, int32(3), metrics.TransgressionsFound)
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

	assert.NoError(t, err)
	_, err = scanner.Scan()
	assert.NoError(t, err)

	metrics := scanner.GetMetrics()
	assert.Equal(t, int32(3), metrics.CommitsProcessed)
	assert.Equal(t, int32(4), metrics.TransgressionsFound)
	assert.Equal(t, int32(0), metrics.TransgressionsIgnored)
	assert.Equal(t, int32(4), metrics.TransgressionsReported)
}
