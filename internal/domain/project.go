package domain

import "context"

// ProjectRepository define o contrato para controle de lote (start/stop) de projetos Docker Compose.
type ProjectRepository interface {
	StartProject(ctx context.Context, project string) error
	StopProject(ctx context.Context, project string) error
}
