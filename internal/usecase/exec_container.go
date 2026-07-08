package usecase

import (
	"context"
	"io"

	"github.com/zamuelfernandes/anchordock/internal/domain"
)

// ExecContainerUseCase coordinates running an interactive shell session inside a container.
type ExecContainerUseCase struct {
	repo domain.ExecRepository
}

func NewExecContainerUseCase(repo domain.ExecRepository) *ExecContainerUseCase {
	return &ExecContainerUseCase{repo: repo}
}

func (u *ExecContainerUseCase) Execute(ctx context.Context, id string, stdin io.Reader, stdout io.Writer) error {
	return u.repo.ExecContainer(ctx, id, stdin, stdout)
}
