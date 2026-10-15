package snapshot

import (
	"encoding/json"
	"time"

	"github.com/itz-prashant/linux-blackbox/internal/tape"
)

type CrimeSceneReport struct {
	Timestamp         time.Time `json:"timestamp"`
	ZScore            float64   `json:"z_score"`
	TriggerReason     string    `json:"trigger_reason"`
	RecentCPUHistory  []float64 `json:"recent_cpu_history"` // Columnar table ke pichle samples
	TopErrorLog       string    `json:"top_error_log"`      // Flight Tape ka sabse bada error
	RepeatCount       uint32    `json:"repeat_count"`       // Wo error kitni baar repeat hua tha
	MemoryAvailableMB uint64    `json:"memory_available_mb"`
}

func GenerateSnapshot(
	zScore float64,
	reason string,
	cpuHistory []float64,
	latestRecord tape.Record,
	availRAM uint64) ([]byte, error) {
	report := CrimeSceneReport{
		Timestamp:         time.Now().UTC(),
		ZScore:            zScore,
		TriggerReason:     reason,
		RecentCPUHistory:  cpuHistory,
		TopErrorLog:       latestRecord.Payload,
		RepeatCount:       latestRecord.Count,
		MemoryAvailableMB: availRAM / (1024 * 1024),
	}
	return json.MarshalIndent(report, "", " ")
}
