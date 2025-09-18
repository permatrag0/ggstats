package http

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"ggstats.com/matches/internal/controller/matches"
	"ggstats.com/matches/pkg/model"
)

type Handler struct {
	ctrl *matches.Controller
}

func New(ctrl *matches.Controller) *Handler {
	return &Handler{ctrl}
}

func (h *Handler) Handle(w http.ResponseWriter, req *http.Request) {
	// recordID := model.RecordID(req.FormValue("id"))

	// if recordID == "" {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	return
	// }
	// recordType := model.RecordType(req.FormValue("type"))
	// if recordType == "" {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	return
	// }

	switch req.Method {
	// case http.MethodGet:
	// 	v, err := h.ctrl.GetAggregatedRating(req.Context(), recordID, recordType)
	// 	if err != nil && errors.Is(err, rating.ErrNotFound) {
	// 		w.WriteHeader(http.StatusNotFound)
	// 		return
	// 	}
	// 	if err := json.NewEncoder(w).Encode(v); err != nil {
	// 		log.Printf("Response encode error: %v \n", err)
	// 	}
	case http.MethodGet:
		// Handle GET request
		recordID := model.RecordID(req.FormValue("id"))
		recordType := model.RecordType(req.FormValue("type"))

		if recordID == "" || recordType == "" {
			http.Error(w, "Both 'id' and 'type' are required query parameters", http.StatusBadRequest)
			return
		}

		m, err := h.ctrl.GetMatches(req.Context(), recordID, recordType) // Call the new Get method in controller
		if err != nil {
			log.Printf("Error fetching matches: %v", err)
			if errors.Is(err, matches.ErrNotFound) {
				http.Error(w, "Record not found", http.StatusNotFound)
			} else {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(m); err != nil {
			log.Printf("Response encoding error: %v\n", err)
			// Don't write header again if already written by Encode
		}
		return // Important: Return after successful GET
	case http.MethodPost:
		var match model.Matches
		if err := json.NewDecoder(req.Body).Decode(&match); err != nil {
			log.Printf("Match decoding error: %v\n", err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Use the recordType and recordID from the decoded match
		recordID := match.RecordID
		recordType := match.RecordType

		if recordID == "" || recordType == "" {
			http.Error(w, "Both 'recordId' and 'recordType' are required in the JSON body", http.StatusBadRequest)
			return
		}

		if err := h.ctrl.PutMatches(req.Context(), model.RecordID(match.RecordID), model.RecordType(match.RecordType), &match); err != nil { // Still uses PutRating
			log.Printf("Controller Post error: %v\n", err)
			http.Error(w, "Failed to save match", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated) // 201 Created for new resource
		return                            // Important: Return after successful POST
		// tournament := model.Tournament(req.FormValue("Tournament"))
		// player1 := model.Player1(req.FormValue("Player1"))
		// player2 := model.Player2(req.FormValue("Player2"))
		// scorep1, _ := strconv.ParseFloat(req.FormValue("Scorep1"), 64)
		// scorep2, err := strconv.ParseFloat(req.FormValue("Scorep2"), 64)
		// if err != nil && errors.Is(err, matches.ErrNotFound) {
		// 	w.WriteHeader(http.StatusBadRequest)
		// 	return
		// }
		// if err := h.ctrl.PutMatches(req.Context(), recordID, recordType, &model.Matches{Tournament: string(tournament), Player1: string(player1),
		// 	Player2: string(player2), Scorep1: int(scorep1), Scorep2: int(scorep2)}); err != nil {
		// 	log.Printf("Repository put error: %v\n", err)
		// 	w.WriteHeader(http.StatusInternalServerError)
		// }
	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}
