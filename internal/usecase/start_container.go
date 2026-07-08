package usecase

import (
	"context"
	"github.com/zamuelfernandes/anchordock/internal/domain"
)

// StartContainerUseCase coordinates starting a single container.
type StartContainerUseCase struct {
	repo domain.ContainerRepository
}

func NewStartContainerUseCase(repo domain.ContainerRepository) *StartContainerUseCase {
	return &StartContainerUseCase{repo: repo}
}

func (u *StartContainerUseCase) Execute(ctx context.Context, id string) error {
	return u.repo.StartContainer(ctx, id)
}
