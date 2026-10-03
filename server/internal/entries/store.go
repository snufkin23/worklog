package entries

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	pool *pgxpool.Pool
}

func NewStore(pool *pgxpool.Pool) *Store {
	return &Store{pool: pool}
}

// AddNote records a standalone note: an event with no task.
func (s *Store) AddNote(ctx context.Context, text string) (Event, error) {
	var e Event
	err := s.pool.QueryRow(ctx,
		`INSERT INTO events (type, text) VALUES ('note', $1)
		 RETURNING id, task_id, type, text, created_at`, text).
		Scan(&e.ID, &e.TaskID, &e.Type, &e.Text, &e.CreatedAt)
	if err != nil {
		return Event{}, fmt.Errorf("insert note: %w", err)
	}
	return e, nil
}

// ListBetween returns events with from <= created_at < to, oldest first,
// with the owning task's title attached (reads may join across tables).
func (s *Store) ListBetween(ctx context.Context, from, to time.Time) ([]Event, error) {
	rows, err := s.pool.Query(ctx,
		`SELECT e.id, e.task_id, t.title, e.type, e.text, e.created_at
		 FROM events e
		 LEFT JOIN tasks t ON t.id = e.task_id
		 WHERE e.created_at >= $1 AND e.created_at < $2
		 ORDER BY e.created_at, e.id`, from, to)
	if err != nil {
		return nil, fmt.Errorf("query events: %w", err)
	}
	defer rows.Close()

	events := []Event{}
	for rows.Next() {
		var e Event
		if err := rows.Scan(&e.ID, &e.TaskID, &e.TaskTitle, &e.Type, &e.Text, &e.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan event: %w", err)
		}
		events = append(events, e)
	}
	return events, rows.Err()
}
