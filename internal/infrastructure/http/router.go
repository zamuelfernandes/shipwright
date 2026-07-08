package http

import (
	"io/fs"
	"net/http"

	"github.com/zamuelfernandes/shipwright/ui"
)

// Router gerencia as rotas HTTP da aplicação, endpoints REST e conexões WebSocket delegando para sub-handlers dedicados.
type Router struct {
	mux              *http.ServeMux
	containerHandler *ContainerHandler
	projectHandler   *ProjectHandler
	imageHandler     *ImageHandler
}

// NewRouter recebe os sub-handlers necessários por injeção de dependência e configura as rotas.
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

	// API REST - Containers individuais (ContainerHandler)
	mux.HandleFunc("GET /api/containers", r.containerHandler.HandleListContainers)
	mux.HandleFunc("POST /api/containers/{id}/start", r.containerHandler.HandleStartContainer)
	mux.HandleFunc("POST /api/containers/{id}/stop", r.containerHandler.HandleStopContainer)
	mux.HandleFunc("DELETE /api/containers/prune", r.containerHandler.HandlePruneContainers)

	// API REST - Projetos Docker Compose (ProjectHandler)
	mux.HandleFunc("POST /api/projects/{name}/start", r.projectHandler.HandleStartProject)
	mux.HandleFunc("POST /api/projects/{name}/stop", r.projectHandler.HandleStopProject)

	// API REST - Imagens Docker (ImageHandler)
	mux.HandleFunc("GET /api/images", r.imageHandler.HandleListImages)

	// API de Streaming SSE (ContainerHandler)
	mux.HandleFunc("GET /api/containers/{id}/logs", r.containerHandler.HandleStreamLogs)
	mux.HandleFunc("GET /api/containers/{id}/stats", r.containerHandler.HandleStreamStats)

	// Terminal Interativo WebSocket (ContainerHandler)
	mux.HandleFunc("GET /api/containers/{id}/exec", r.containerHandler.HandleExecContainer)

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
