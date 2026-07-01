package domain

import "context"

// Container representa os dados simplificados de um container Docker para o Shipwright.
// Em Flutter, isso seria um Modelo de Dados (Model ou Entity) contendo campos e possivelmente serialização JSON.
type Container struct {
	ID     string   `json:"id"`
	Names  []string `json:"names"`
	Image  string   `json:"image"`
	State  string   `json:"state"`
	Status string   `json:"status"`
}

// ContainerRepository define o contrato (interface) para manipulação de dados dos containers.
// Em Flutter/Dart, isso é exatamente igual a uma 'abstract class ContainerRepository'.
// No Go, usamos interfaces. A principal diferença é que no Go a implementação é implícita:
// qualquer struct que implemente os métodos definidos aqui satisfará a interface automaticamente,
// sem precisar de palavras-chave como 'implements'.
type ContainerRepository interface {
	ListContainers(ctx context.Context) ([]Container, error)
	StartContainer(ctx context.Context, id string) error
	StopContainer(ctx context.Context, id string) error
	PruneContainers(ctx context.Context) error
	StreamLogs(ctx context.Context, id string, logsChan chan<- string, errChan chan<- error)
	StreamStats(ctx context.Context, id string, statsChan chan<- ContainerStats, errChan chan<- error)
}
