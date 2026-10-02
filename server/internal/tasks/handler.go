package tasks

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
	mux.HandleFunc("POST /tasks", h.create)
	mux.HandleFunc("GET /tasks", h.list)
}

type createRequest struct {
	Title     string `json:"title"`
	Important bool   `json:"important"`
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	var req createRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		platform.WriteError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	title := strings.TrimSpace(req.Title)
	if title == "" || len(title) > 200 {
		platform.WriteError(w, http.StatusBadRequest, "title is required and must be at most 200 characters")
		return
	}

	t, err := h.store.Create(r.Context(), title, req.Important)
	if err != nil {
		log.Printf("create task: %v", err)
		platform.WriteError(w, http.StatusInternalServerError, "could not create task")
		return
	}
	platform.WriteJSON(w, http.StatusCreated, t)
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	items, err := h.store.ListActive(r.Context())
	if err != nil {
		log.Printf("list tasks: %v", err)
		platform.WriteError(w, http.StatusInternalServerError, "could not list tasks")
		return
	}
	platform.WriteJSON(w, http.StatusOK, items)
}
