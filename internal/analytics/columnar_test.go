package analytics

import (
	"fmt"
	"testing"
	"time"
)

func TestColumnarSpeedAndAnalytics(t *testing.T) {
	// 1. Ek columnar table banayi:
	table := NewColumnarTable(10000)

	// 2. 10,000 records insert kiye (Simulating 5+ hours of metrics!):
	now := time.Now().Unix()
	for i := 0; i < 10000; i++ {
		// CPU values oscillate between 10.0% and 50.0%
		cpu := 10.0 + float64(i%40)
		ram := uint64(7000 * 1024 * 1024)
		table.Append(now+int64(i*2), cpu, ram)
	}

	// 3. TIMER SHURU: 10,000 records ka Average aur Max nikalo:
	start := time.Now()

	avgCPU := table.AvgCPU()
	maxCPU := table.MaxCPU()

	duration := time.Since(start)

	// 4. Results Print karo:
	fmt.Println("\n==================================================================")
	fmt.Printf("[CLICKHOUSE COLUMNAR ANALYTICS SPEED TEST]\n")
	fmt.Printf("  Total Records Scanned : %d\n", table.Count())
	fmt.Printf("  Average CPU           : %.2f%%\n", avgCPU)
	fmt.Printf("  Peak Max CPU Spike    : %.2f%%\n", maxCPU)
	fmt.Printf("  TIME TAKEN            : %v (%d Microseconds!)\n", duration, duration.Microseconds())
	fmt.Println("==================================================================")

	// Basic assertion
	if avgCPU == 0 || maxCPU == 0 {
		t.Fatalf("Expected non-zero analytics results")
	}
}

// Benchmark: Memory check (0 B/op hona chahiye!)
func BenchmarkColumnarAvgCPU(b *testing.B) {
	table := NewColumnarTable(10000)
	for i := 0; i < 10000; i++ {
		table.Append(int64(i), 25.5, 7000)
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; b.Loop(); i++ {
		_ = table.AvgCPU()
	}
}