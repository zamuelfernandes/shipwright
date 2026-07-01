package usecase

import (
	"context"
	"github.com/zamuelfernandes/shipwright/internal/domain"
)

// PruneContainersUseCase coordena a ação de remover todos os containers parados.
type PruneContainersUseCase struct {
	repo domain.ContainerRepository
}

func NewPruneContainersUseCase(repo domain.ContainerRepository) *PruneContainersUseCase {
	return &PruneContainersUseCase{repo: repo}
}

func (u *PruneContainersUseCase) Execute(ctx context.Context) error {
	return u.repo.PruneContainers(ctx)
}
