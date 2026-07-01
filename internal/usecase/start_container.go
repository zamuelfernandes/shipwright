package usecase

import (
	"context"
	"github.com/zamuelfernandes/shipwright/internal/domain"
)

// StartContainerUseCase coordena a ação de iniciar um container.
// Em Flutter, isso seria o caso de uso 'StartContainerUseCase'.
type StartContainerUseCase struct {
	repo domain.ContainerRepository
}

func NewStartContainerUseCase(repo domain.ContainerRepository) *StartContainerUseCase {
	return &StartContainerUseCase{repo: repo}
}

func (u *StartContainerUseCase) Execute(ctx context.Context, id string) error {
	return u.repo.StartContainer(ctx, id)
}
