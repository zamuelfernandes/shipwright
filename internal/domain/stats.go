package domain

import "context"

// ContainerStats represents simplified real-time telemetry (CPU/RAM) of a container.
type ContainerStats struct {
	CPUPercent       float64 `json:"cpu_percent"`
	MemoryUsageBytes uint64  `json:"memory_usage_bytes"`
	MemoryLimitBytes uint64  `json:"memory_limit_bytes"`
	MemoryPercent    float64 `json:"memory_percent"`
}

// TelemetryRepository defines the contract for container logs and telemetry streaming.
type TelemetryRepository interface {
	StreamLogs(ctx context.Context, id string, logsChan chan<- string, errChan chan<- error)
	StreamStats(ctx context.Context, id string, statsChan chan<- ContainerStats, errChan chan<- error)
}
