package usecase

import (
	"context"
	"io"

	"github.com/zamuelfernandes/shipwright/internal/domain"
)

// ExecContainerUseCase coordena a execução de um shell interativo no container.
type ExecContainerUseCase struct {
	repo domain.ContainerRepository
}

func NewExecContainerUseCase(repo domain.ContainerRepository) *ExecContainerUseCase {
	return &ExecContainerUseCase{repo: repo}
}

func (u *ExecContainerUseCase) Execute(ctx context.Context, id string, stdin io.Reader, stdout io.Writer) error {
	return u.repo.ExecContainer(ctx, id, stdin, stdout)
}
