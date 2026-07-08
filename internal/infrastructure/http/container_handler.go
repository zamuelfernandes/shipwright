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

// upgrader configures WebSocket connections.
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for local development.
	},
}

// ContainerHandler manages container-related HTTP routes, SSE telemetry streams, and WebSockets.
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

// HandleListContainers handles GET /api/containers.
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

// HandleStartContainer handles POST /api/containers/{id}/start.
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
	json.NewEncoder(w).Encode(map[string]string{"message": "container started successfully"})
}

// HandleStopContainer handles POST /api/containers/{id}/stop.
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
	json.NewEncoder(w).Encode(map[string]string{"message": "container stopped successfully"})
}

// HandlePruneContainers handles DELETE /api/containers/prune.
func (h *ContainerHandler) HandlePruneContainers(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	err := h.pruneUseCase.Execute(req.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "prune operation finished successfully"})
}

// HandleStreamLogs handles SSE logs stream on GET /api/containers/{id}/logs.
func (h *ContainerHandler) HandleStreamLogs(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	id := req.PathValue("id")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming not supported", http.StatusInternalServerError)
		return
	}

	logsChan := make(chan string, 10)
	errChan := make(chan error, 1)

	ctx, cancel := context.WithCancel(req.Context())
	defer cancel()

	go h.streamLogsUseCase.Execute(ctx, id, logsChan, errChan)

	fmt.Fprint(w, "event: open\ndata: connected to logs\n\n")
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

// HandleStreamStats handles SSE stats stream on GET /api/containers/{id}/stats.
func (h *ContainerHandler) HandleStreamStats(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	id := req.PathValue("id")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming not supported", http.StatusInternalServerError)
		return
	}

	statsChan := make(chan domain.ContainerStats, 10)
	errChan := make(chan error, 1)

	ctx, cancel := context.WithCancel(req.Context())
	defer cancel()

	go h.streamStatsUseCase.Execute(ctx, id, statsChan, errChan)

	fmt.Fprint(w, "event: open\ndata: connected to stats\n\n")
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

// HandleExecContainer handles WebSocket connections for interactive terminal.
func (h *ContainerHandler) HandleExecContainer(w http.ResponseWriter, req *http.Request) {
	id := req.PathValue("id")
	log.Printf("[HTTP] Exec request received for container: %s", id)

	conn, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Printf("[HTTP] WebSocket upgrade failed: %v", err)
		return
	}
	defer conn.Close()
	log.Printf("[HTTP] WebSocket connection active for container: %s", id)

	stdinReader, stdinWriter := io.Pipe()

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

	wsWriter := &wsWriterWrapper{conn: conn}

	err = h.execUseCase.Execute(req.Context(), id, stdinReader, wsWriter)
	if err != nil {
		log.Printf("[HTTP] Exec session error: %v", err)
		_ = conn.WriteMessage(websocket.TextMessage, []byte("\r\nTerminal initialization failed: "+err.Error()+"\r\n"))
	} else {
		log.Printf("[HTTP] Exec session closed for container: %s", id)
	}
}

// wsWriterWrapper wraps WebSocket connections to satisfy io.Writer interface.
type wsWriterWrapper struct {
	conn *websocket.Conn
}

func (w *wsWriterWrapper) Write(p []byte) (n int, err error) {
	err = w.conn.WriteMessage(websocket.TextMessage, p)
	if err != nil {
		return 0, err
	}
	return len(p), nil
}
