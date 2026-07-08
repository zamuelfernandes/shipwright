package http

import (
	"encoding/json"
	"net/http"

	"github.com/zamuelfernandes/shipwright/internal/usecase"
)

// ImageHandler manages Docker image HTTP routes.
type ImageHandler struct {
	listImagesUseCase *usecase.ListImagesUseCase
}

func NewImageHandler(listImagesUC *usecase.ListImagesUseCase) *ImageHandler {
	return &ImageHandler{
		listImagesUseCase: listImagesUC,
	}
}

// HandleListImages handles GET /api/images.
func (h *ImageHandler) HandleListImages(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	images, err := h.listImagesUseCase.Execute(req.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	json.NewEncoder(w).Encode(images)
}
