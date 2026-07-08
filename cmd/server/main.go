package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/zamuelfernandes/shipwright/internal/infrastructure/docker"
	infraHttp "github.com/zamuelfernandes/shipwright/internal/infrastructure/http"
	"github.com/zamuelfernandes/shipwright/internal/usecase"
)

// main é o ponto de entrada da aplicação.
func main() {
	fmt.Println("=== Shipwright (Fase V2: Compose Projects & Terminal Interativo) ===")

	// 1. Inicialização da Infraestrutura (Docker Client)
	dockerClient, err := docker.NewDockerClient()
	if err != nil {
		log.Fatalf("Erro ao inicializar cliente Docker: %v", err)
	}

	// 2. Instanciação de todos os Casos de Uso
	listUC := usecase.NewListContainersUseCase(dockerClient)
	startUC := usecase.NewStartContainerUseCase(dockerClient)
	stopUC := usecase.NewStopContainerUseCase(dockerClient)
	pruneUC := usecase.NewPruneContainersUseCase(dockerClient)
	
	// Streaming SSE
	streamLogsUC := usecase.NewStreamLogsUseCase(dockerClient)
	streamStatsUC := usecase.NewStreamStatsUseCase(dockerClient)

	// V2: Compose e Terminal Exec
	startProjUC := usecase.NewStartProjectUseCase(dockerClient)
	stopProjUC := usecase.NewStopProjectUseCase(dockerClient)
	execUC := usecase.NewExecContainerUseCase(dockerClient)
	listImagesUC := usecase.NewListImagesUseCase(dockerClient)

	// Teste rápido de conexão para garantir que o socket está acessível
	containers, err := listUC.Execute(context.Background())
	if err != nil {
		log.Printf("Aviso: Falha ao conectar ao Docker: %v", err)
	} else {
		log.Printf("Sucesso: Conectado ao Docker Daemon local (%d containers encontrados).", len(containers))
	}

	// 3. Inicialização dos Handlers HTTP segregados (SRP)
	containerHandler := infraHttp.NewContainerHandler(
		listUC,
		startUC,
		stopUC,
		pruneUC,
		streamLogsUC,
		streamStatsUC,
		execUC,
	)
	projectHandler := infraHttp.NewProjectHandler(
		startProjUC,
		stopProjUC,
	)
	imageHandler := infraHttp.NewImageHandler(
		listImagesUC,
	)

	// 4. Inicialização do roteador HTTP injetando os handlers
	router := infraHttp.NewRouter(
		containerHandler,
		projectHandler,
		imageHandler,
	)

	// 4. Inicialização do Servidor HTTP na porta local 8080
	addr := ":8080"
	fmt.Printf("\n⚓ Shipwright rodando em: http://localhost%s\n", addr)
	fmt.Println("Pressione Ctrl+C para encerrar o servidor.")

	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("Erro fatal no servidor HTTP: %v", err)
	}
}
