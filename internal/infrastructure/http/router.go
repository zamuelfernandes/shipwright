package http

import (
	"io/fs"
	"net/http"

	"github.com/zamuelfernandes/shipwright/ui"
)

// Router maps HTTP routes, REST endpoints, WebSockets, and UI file servers to resource handlers.
type Router struct {
	mux              *http.ServeMux
	containerHandler *ContainerHandler
	projectHandler   *ProjectHandler
	imageHandler     *ImageHandler
}

func NewRouter(
	containerH *ContainerHandler,
	projectH *ProjectHandler,
	imageH *ImageHandler,
) *Router {
	mux := http.NewServeMux()

	r := &Router{
		mux:              mux,
		containerHandler: containerH,
		projectHandler:   projectH,
		imageHandler:     imageH,
	}

	// REST API - Container Lifecycle (ContainerHandler)
	mux.HandleFunc("GET /api/containers", r.containerHandler.HandleListContainers)
	mux.HandleFunc("POST /api/containers/{id}/start", r.containerHandler.HandleStartContainer)
	mux.HandleFunc("POST /api/containers/{id}/stop", r.containerHandler.HandleStopContainer)
	mux.HandleFunc("DELETE /api/containers/prune", r.containerHandler.HandlePruneContainers)

	// REST API - Docker Compose Projects (ProjectHandler)
	mux.HandleFunc("POST /api/projects/{name}/start", r.projectHandler.HandleStartProject)
	mux.HandleFunc("POST /api/projects/{name}/stop", r.projectHandler.HandleStopProject)

	// REST API - Docker Local Images (ImageHandler)
	mux.HandleFunc("GET /api/images", r.imageHandler.HandleListImages)

	// SSE Streaming Telemetry & Logs (ContainerHandler)
	mux.HandleFunc("GET /api/containers/{id}/logs", r.containerHandler.HandleStreamLogs)
	mux.HandleFunc("GET /api/containers/{id}/stats", r.containerHandler.HandleStreamStats)

	// WebSocket Terminal Emulator Integration (ContainerHandler)
	mux.HandleFunc("GET /api/containers/{id}/exec", r.containerHandler.HandleExecContainer)

	// Embedded Static assets (embedded HTML/CSS/JS files)
	distFS, err := fs.Sub(ui.DistFS, "dist")
	if err != nil {
		panic("error mapping dist UI folder: " + err.Error())
	}
	fileServer := http.FileServer(http.FS(distFS))
	mux.Handle("/", fileServer)

	return r
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}
