package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/itz-prashant/linux-blackbox/internal/collector"
	"github.com/itz-prashant/linux-blackbox/internal/engine"
	"github.com/itz-prashant/linux-blackbox/internal/sentry"
	"github.com/itz-prashant/linux-blackbox/internal/tape"
)

func main() {
	fmt.Println("==================================================================================")
	fmt.Println("               THE LINUX BLACKBOX - FLIGHT RECORDER (MILESTONE 1 - 4)             ")
	fmt.Println("==================================================================================")

	// 1. Health Collectors (M1)
	cpuCollector := collector.NewCPUCollector("/proc/stat")
	memCollector := collector.NewMemoryCollector("/proc/meminfo")
	diskCollector := collector.NewDiskCollector("/")

	// 2. Dummy Log File Setup
	logPath := "app.log"
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		_ = os.WriteFile(logPath, []byte(""), 0644)
	}

	// 3. Sentry Setup (M2)
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

	eventsChan, errChan := fileSentry.Watch(ctx)

	// 4. Milestone 4: Fixed-Memory Ring Buffer & 5-Second Duplicate Folder
	flightTape := tape.NewRingBuffer(100)                         // 100 slots max (Flat RAM)
	eventFolder := tape.NewEventFolder(flightTape, 5*time.Second) // 5s folding window

	// 5. Concurrency Engine with 4 Workers & sync.Pool (M3)
	eng := engine.NewEngine(eventsChan, 4, eventFolder)
	eng.Start(ctx)

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	fmt.Printf("[READY] Sentry watching '%s' (Inode #%d)\n", logPath, fileSentry.Inode)
	fmt.Println("[READY] Tape-Recorder: 100 slots | 5s Duplicate Folding Active")
	fmt.Println("----------------------------------------------------------------------------------")

	// 6. Master Event Loop
	for {
		select {
		case <-sigChan:
			fmt.Println("\n\n==================================================================================")
			fmt.Println("[SHUTDOWN] Received SIGINT (Ctrl+C). Flushing workers and tape recorder...")
			cancel()
			eng.Wait() // Workers finish, folder flushed to tape!

			// ASLI FLIGHT RECORDER TAPE DUMP:
			records := flightTape.GetAll()
			fmt.Printf("[FLIGHT TAPE RECORDER DUMP] Captured %d Unique Incidents:\n", len(records))
			for i, r := range records {
				fmt.Printf("  #%02d [%s] Repeated: %dx | Log: \"%s\"\n",
					i+1, r.Timestamp.Format("15:04:05"), r.Count, r.Payload)
			}
			fmt.Println("==================================================================================")
			fmt.Println("[SHUTDOWN] Blackbox stopped gracefully. Bye!")
			return

		case err := <-errChan:
			if err != nil {
				fmt.Printf("[ERROR] Sentry: %v\n", err)
			}

		case <-ticker.C:
			cpuMetrics, _ := cpuCollector.Collect()
			memMetrics, _ := memCollector.Collect()
			diskMetrics, _ := diskCollector.Collect()

			availMB := memMetrics.MemAvailable / (1024 * 1024)
			freeGB := diskMetrics.DiskFree / (1024 * 1024 * 1024)

			processed, bytesTotal, dropped := eng.Metrics().GetStats()

			fmt.Printf("[HEALTH] CPU: %5.2f%% | RAM: %4dMB Avail | DISK: %3dGB Free\n",
				cpuMetrics.CPUPercent, availMB, freeGB)
			fmt.Printf("[ENGINE] Workers: 4 | Processed: %d logs (%d bytes) | Dropped: %d | Tape Slots: %d/100\n",
				processed, bytesTotal, dropped, flightTape.Size())
			fmt.Println("----------------------------------------------------------------------------------")
		}
	}
}
