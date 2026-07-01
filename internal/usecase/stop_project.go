package usecase

import (
	"context"
	"github.com/zamuelfernandes/shipwright/internal/domain"
)

// StopProjectUseCase coordena a parada de todos os containers de um projeto do Compose.
type StopProjectUseCase struct {
	repo domain.ContainerRepository
}

func NewStopProjectUseCase(repo domain.ContainerRepository) *StopProjectUseCase {
	return &StopProjectUseCase{repo: repo}
}

func (u *StopProjectUseCase) Execute(ctx context.Context, project string) error {
	return u.repo.StopProject(ctx, project)
}
