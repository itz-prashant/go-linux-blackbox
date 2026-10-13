package main

import (
	"fmt"
	"time"

	"github.com/itz-prashant/linux-blackbox/internal/collector"
)

func main() {
	cpuCollector := collector.NewCPUCollector("/proc/stat")
	memCollector := collector.NewMemoryCollector("/proc/meminfo")
	diskCollector := collector.NewDiskCollector("/")

	for {
		cpuMetrics, err := cpuCollector.Collect()

		if err != nil {
			fmt.Printf("CPU Error: %v\n", err)
		}

		memMetrics, err := memCollector.Collect()

		if err != nil {
			fmt.Printf("Mem Error: %v\n", err)
		}

		diskMetrics, err := diskCollector.Collect()

		if err != nil {
			fmt.Printf("disk Error: %v\n", err)
		}

		totalMB := memMetrics.MemTotal / (1024 * 1024)
		availMB := memMetrics.MemAvailable / (1024 * 1024)
		freeMB := memMetrics.MemFree / (1024 * 1024)

		totalGB := diskMetrics.DiskTotal / (1024 * 1024 * 1024)
		freeGB := diskMetrics.DiskFree / (1024 * 1024 * 1024)

		// fmt.Printf("[HEALTH] CPU: %.2f%% | IOwait: %.2f%%\n", cpuMetrics.CPUPercent, cpuMetrics.IOWaitPercent)

		// fmt.Printf("[HEALTH] CPU: %5.2f%% | IOwait: %5.2f%% |  RAM: %dMB Avail / %dMB Total (Free: %dMB)\n",
		// 	cpuMetrics.CPUPercent, cpuMetrics.IOWaitPercent, availMB, totalMB, freeMB)

		fmt.Printf("[HEALTH] CPU: %5.2f%% (IO: %4.2f%%) | RAM: %4dMB / %4dMB Avail (Free: %dMB) | DISK: %3dGB Free / %3dGB (%4.1f%% used)\n",
			cpuMetrics.CPUPercent,
			cpuMetrics.IOWaitPercent,
			availMB,
			totalMB,
			freeMB,
			freeGB,
			totalGB,
			diskMetrics.DiskUsedPercent,
		)

		time.Sleep(2 * time.Second)
	}
}
