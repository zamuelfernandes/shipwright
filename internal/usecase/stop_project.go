package usecase

import (
	"context"
	"github.com/zamuelfernandes/shipwright/internal/domain"
)

// StopProjectUseCase coordena a parada de todos os containers de um projeto do Compose.
type StopProjectUseCase struct {
	repo domain.ProjectRepository
}

func NewStopProjectUseCase(repo domain.ProjectRepository) *StopProjectUseCase {
	return &StopProjectUseCase{repo: repo}
}

func (u *StopProjectUseCase) Execute(ctx context.Context, project string) error {
	return u.repo.StopProject(ctx, project)
}
