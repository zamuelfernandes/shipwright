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

func main() {
	fmt.Println("=== Shipwright (Fase V2: Compose Projects & Terminal Interativo) ===")

	loadEnv(".env")

	// 1. Initialize Infrastructure (Docker Client connection)
	dockerClient, err := docker.NewDockerClient()
	if err != nil {
		log.Fatalf("Error connecting to Docker Client: %v", err)
	}

	// 2. Instantiate Usecases
	listUC := usecase.NewListContainersUseCase(dockerClient)
	startUC := usecase.NewStartContainerUseCase(dockerClient)
	stopUC := usecase.NewStopContainerUseCase(dockerClient)
	pruneUC := usecase.NewPruneContainersUseCase(dockerClient)
	
	streamLogsUC := usecase.NewStreamLogsUseCase(dockerClient)
	streamStatsUC := usecase.NewStreamStatsUseCase(dockerClient)

	startProjUC := usecase.NewStartProjectUseCase(dockerClient)
	stopProjUC := usecase.NewStopProjectUseCase(dockerClient)
	execUC := usecase.NewExecContainerUseCase(dockerClient)
	listImagesUC := usecase.NewListImagesUseCase(dockerClient)

	// Quick connection test to verify Docker Socket accessibility
	containers, err := listUC.Execute(context.Background())
	if err != nil {
		log.Printf("Warning: Failed to connect to Docker daemon: %v", err)
	} else {
		log.Printf("Success: Connected to local Docker daemon (%d containers detected).", len(containers))
	}

	// 3. Initialize HTTP resource handlers (SRP Controllers)
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

	// 4. Initialize HTTP multiplexer router
	router := infraHttp.NewRouter(
		containerHandler,
		projectHandler,
		imageHandler,
	)

	// 5. Start HTTP server on the configured port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := ":" + port
	fmt.Printf("\n⚓ Shipwright running at: http://localhost%s\n", addr)
	fmt.Println("Press Ctrl+C to stop the server.")

	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("Fatal HTTP Server error: %v", err)
	}
}

// loadEnv reads environment variables from a local .env file.
func loadEnv(filepath string) {
	file, err := os.Open(filepath)
	if err != nil {
		return
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
			if os.Getenv(key) == "" {
				os.Setenv(key, value)
			}
		}
	}
}
