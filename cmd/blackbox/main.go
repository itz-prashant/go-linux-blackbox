package main

import (
	"fmt"
	"time"

	"github.com/itz-prashant/linux-blackbox/internal/collector"
)

func main() {
	cpuCollector := collector.NewCPUCollector("/proc/stat")
	memCollector := collector.NewMemoryCollector("/proc/meminfo")

	for {
		cpuMetrics, err := cpuCollector.Collect()

		if err != nil {
			fmt.Printf("CPU Error: %v\n", err)
		}

		memMetrics, err := memCollector.Collect()

		if err != nil {
			fmt.Printf("Mem Error: %v\n", err)
		}

		totalMB := memMetrics.MemTotal / (1024 * 1024)
		availMB := memMetrics.MemAvailable / (1024 * 1024)
		freeMB := memMetrics.MemFree / (1024 * 1024)

		// fmt.Printf("[HEALTH] CPU: %.2f%% | IOwait: %.2f%%\n", cpuMetrics.CPUPercent, cpuMetrics.IOWaitPercent)

		fmt.Printf("[HEALTH] CPU: %5.2f%% | IOwait: %5.2f%% |  RAM: %dMB Avail / %dMB Total (Free: %dMB)\n",
			cpuMetrics.CPUPercent, cpuMetrics.IOWaitPercent, availMB, totalMB, freeMB)

		time.Sleep(2 * time.Second)
	}
}
