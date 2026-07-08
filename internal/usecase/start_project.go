package usecase

import (
	"context"
	"github.com/zamuelfernandes/anchordock/internal/domain"
)

// StartProjectUseCase coordinates batch start commands of Compose projects.
type StartProjectUseCase struct {
	repo domain.ProjectRepository
}

func NewStartProjectUseCase(repo domain.ProjectRepository) *StartProjectUseCase {
	return &StartProjectUseCase{repo: repo}
}

func (u *StartProjectUseCase) Execute(ctx context.Context, project string) error {
	return u.repo.StartProject(ctx, project)
}
