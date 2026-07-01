package usecase

import (
	"context"
	"github.com/zamuelfernandes/shipwright/internal/domain"
)

// ListContainersUseCase representa o caso de uso para listar containers.
// Em Flutter/Dart com Clean Architecture, isso seria uma classe 'ListContainersUseCase'
// contendo uma dependência de 'ContainerRepository' injetada no construtor.
type ListContainersUseCase struct {
	repo domain.ContainerRepository
}

// NewListContainersUseCase funciona como um "construtor".
// No Dart/Flutter, você usaria 'ListContainersUseCase(this.repo)'.
// Em Go, não temos classes com construtores, então é padrão criar uma função externa
// iniciada com "New" que inicializa e retorna um ponteiro para a struct.
func NewListContainersUseCase(repo domain.ContainerRepository) *ListContainersUseCase {
	return &ListContainersUseCase{repo: repo}
}

// Execute roda a lógica de negócio do caso de uso.
// Em Flutter/Dart, poderíamos usar um método 'call()' para que a classe seja executável
// diretamente como uma função: 'useCase()'.
// No Go, costuma-se criar um método explícito como 'Execute' ou 'Run'.
func (u *ListContainersUseCase) Execute(ctx context.Context) ([]domain.Container, error) {
	return u.repo.ListContainers(ctx)
}
