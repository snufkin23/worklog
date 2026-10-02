package tasks

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"
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
	for _, a := range []Action{Progress, Block, Unblock, Done, Drop} {
		mux.HandleFunc("POST /tasks/{id}/"+string(a), h.act(a))
	}
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

type actionRequest struct {
	Text string `json:"text"`
}

func (h *Handler) act(action Action) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil || id < 1 {
			platform.WriteError(w, http.StatusBadRequest, "invalid task id")
			return
		}

		r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
		var req actionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil && !errors.Is(err, io.EOF) {
			platform.WriteError(w, http.StatusBadRequest, "invalid JSON body")
			return
		}

		text := strings.TrimSpace(req.Text)
		if len(text) > 1000 {
			platform.WriteError(w, http.StatusBadRequest, "text must be at most 1000 characters")
			return
		}
		if (action == Progress || action == Block) && text == "" {
			platform.WriteError(w, http.StatusBadRequest, "text is required for this action")
			return
		}

		t, err := h.store.Apply(r.Context(), id, action, text)
		switch {
		case errors.Is(err, ErrNotFound):
			platform.WriteError(w, http.StatusNotFound, err.Error())
		case errors.Is(err, ErrInvalidState):
			platform.WriteError(w, http.StatusConflict, err.Error())
		case err != nil:
			log.Printf("%s task %d: %v", action, id, err)
			platform.WriteError(w, http.StatusInternalServerError, "could not update task")
		default:
			platform.WriteJSON(w, http.StatusOK, t)
		}
	}
}
