package usecase

import (
	"context"
	"github.com/zamuelfernandes/anchordock/internal/domain"
)

// StopContainerUseCase coordinates stopping a single container.
type StopContainerUseCase struct {
	repo domain.ContainerRepository
}

func NewStopContainerUseCase(repo domain.ContainerRepository) *StopContainerUseCase {
	return &StopContainerUseCase{repo: repo}
}

func (u *StopContainerUseCase) Execute(ctx context.Context, id string) error {
	return u.repo.StopContainer(ctx, id)
}
