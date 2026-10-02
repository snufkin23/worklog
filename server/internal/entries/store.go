package entries

import (
	"context"
	"fmt"

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
