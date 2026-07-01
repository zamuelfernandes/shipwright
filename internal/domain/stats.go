package domain

// ContainerStats representa a telemetria (uso de CPU e memória RAM) simplificada de um container.
type ContainerStats struct {
	CPUPercent       float64 `json:"cpu_percent"`
	MemoryUsageBytes uint64  `json:"memory_usage_bytes"`
	MemoryLimitBytes uint64  `json:"memory_limit_bytes"`
	MemoryPercent    float64 `json:"memory_percent"`
}
