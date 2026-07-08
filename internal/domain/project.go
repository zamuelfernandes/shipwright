package domain

import "context"

// ProjectRepository defines the contract for batch start/stop control of Docker Compose projects.
type ProjectRepository interface {
	StartProject(ctx context.Context, project string) error
	StopProject(ctx context.Context, project string) error
}
