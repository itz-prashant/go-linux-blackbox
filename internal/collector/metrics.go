package collector

import "time"

type CPUMetrics struct{
	CPUPercent float64 `json:"cpu_percent"`
	IOWaitPercent float64 `json:"io_wait_percent"`
}

type MemoryMetrics struct {
	MemTotal uint64 `json:"mem_total"`
	MemFree uint64 `json:"mem_free"`
	MemAvailable uint64 `json:"mem_available"`
}

type DiskMetrics struct {
	DiskTotal uint64 `json:"disk_total"`
	DiskFree uint64 `json:"disk_free"`
	DiskUsedPercent float64 `json:"disk_used_percent"`
}

type SystemSnapshot struct {
	CPUMetrics	`json:"cpu"`
	MemoryMetrics `json:"memory"`
	DiskMetrics	`json:"disk"`
	TimeStamp time.Time `json:"time_stamp"`
}
