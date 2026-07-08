package usecase

import (
	"context"
	"github.com/zamuelfernandes/anchordock/internal/domain"
)

// ListContainersUseCase coordinates listing all containers.
type ListContainersUseCase struct {
	repo domain.ContainerRepository
}

func NewListContainersUseCase(repo domain.ContainerRepository) *ListContainersUseCase {
	return &ListContainersUseCase{repo: repo}
}

func (u *ListContainersUseCase) Execute(ctx context.Context) ([]domain.Container, error) {
	return u.repo.ListContainers(ctx)
}
