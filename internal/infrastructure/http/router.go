package http

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/zamuelfernandes/shipwright/internal/domain"
	"github.com/zamuelfernandes/shipwright/internal/usecase"
	"github.com/zamuelfernandes/shipwright/ui"
)

// Upgrader para configurar conexões WebSocket
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Permite qualquer origem em ambiente de desenvolvimento local
	},
}

// Router gerencia as rotas HTTP da aplicação, endpoints REST e conexões WebSocket.
type Router struct {
	mux                 *http.ServeMux
	listUseCase         *usecase.ListContainersUseCase
	startUseCase        *usecase.StartContainerUseCase
	stopUseCase         *usecase.StopContainerUseCase
	pruneUseCase        *usecase.PruneContainersUseCase
	streamLogsUseCase   *usecase.StreamLogsUseCase
	streamStatsUseCase  *usecase.StreamStatsUseCase
	startProjectUseCase *usecase.StartProjectUseCase
	stopProjectUseCase  *usecase.StopProjectUseCase
	execUseCase         *usecase.ExecContainerUseCase
}

// NewRouter recebe os casos de uso necessários por injeção de dependência e configura as rotas.
func NewRouter(
	listUC *usecase.ListContainersUseCase,
	startUC *usecase.StartContainerUseCase,
	stopUC *usecase.StopContainerUseCase,
	pruneUC *usecase.PruneContainersUseCase,
	streamLogsUC *usecase.StreamLogsUseCase,
	streamStatsUC *usecase.StreamStatsUseCase,
	startProjUC *usecase.StartProjectUseCase,
	stopProjUC *usecase.StopProjectUseCase,
	execUC *usecase.ExecContainerUseCase,
) *Router {
	mux := http.NewServeMux()

	r := &Router{
		mux:                 mux,
		listUseCase:         listUC,
		startUseCase:        startUC,
		stopUseCase:         stopUC,
		pruneUseCase:        pruneUC,
		streamLogsUseCase:   streamLogsUC,
		streamStatsUseCase:  streamStatsUC,
		startProjectUseCase: startProjUC,
		stopProjectUseCase:  stopProjUC,
		execUseCase:         execUC,
	}

	// API REST - Containers individuais
	mux.HandleFunc("GET /api/containers", r.handleListContainers)
	mux.HandleFunc("POST /api/containers/{id}/start", r.handleStartContainer)
	mux.HandleFunc("POST /api/containers/{id}/stop", r.handleStopContainer)
	mux.HandleFunc("DELETE /api/containers/prune", r.handlePruneContainers)

	// API REST - Projetos Docker Compose (V2)
	mux.HandleFunc("POST /api/projects/{name}/start", r.handleStartProject)
	mux.HandleFunc("POST /api/projects/{name}/stop", r.handleStopProject)

	// API de Streaming SSE (Server-Sent Events)
	mux.HandleFunc("GET /api/containers/{id}/logs", r.handleStreamLogs)
	mux.HandleFunc("GET /api/containers/{id}/stats", r.handleStreamStats)

	// Terminal Interativo WebSocket (V2)
	mux.HandleFunc("GET /api/containers/{id}/exec", r.handleExecContainer)

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

// handleStartProject lida com o endpoint POST /api/projects/{name}/start (V2)
func (r *Router) handleStartProject(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	name := req.PathValue("name")

	err := r.startProjectUseCase.Execute(req.Context(), name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "projeto inicializado com sucesso"})
}

// handleStopProject lida com o endpoint POST /api/projects/{name}/stop (V2)
func (r *Router) handleStopProject(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	name := req.PathValue("name")

	err := r.stopProjectUseCase.Execute(req.Context(), name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "projeto pausado com sucesso"})
}

// handleStreamLogs gerencia o streaming SSE de logs do container
func (r *Router) handleStreamLogs(w http.ResponseWriter, req *http.Request) {
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

	logsChan := make(chan string, 10)
	errChan := make(chan error, 1)

	ctx, cancel := context.WithCancel(req.Context())
	defer cancel()

	go r.streamLogsUseCase.Execute(ctx, id, logsChan, errChan)

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
			fmt.Fprintf(w, "data: %s\n\n", logLine)
			flusher.Flush()
		}
	}
}

// handleStreamStats gerencia o streaming SSE de estatísticas de CPU/RAM
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

// handleExecContainer lida com a conexão WebSocket para o terminal interativo (V2)
func (r *Router) handleExecContainer(w http.ResponseWriter, req *http.Request) {
	id := req.PathValue("id")
	log.Printf("[HTTP] Recebido pedido de terminal WebSocket para o container: %s", id)

	// Upgrade da conexão HTTP comum para WebSocket
	conn, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Printf("[HTTP] Erro no upgrade do websocket para exec: %v", err)
		return
	}
	defer conn.Close()
	log.Printf("[HTTP] Conexão WebSocket de terminal ativa para o container %s", id)

	// io.Pipe liga a escrita do WebSocket à leitura do stdin do Exec do Docker
	stdinReader, stdinWriter := io.Pipe()

	// Goroutine para ler teclas do WebSocket -> escreve no Pipe (stdinWriter)
	go func() {
		defer stdinWriter.Close()
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				break
			}
			_, err = stdinWriter.Write(msg)
			if err != nil {
				break
			}
		}
	}()

	// Adaptador io.Writer customizado para escrever de volta no WebSocket
	wsWriter := &wsWriterWrapper{conn: conn}

	// Executa a sessão interativa dentro do container. Bloqueia até a sessão encerrar.
	err = r.execUseCase.Execute(req.Context(), id, stdinReader, wsWriter)
	if err != nil {
		log.Printf("[HTTP] Erro na execução do terminal no container %s: %v", id, err)
		_ = conn.WriteMessage(websocket.TextMessage, []byte("\r\nErro ao iniciar terminal: "+err.Error()+"\r\n"))
	} else {
		log.Printf("[HTTP] Conexão WebSocket de terminal encerrada para o container %s", id)
	}
}

// wsWriterWrapper adapta a conexão WebSocket para que satisfaça a interface io.Writer
type wsWriterWrapper struct {
	conn *websocket.Conn
}

func (w *wsWriterWrapper) Write(p []byte) (n int, err error) {
	// Escreve os bytes recebidos da saída do terminal do container para o WebSocket como texto
	err = w.conn.WriteMessage(websocket.TextMessage, p)
	if err != nil {
		return 0, err
	}
	return len(p), nil
}
