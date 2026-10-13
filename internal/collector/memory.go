package collector

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type MemoryCollector struct {
	memInfoPath string
}

func NewMemoryCollector(memInfoPath string) *MemoryCollector {
	return &MemoryCollector{
		memInfoPath: memInfoPath,
	}
}

func (m *MemoryCollector) Collect() (*MemoryMetrics, error) {
	file, err := os.Open(m.memInfoPath)

	if err !=nil {
		return nil, fmt.Errorf("failed to open mem file: %w", err)
	}

	defer file.Close()

	var metrics MemoryMetrics

	scanner := bufio.NewScanner(file)

	for scanner.Scan(){
		line := scanner.Text()

		if strings.HasPrefix(line, "MemTotal:"){
			var val uint64
			_, err := fmt.Sscanf(line, "MemTotal: %d kB", &val)
			if err == nil {
				metrics.MemTotal = val * 1024
			}
		}else if strings.HasPrefix(line, "MemFree:"){
			var val uint64
			_,err := fmt.Sscanf(line, "MemFree: %d kB", &val)

			if err == nil {
				metrics.MemFree = val * 1024
			}
		}else if strings.HasPrefix(line, "MemAvailable:"){
			var val uint64
			_, err := fmt.Sscanf(line, "MemAvailable: %d kB", &val)
			if err == nil {
				metrics.MemAvailable = val * 1024
			}
		}

		if metrics.MemTotal > 0 && metrics.MemFree > 0 && metrics.MemAvailable > 0 {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed something mem file: %w", err)
	}

	return &metrics, nil
}