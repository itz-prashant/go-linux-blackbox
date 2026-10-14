package engine

import "sync/atomic"

type EngineMetrics struct {
	totalProcessed uint64
	totalBytes     uint64
	droppedEvents  uint64
}

func (m *EngineMetrics) IncrementProcessed() {
	atomic.AddUint64(&m.totalProcessed, 1)
}

func (m *EngineMetrics) AddBytes(n uint64) {
	atomic.AddUint64(&m.totalBytes, n)
}

func (m *EngineMetrics) IncrementDropped() {
	atomic.AddUint64(&m.droppedEvents, 1)
}

func (m *EngineMetrics) GetStats() (prcessed, bytes, dropped uint64) {
	return atomic.LoadUint64(&m.totalProcessed),
		atomic.LoadUint64(&m.totalBytes),
		atomic.LoadUint64(&m.droppedEvents)
}
