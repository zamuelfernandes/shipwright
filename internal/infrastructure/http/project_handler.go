package http

import (
	"encoding/json"
	"net/http"

	"github.com/zamuelfernandes/shipwright/internal/usecase"
)

// ProjectHandler gerencia as rotas HTTP associadas a projetos Docker Compose.
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

// HandleStartProject lida com o endpoint POST /api/projects/{name}/start
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
	json.NewEncoder(w).Encode(map[string]string{"message": "projeto inicializado com sucesso"})
}

// HandleStopProject lida com o endpoint POST /api/projects/{name}/stop
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
	json.NewEncoder(w).Encode(map[string]string{"message": "projeto pausado com sucesso"})
}
