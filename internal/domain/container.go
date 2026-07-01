package domain

import (
	"context"
	"io"
)

// Container representa os dados simplificados de um container Docker para o Shipwright.
// Em Flutter, isso seria um Modelo de Dados (Model ou Entity) contendo campos e possivelmente serialização JSON.
type Container struct {
	ID          string   `json:"id"`
	Names       []string `json:"names"`
	Image       string   `json:"image"`
	State       string   `json:"state"`
	Status      string   `json:"status"`
	ComposeProj string   `json:"compose_project"`
	ComposeServ string   `json:"compose_service"`
}

// ContainerRepository define o contrato (interface) para manipulação de dados dos containers.
type ContainerRepository interface {
	ListContainers(ctx context.Context) ([]Container, error)
	StartContainer(ctx context.Context, id string) error
	StopContainer(ctx context.Context, id string) error
	PruneContainers(ctx context.Context) error
	StreamLogs(ctx context.Context, id string, logsChan chan<- string, errChan chan<- error)
	StreamStats(ctx context.Context, id string, statsChan chan<- ContainerStats, errChan chan<- error)

	// Métodos da V2
	StartProject(ctx context.Context, project string) error
	StopProject(ctx context.Context, project string) error
	ExecContainer(ctx context.Context, id string, stdin io.Reader, stdout io.Writer) error
}
