package entries

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/snufkin23/worklog/server/internal/platform"
)

type Handler struct {
	store *Store
}

func NewHandler(store *Store) *Handler {
	return &Handler{store: store}
}

func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("POST /notes", h.addNote)
}

type noteRequest struct {
	Text string `json:"text"`
}

func (h *Handler) addNote(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	var req noteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		platform.WriteError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	text := strings.TrimSpace(req.Text)
	if text == "" || len(text) > 1000 {
		platform.WriteError(w, http.StatusBadRequest, "text is required and must be at most 1000 characters")
		return
	}

	e, err := h.store.AddNote(r.Context(), text)
	if err != nil {
		log.Printf("add note: %v", err)
		platform.WriteError(w, http.StatusInternalServerError, "could not add note")
		return
	}
	platform.WriteJSON(w, http.StatusCreated, e)
}
