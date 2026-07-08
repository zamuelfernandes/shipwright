package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/zamuelfernandes/shipwright/internal/infrastructure/docker"
	infraHttp "github.com/zamuelfernandes/shipwright/internal/infrastructure/http"
	"github.com/zamuelfernandes/shipwright/internal/usecase"
)

// main é o ponto de entrada da aplicação.
func main() {
	fmt.Println("=== Shipwright (Fase V2: Compose Projects & Terminal Interativo) ===")

	// Carrega arquivo .env local se existir
	loadEnv(".env")

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

	// 5. Inicialização do Servidor HTTP na porta local configurada
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := ":" + port
	fmt.Printf("\n⚓ Shipwright rodando em: http://localhost%s\n", addr)
	fmt.Println("Pressione Ctrl+C para encerrar o servidor.")

	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("Erro fatal no servidor HTTP: %v", err)
	}
}

// loadEnv lê um arquivo .env se ele existir e carrega as variáveis no ambiente.
func loadEnv(filepath string) {
	file, err := os.Open(filepath)
	if err != nil {
		return // Se não existir, apenas ignora e usa as variáveis de ambiente existentes
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			// Se a variável já estiver definida no ambiente do OS, não sobrescreve
			if os.Getenv(key) == "" {
				os.Setenv(key, value)
			}
		}
	}
}
