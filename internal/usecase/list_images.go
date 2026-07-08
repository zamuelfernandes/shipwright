package usecase

import (
	"context"
	"github.com/zamuelfernandes/shipwright/internal/domain"
)

// ListImagesUseCase coordena a listagem de imagens do Docker.
type ListImagesUseCase struct {
	repo domain.ImageRepository
}

func NewListImagesUseCase(repo domain.ImageRepository) *ListImagesUseCase {
	return &ListImagesUseCase{repo: repo}
}

func (u *ListImagesUseCase) Execute(ctx context.Context) ([]domain.Image, error) {
	return u.repo.ListImages(ctx)
}
