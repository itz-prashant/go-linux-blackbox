package collector

import (
	"fmt"
	"syscall"
)

type DiskCollector struct {
	mountPath string
}

func NewDiskCollector(mountPath string) *DiskCollector {
	return &DiskCollector{
		mountPath: mountPath,
	}
}

func (d *DiskCollector) Collect() (*DiskMetrics, error) {
	var stat syscall.Statfs_t

	err := syscall.Statfs(d.mountPath, &stat)

	if err != nil {
		return nil, fmt.Errorf("failed to get disk status: %w", err)
	}

	diskTotal := stat.Blocks * uint64(stat.Bsize)
	diskFree := stat.Bavail * uint64(stat.Bsize)
	diskUsed := diskTotal - diskFree

	var diskUsedPercentage float64

	if diskTotal > 0 {
		diskUsedPercentage = (float64(diskUsed) / float64(diskTotal)) * 100.0
	}

	return &DiskMetrics{DiskTotal: diskTotal, DiskFree: diskFree, DiskUsedPercent: diskUsedPercentage}, nil
}
