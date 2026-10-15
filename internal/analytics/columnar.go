package analytics

import "sync"

type ColumnarTable struct {
	mu         sync.RWMutex
	timeStamps []int64
	cpuValues  []float64
	ramValues  []uint64
}

func NewColumnarTable(capacity int) *ColumnarTable {
	return &ColumnarTable{
		timeStamps: make([]int64, 0, capacity),
		cpuValues:  make([]float64, 0, capacity),
		ramValues:  make([]uint64, 0, capacity),
	}
}

func (t *ColumnarTable) Append(ts int64, cpu float64, ram uint64) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.timeStamps = append(t.timeStamps, ts)
	t.cpuValues = append(t.cpuValues, cpu)
	t.ramValues = append(t.ramValues, ram)
}

func (t *ColumnarTable) AvgCPU() float64 {
	t.mu.RLock()
	defer t.mu.RUnlock()

	if len(t.cpuValues) == 0 {
		return 0.0
	}
	var sum float64
	for _,v := range t.cpuValues {
		sum += v
	}
	return sum / float64(len(t.cpuValues))
}

func (t *ColumnarTable) MaxCPU() float64 {
    t.mu.RLock()
    defer t.mu.RUnlock()

    var max float64
    for _, v := range t.cpuValues {
        if v > max {
            max = v
        }
    }
    return max
}

func (t *ColumnarTable) Count() int {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return len(t.cpuValues)
}