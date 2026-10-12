package main

import (
	"fmt"
	"time"

	"github.com/itz-prashant/linux-blackbox/internal/collector"
)

func main() {
	cpuCollector := collector.NewCPUCollector("/proc/stat")

	for {
		metrics, err := cpuCollector.Collect()

		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("[HEALTH] CPU: %.2f%% | IOwait: %.2f%%\n", metrics.CPUPercent, metrics.IOWaitPercent)

		time.Sleep(2 * time.Second)
	}
}
