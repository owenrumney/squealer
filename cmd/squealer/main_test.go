package main

import (
	"github.com/owenrumney/squealer/internal/app/squealer/mertics"
	"github.com/stretchr/testify/assert"
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
