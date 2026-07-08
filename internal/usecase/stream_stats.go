package usecase

import (
	"context"
	"github.com/zamuelfernandes/shipwright/internal/domain"
)

// StreamStatsUseCase coordena a transmissão contínua de telemetria (CPU/RAM).
type StreamStatsUseCase struct {
	repo domain.TelemetryRepository
}

func NewStreamStatsUseCase(repo domain.TelemetryRepository) *StreamStatsUseCase {
	return &StreamStatsUseCase{repo: repo}
}

func (u *StreamStatsUseCase) Execute(ctx context.Context, id string, statsChan chan<- domain.ContainerStats, errChan chan<- error) {
	u.repo.StreamStats(ctx, id, statsChan, errChan)
}
