package today

import (
	"log"
	"net/http"
	"time"

	"github.com/snufkin23/worklog/server/internal/entries"
	"github.com/snufkin23/worklog/server/internal/platform"
	"github.com/snufkin23/worklog/server/internal/tasks"
)

type Handler struct {
	taskStore  *tasks.Store
	entryStore *entries.Store
	loc        *time.Location
}

func NewHandler(taskStore *tasks.Store, entryStore *entries.Store, loc *time.Location) *Handler {
	return &Handler{taskStore: taskStore, entryStore: entryStore, loc: loc}
}

func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /today", h.get)
}

type response struct {
	Date        string          `json:"date"`
	ActiveTasks []tasks.Task    `json:"active_tasks"`
	Events      []entries.Event `json:"events"`
}

func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	start, end := platform.DayBounds(time.Now(), h.loc)

	active, err := h.taskStore.ListActive(r.Context())
	if err != nil {
		log.Printf("today: list tasks: %v", err)
		platform.WriteError(w, http.StatusInternalServerError, "could not load today")
		return
	}

	events, err := h.entryStore.ListBetween(r.Context(), start, end)
	if err != nil {
		log.Printf("today: list events: %v", err)
		platform.WriteError(w, http.StatusInternalServerError, "could not load today")
		return
	}

	platform.WriteJSON(w, http.StatusOK, response{
		Date:        start.Format("2006-01-02"),
		ActiveTasks: active,
		Events:      events,
	})
}
