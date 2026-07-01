package usecase

import (
	"context"
	"github.com/zamuelfernandes/shipwright/internal/domain"
)

// StopContainerUseCase coordena a ação de parar um container.
type StopContainerUseCase struct {
	repo domain.ContainerRepository
}

func NewStopContainerUseCase(repo domain.ContainerRepository) *StopContainerUseCase {
	return &StopContainerUseCase{repo: repo}
}

func (u *StopContainerUseCase) Execute(ctx context.Context, id string) error {
	return u.repo.StopContainer(ctx, id)
}
