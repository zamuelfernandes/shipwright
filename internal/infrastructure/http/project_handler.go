package http

import (
	"encoding/json"
	"net/http"

	"github.com/zamuelfernandes/shipwright/internal/usecase"
)

// ProjectHandler manages Docker Compose project HTTP routes.
type ProjectHandler struct {
	startProjectUseCase *usecase.StartProjectUseCase
	stopProjectUseCase  *usecase.StopProjectUseCase
}

func NewProjectHandler(
	startProjUC *usecase.StartProjectUseCase,
	stopProjUC *usecase.StopProjectUseCase,
) *ProjectHandler {
	return &ProjectHandler{
		startProjectUseCase: startProjUC,
		stopProjectUseCase:  stopProjUC,
	}
}

// HandleStartProject handles POST /api/projects/{name}/start.
func (h *ProjectHandler) HandleStartProject(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	name := req.PathValue("name")

	err := h.startProjectUseCase.Execute(req.Context(), name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "project started successfully"})
}

// HandleStopProject handles POST /api/projects/{name}/stop.
func (h *ProjectHandler) HandleStopProject(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	name := req.PathValue("name")

	err := h.stopProjectUseCase.Execute(req.Context(), name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "project stopped successfully"})
}
