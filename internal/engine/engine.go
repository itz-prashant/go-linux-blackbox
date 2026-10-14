package engine

import (
	"context"
	"strconv"
	"sync"

	"github.com/itz-prashant/linux-blackbox/internal/sentry"
)

type Engine struct {
	inchan      <-chan sentry.LogEvent
	pool        *BufferPool
	metrics     *EngineMetrics
	workerCount int
	wg          sync.WaitGroup
}

func NewEngine(inchan <-chan sentry.LogEvent, workerCount int) *Engine {
	return &Engine{
		inchan:      inchan,
		pool:        NewBufferPool(),
		metrics:     &EngineMetrics{},
		workerCount: workerCount,
	}
}

// Start saare workers ko background goroutines mein launch karta hai

func (e *Engine) Start(ctx context.Context) {
	for i := 0; i < e.workerCount; i++ {
		e.wg.Add(1)
		go e.worker(ctx, i+1)
	}
}

// worker har incoming event ko channel se uthata hai
func (e *Engine) worker(ctx context.Context, id int) {
	defer e.wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case event, ok := <-e.inchan:
			if !ok {
				return
			}
			e.process(event, id)
		}
	}
}

// process sync.Pool ka recycled buffer use karke bina heap allocation ke log process karta hai
func (e *Engine) process(event sentry.LogEvent, workerID int) {
	buf := e.pool.Get()

	// 2. Log format kiya buffer ke andar:
	var numBuf [20]byte
	// 3. Direct Zero-Allocation Writing (Bypassing fmt.Fprintf):
	buf.WriteString("[WORKER ")
	buf.Write(strconv.AppendInt(numBuf[:0], int64(workerID), 10))
	buf.WriteString("] Inode #")
	buf.Write(strconv.AppendUint(numBuf[:0], event.Inode, 10))
	buf.WriteString(": ")
	buf.WriteString(event.Line)
	// 4. Atomic counters
	e.metrics.IncrementProcessed()
	e.metrics.AddBytes(uint64(buf.Len()))
	// 5. Buffer wapas pool mein
	e.pool.Put(buf)
}

// Wait program band hote waqt sabhi workers ka kaam khatam hone ka intezaar karta hai
func (e *Engine) Wait() {
	e.wg.Wait()
}

// Metrics engine ke current atomic stats return karta hai
func (e *Engine) Metrics() *EngineMetrics {
	return e.metrics
}
