package metrics

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConcurrentMetricUpdatesProcessed(t *testing.T) {
	concurrentTests := 5
	iterations := 10

	metrics := runConcurrentMetricTest(concurrentTests, iterations, updateProcessedMetrics)

	assert.Equal(t, int32(iterations), metrics.CommitsProcessed)
	assert.Equal(t, int32(iterations), metrics.FilesProcessed)
	assert.Equal(t, int32(0), metrics.TransgressionsFound)
	assert.Equal(t, int32(0), metrics.TransgressionsIgnored)
	assert.Equal(t, int32(0), metrics.TransgressionsReported)
}

func TestConcurrentMetricUpdatesTransgressions(t *testing.T) {
	concurrentTests := 5
	iterations := 15

	metrics := runConcurrentMetricTest(concurrentTests, iterations, updateTransgressionMetrics)

	assert.Equal(t, int32(0), metrics.CommitsProcessed)
	assert.Equal(t, int32(0), metrics.FilesProcessed)
	assert.Equal(t, int32(iterations), metrics.TransgressionsFound)
	assert.Equal(t, int32(iterations), metrics.TransgressionsIgnored)
	assert.Equal(t, int32(iterations), metrics.TransgressionsReported)
}

func runConcurrentMetricTest(concurrentTests int, iterations int, tester func(metrics *Metrics)) *Metrics {
	metricUpdated := make(chan bool, 10)
	wg := sync.WaitGroup{}
	metrics := NewMetrics()

	wg.Add(concurrentTests)
	for i := 0; i < concurrentTests; i++ {
		go func() {
			for {
				_, ok := <-metricUpdated
				if !ok {
					wg.Done()
					return
				}
				tester(metrics)
			}
		}()
	}

	for i := 0; i < iterations; i++ {
		metricUpdated <- true
	}
	close(metricUpdated)
	wg.Wait()
	return metrics
}

func updateProcessedMetrics(metrics *Metrics) {
	metrics.IncrementCommitsProcessed()
	metrics.IncrementFilesProcessed()
}

func updateTransgressionMetrics(metrics *Metrics) {
	metrics.IncrementTransgressionsFound()
	metrics.IncrementTransgressionsIgnored()
	metrics.IncrementTransgressionsReported()
}
