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
	fmt.Println("=== Shipwright (Fase 4: Streaming de Logs & Stats em Tempo Real) ===")

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
	
	// Casos de Uso da Fase 4 para Streaming SSE
	streamLogsUC := usecase.NewStreamLogsUseCase(dockerClient)
	streamStatsUC := usecase.NewStreamStatsUseCase(dockerClient)

	// Teste rápido de conexão para garantir que o socket está acessível
	containers, err := listUC.Execute(context.Background())
	if err != nil {
		log.Printf("Aviso: Falha ao conectar ao Docker: %v", err)
	} else {
		log.Printf("Sucesso: Conectado ao Docker Daemon local (%d containers encontrados).", len(containers))
	}

	// 3. Inicialização do roteador HTTP injetando todas as dependências de regras de negócio
	router := infraHttp.NewRouter(listUC, startUC, stopUC, pruneUC, streamLogsUC, streamStatsUC)

	// 4. Inicialização do Servidor HTTP na porta local 8080
	addr := ":8080"
	fmt.Printf("\n⚓ Shipwright rodando em: http://localhost%s\n", addr)
	fmt.Println("Pressione Ctrl+C para encerrar o servidor.")

	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("Erro fatal no servidor HTTP: %v", err)
	}
}
