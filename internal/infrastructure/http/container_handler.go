package http

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/zamuelfernandes/shipwright/internal/domain"
	"github.com/zamuelfernandes/shipwright/internal/usecase"
)

// Upgrader para configurar conexões WebSocket
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Permite qualquer origem em ambiente de desenvolvimento local
	},
}

// ContainerHandler gerencia as rotas HTTP e WebSocket associadas a containers.
type ContainerHandler struct {
	listUseCase        *usecase.ListContainersUseCase
	startUseCase       *usecase.StartContainerUseCase
	stopUseCase        *usecase.StopContainerUseCase
	pruneUseCase       *usecase.PruneContainersUseCase
	streamLogsUseCase  *usecase.StreamLogsUseCase
	streamStatsUseCase *usecase.StreamStatsUseCase
	execUseCase        *usecase.ExecContainerUseCase
}

func NewContainerHandler(
	listUC *usecase.ListContainersUseCase,
	startUC *usecase.StartContainerUseCase,
	stopUC *usecase.StopContainerUseCase,
	pruneUC *usecase.PruneContainersUseCase,
	streamLogsUC *usecase.StreamLogsUseCase,
	streamStatsUC *usecase.StreamStatsUseCase,
	execUC *usecase.ExecContainerUseCase,
) *ContainerHandler {
	return &ContainerHandler{
		listUseCase:        listUC,
		startUseCase:       startUC,
		stopUseCase:        stopUC,
		pruneUseCase:       pruneUC,
		streamLogsUseCase:  streamLogsUC,
		streamStatsUseCase: streamStatsUC,
		execUseCase:        execUC,
	}
}

// HandleListContainers lida com o endpoint GET /api/containers
func (h *ContainerHandler) HandleListContainers(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	containers, err := h.listUseCase.Execute(req.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	json.NewEncoder(w).Encode(containers)
}

// HandleStartContainer lida com o endpoint POST /api/containers/{id}/start
func (h *ContainerHandler) HandleStartContainer(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := req.PathValue("id")

	err := h.startUseCase.Execute(req.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "container iniciado com sucesso"})
}

// HandleStopContainer lida com o endpoint POST /api/containers/{id}/stop
func (h *ContainerHandler) HandleStopContainer(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := req.PathValue("id")

	err := h.stopUseCase.Execute(req.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "container parado com sucesso"})
}

// HandlePruneContainers lida com o endpoint DELETE /api/containers/prune
func (h *ContainerHandler) HandlePruneContainers(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	err := h.pruneUseCase.Execute(req.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "limpeza de containers concluída"})
}

// HandleStreamLogs gerencia o streaming SSE de logs do container
func (h *ContainerHandler) HandleStreamLogs(w http.ResponseWriter, req *http.Request) {
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

	go h.streamLogsUseCase.Execute(ctx, id, logsChan, errChan)

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

// HandleStreamStats gerencia o streaming SSE de estatísticas de CPU/RAM
func (h *ContainerHandler) HandleStreamStats(w http.ResponseWriter, req *http.Request) {
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

	go h.streamStatsUseCase.Execute(ctx, id, statsChan, errChan)

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

// HandleExecContainer lida com a conexão WebSocket para o terminal interativo (V2)
func (h *ContainerHandler) HandleExecContainer(w http.ResponseWriter, req *http.Request) {
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
	err = h.execUseCase.Execute(req.Context(), id, stdinReader, wsWriter)
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
