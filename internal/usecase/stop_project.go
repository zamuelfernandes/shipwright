package usecase

import (
	"context"
	"github.com/zamuelfernandes/anchordock/internal/domain"
)

// StopProjectUseCase coordinates batch stop commands of Compose projects.
type StopProjectUseCase struct {
	repo domain.ProjectRepository
}

func NewStopProjectUseCase(repo domain.ProjectRepository) *StopProjectUseCase {
	return &StopProjectUseCase{repo: repo}
}

func (u *StopProjectUseCase) Execute(ctx context.Context, project string) error {
	return u.repo.StopProject(ctx, project)
}
