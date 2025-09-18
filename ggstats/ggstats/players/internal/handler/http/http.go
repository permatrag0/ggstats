package http

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"ggstats.com/players/internal/controller/players"
)

type Handler struct {
	ctrl *players.Controller
}

func New(ctrl *players.Controller) *Handler {
	return &Handler{ctrl}
}

func (h *Handler) GetPlayersDetails(w http.ResponseWriter, req *http.Request) {
	id := req.FormValue("id")
	details, err := h.ctrl.Get(req.Context(), id)
	if err != nil && errors.Is(err, players.ErrNotFound) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		log.Printf("Repository get error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(details); err != nil {
		log.Printf("Response encode error: %v\n", err)
	}
}
