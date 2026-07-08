package usecase

import (
	"context"
	"github.com/zamuelfernandes/shipwright/internal/domain"
)

// StartProjectUseCase coordena a inicialização de todos os containers de um projeto do Compose.
type StartProjectUseCase struct {
	repo domain.ProjectRepository
}

func NewStartProjectUseCase(repo domain.ProjectRepository) *StartProjectUseCase {
	return &StartProjectUseCase{repo: repo}
}

func (u *StartProjectUseCase) Execute(ctx context.Context, project string) error {
	return u.repo.StartProject(ctx, project)
}
