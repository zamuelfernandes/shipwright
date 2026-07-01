package http

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"

	"github.com/zamuelfernandes/shipwright/internal/domain"
	"github.com/zamuelfernandes/shipwright/internal/usecase"
	"github.com/zamuelfernandes/shipwright/ui"
)

// Router gerencia as rotas HTTP da aplicação, tanto a entrega dos arquivos estáticos da interface
// quanto os endpoints da API REST do Shipwright.
type Router struct {
	mux                *http.ServeMux
	listUseCase        *usecase.ListContainersUseCase
	startUseCase       *usecase.StartContainerUseCase
	stopUseCase        *usecase.StopContainerUseCase
	pruneUseCase       *usecase.PruneContainersUseCase
	streamLogsUseCase  *usecase.StreamLogsUseCase
	streamStatsUseCase *usecase.StreamStatsUseCase
}

// NewRouter recebe os casos de uso necessários por injeção de dependência e configura as rotas.
func NewRouter(
	listUC *usecase.ListContainersUseCase,
	startUC *usecase.StartContainerUseCase,
	stopUC *usecase.StopContainerUseCase,
	pruneUC *usecase.PruneContainersUseCase,
	streamLogsUC *usecase.StreamLogsUseCase,
	streamStatsUC *usecase.StreamStatsUseCase,
) *Router {
	mux := http.NewServeMux()

	r := &Router{
		mux:                mux,
		listUseCase:        listUC,
		startUseCase:       startUC,
		stopUseCase:        stopUC,
		pruneUseCase:       pruneUC,
		streamLogsUseCase:  streamLogsUC,
		streamStatsUseCase: streamStatsUC,
	}

	// API REST
	mux.HandleFunc("GET /api/containers", r.handleListContainers)
	mux.HandleFunc("POST /api/containers/{id}/start", r.handleStartContainer)
	mux.HandleFunc("POST /api/containers/{id}/stop", r.handleStopContainer)
	mux.HandleFunc("DELETE /api/containers/prune", r.handlePruneContainers)

	// API de Streaming SSE (Server-Sent Events)
	mux.HandleFunc("GET /api/containers/{id}/logs", r.handleStreamLogs)
	mux.HandleFunc("GET /api/containers/{id}/stats", r.handleStreamStats)

	// Arquivos estáticos da interface embutida
	distFS, err := fs.Sub(ui.DistFS, "dist")
	if err != nil {
		panic("erro ao extrair sub-filesystem dist da interface: " + err.Error())
	}
	fileServer := http.FileServer(http.FS(distFS))
	mux.Handle("/", fileServer)

	return r
}

// ServeHTTP delega o tratamento HTTP da struct para o 'mux' interno.
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}

// handleListContainers lida com o endpoint GET /api/containers
func (r *Router) handleListContainers(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	containers, err := r.listUseCase.Execute(req.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	json.NewEncoder(w).Encode(containers)
}

// handleStartContainer lida com o endpoint POST /api/containers/{id}/start
func (r *Router) handleStartContainer(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := req.PathValue("id")

	err := r.startUseCase.Execute(req.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "container iniciado com sucesso"})
}

// handleStopContainer lida com o endpoint POST /api/containers/{id}/stop
func (r *Router) handleStopContainer(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := req.PathValue("id")

	err := r.stopUseCase.Execute(req.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "container parado com sucesso"})
}

// handlePruneContainers lida com o endpoint DELETE /api/containers/prune
func (r *Router) handlePruneContainers(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	err := r.pruneUseCase.Execute(req.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "limpeza de containers concluída"})
}

// handleStreamLogs gerencia o streaming SSE de logs do container
func (r *Router) handleStreamLogs(w http.ResponseWriter, req *http.Request) {
	// Configura cabeçalhos necessários para manter a conexão aberta (SSE)
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	id := req.PathValue("id")

	// A interface Flusher permite enviar dados imediatamente para o cliente
	// sem esperar que a requisição seja concluída.
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming não suportado", http.StatusInternalServerError)
		return
	}

	logsChan := make(chan string, 10)
	errChan := make(chan error, 1)

	// Usamos o contexto da requisição. Se o usuário fechar a aba/conexão no navegador,
	// req.Context().Done() é ativado e nós cancelamos a goroutine secundária.
	ctx, cancel := context.WithCancel(req.Context())
	defer cancel()

	// Dispara o caso de uso em uma Goroutine paralela (semelhante a criar um isolate ou stream no Dart)
	go r.streamLogsUseCase.Execute(ctx, id, logsChan, errChan)

	// Envia um evento inicial de conexão estabelecida
	fmt.Fprint(w, "event: open\ndata: conectado aos logs\n\n")
	flusher.Flush()

	for {
		select {
		case <-ctx.Done():
			return
		case err := <-errChan:
			fmt.Fprintf(w, "event: error\ndata: %s\n\n", err.Error())
			flusher.Flush()
			return
		case logLine, ok := <-logsChan:
			if !ok {
				return
			}
			// Formato padrão de mensagens SSE: data: <valor>\n\n
			fmt.Fprintf(w, "data: %s\n\n", logLine)
			flusher.Flush()
		}
	}
}

// handleStreamStats gerencia o streaming SSE de estatísticas de CPU/RAM do container
func (r *Router) handleStreamStats(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	id := req.PathValue("id")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming não suportado", http.StatusInternalServerError)
		return
	}

	statsChan := make(chan domain.ContainerStats, 10)
	errChan := make(chan error, 1)

	ctx, cancel := context.WithCancel(req.Context())
	defer cancel()

	go r.streamStatsUseCase.Execute(ctx, id, statsChan, errChan)

	fmt.Fprint(w, "event: open\ndata: conectado à telemetria\n\n")
	flusher.Flush()

	for {
		select {
		case <-ctx.Done():
			return
		case err := <-errChan:
			fmt.Fprintf(w, "event: error\ndata: %s\n\n", err.Error())
			flusher.Flush()
			return
		case stats, ok := <-statsChan:
			if !ok {
				return
			}
			data, err := json.Marshal(stats)
			if err != nil {
				return
			}
			fmt.Fprintf(w, "data: %s\n\n", string(data))
			flusher.Flush()
		}
	}
}
