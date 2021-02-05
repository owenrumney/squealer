package mertics

import (
	"sync/atomic"
)

type Metrics struct {
	CommitsProcessed       int32
	FilesProcessed         int32
	TransgressionsFound    int32
	TransgressionsIgnored  int32
	TransgressionsReported int32
}

func NewMetrics() *Metrics {
	return &Metrics{}
}

func (m *Metrics) UpdateCommitsProcessed() {
	atomic.AddInt32(&m.CommitsProcessed, 1)
}

func (m *Metrics) UpdateFilesProcessed() {
	atomic.AddInt32(&m.FilesProcessed, 1)
}

func (m *Metrics) UpdateTransgressionsFound() {
	atomic.AddInt32(&m.TransgressionsFound, 1)
}


func (m *Metrics) UpdateTransgressionsIgnored() {
	atomic.AddInt32(&m.TransgressionsIgnored, 1)
}

func (m *Metrics) UpdateTransgressionsReported() {
	atomic.AddInt32(&m.TransgressionsReported, 1)
}

