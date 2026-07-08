package domain

import "context"

// Container represents simplified Docker container metadata.
type Container struct {
	ID          string   `json:"id"`
	Names       []string `json:"names"`
	Image       string   `json:"image"`
	State       string   `json:"state"`
	Status      string   `json:"status"`
	ComposeProj string   `json:"compose_project"`
	ComposeServ string   `json:"compose_service"`
}

// ContainerRepository defines the contract for container lifecycle operations.
type ContainerRepository interface {
	ListContainers(ctx context.Context) ([]Container, error)
	StartContainer(ctx context.Context, id string) error
	StopContainer(ctx context.Context, id string) error
	PruneContainers(ctx context.Context) error
}
