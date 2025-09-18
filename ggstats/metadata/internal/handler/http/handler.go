package http

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"ggstats.com/metadata/internal/controller/metadata"
	"ggstats.com/metadata/internal/repository"
	model "ggstats.com/metadata/pkg"
)

type Handler struct {
	ctrl *metadata.Controller
}

func New(ctrl *metadata.Controller) *Handler {
	return &Handler{ctrl}
}

func (h *Handler) GetMetadata(w http.ResponseWriter, req *http.Request) {
	id := req.FormValue("id")

	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
	}
	ctx := req.Context()
	m, err := h.ctrl.Get(ctx, id)
	if err != nil && errors.Is(err, repository.ErrNotFound) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		log.Printf("Repository error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(m); err != nil {
		log.Printf("Response error: v%\n", err)
	}
}

func (h *Handler) CreateMetadata(w http.ResponseWriter, r *http.Request) {
	var metadata model.Metadata
	if err := json.NewDecoder(r.Body).Decode(&metadata); err != nil {
		log.Printf("failed to decode metadata: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	if err := h.ctrl.Create(ctx, &metadata); err != nil {
		log.Printf("failed to create metadata: %v", err)
		http.Error(w, "Failed to create metadata", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
