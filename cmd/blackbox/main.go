package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/itz-prashant/linux-blackbox/internal/collector"
	"github.com/itz-prashant/linux-blackbox/internal/sentry"
)

func main() {
	cpuCollector := collector.NewCPUCollector("/proc/stat")
	memCollector := collector.NewMemoryCollector("/proc/meminfo")
	diskCollector := collector.NewDiskCollector("/")

	logPath := "app.log"
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		_ = os.WriteFile(logPath, []byte(""), 0644)
	}

	fileSentry, err := sentry.NewSentry(logPath)
	if err != nil {
		fmt.Printf("Fatal: Sentry init failed: %v\n", err)
		return
	}
	defer fileSentry.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	eventChan, errChan := fileSentry.Watch(ctx)

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	fmt.Printf("[READY] Sentry is actively watching '%s' on Inode #%d...\n\n", logPath, fileSentry.Inode)

	for {
		select {
		case <-sigChan:
			fmt.Println("\n\n[SHUTDOWN] Received SIGINT (Ctrl+C). Flushing buffers...")
			cancel() // Sentry goroutine ko stop signal bheja

			time.Sleep(200 * time.Millisecond)
			fmt.Println("[SHUTDOWN] Inotify closed. Blackbox stopped gracefully. Bye!")
			return
		case err := <-errChan:
			if err != nil {
				fmt.Printf("[ERROR] Sentry: %v\n", err)
			}
		case event, ok := <-eventChan:
			if ok {
				fmt.Printf("\n>>> [EVENT CAPTURED] Inode #%d | Time: %s\n    Log: \"%s\"\n\n",
					event.Inode, event.Timestamp.Format("15:04:05"), event.Line)
			}

		case <-ticker.C:
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
		}
	}
}
