package usecase

import (
	"context"
	"github.com/zamuelfernandes/shipwright/internal/domain"
)

// StreamLogsUseCase coordena a transmissão contínua de logs.
type StreamLogsUseCase struct {
	repo domain.ContainerRepository
}

func NewStreamLogsUseCase(repo domain.ContainerRepository) *StreamLogsUseCase {
	return &StreamLogsUseCase{repo: repo}
}

func (u *StreamLogsUseCase) Execute(ctx context.Context, id string, logsChan chan<- string, errChan chan<- error) {
	u.repo.StreamLogs(ctx, id, logsChan, errChan)
}
