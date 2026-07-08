package domain

import (
	"context"
	"io"
)

// ExecRepository define o contrato para execução de comandos interativos em containers.
type ExecRepository interface {
	ExecContainer(ctx context.Context, id string, stdin io.Reader, stdout io.Writer) error
}
