package domain

import (
	"context"
	"io"
)

// ExecRepository defines the contract for running interactive commands inside containers.
type ExecRepository interface {
	ExecContainer(ctx context.Context, id string, stdin io.Reader, stdout io.Writer) error
}
