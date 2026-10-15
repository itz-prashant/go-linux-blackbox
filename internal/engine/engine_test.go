package engine

import (
	"testing"
	"time"

	"github.com/itz-prashant/linux-blackbox/internal/sentry"
	"github.com/itz-prashant/linux-blackbox/internal/tape"
)

func BenchmarkEngineProcess(b *testing.B) {
	eventChan := make(chan sentry.LogEvent, 100)
	ring := tape.NewRingBuffer(100)
	folder := tape.NewEventFolder(ring, 5*time.Second)
	eng := NewEngine(eventChan, 4, folder)

	dummyEvent := sentry.LogEvent{
		Line:      "FATAL: Database connection timeout on port 5432",
		Inode:     77,
		Timestamp: time.Now(),
	}
	b.ResetTimer()   // Setup ka time ignore karo, sirf processing ka time count karo:
	b.ReportAllocs() // HAR OPERATION MEIN MEMORY ALLOCATIONS COUNT KARO!

	for b.Loop() {
		eng.process(dummyEvent, 1)
	}
}
