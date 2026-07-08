package usecase

import (
	"context"
	"github.com/zamuelfernandes/shipwright/internal/domain"
)

// StreamLogsUseCase coordinates logging streams.
type StreamLogsUseCase struct {
	repo domain.TelemetryRepository
}

func NewStreamLogsUseCase(repo domain.TelemetryRepository) *StreamLogsUseCase {
	return &StreamLogsUseCase{repo: repo}
}

func (u *StreamLogsUseCase) Execute(ctx context.Context, id string, logsChan chan<- string, errChan chan<- error) {
	u.repo.StreamLogs(ctx, id, logsChan, errChan)
}
