package tasks

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrNotFound     = errors.New("task not found")
	ErrInvalidState = errors.New("action not allowed in the task's current state")
)

const taskColumns = `id, title, status, blocker_reason, important, created_at, updated_at, closed_at`

type Store struct {
	pool *pgxpool.Pool
}

func NewStore(pool *pgxpool.Pool) *Store {
	return &Store{pool: pool}
}

func scanTask(row pgx.Row) (Task, error) {
	var t Task
	err := row.Scan(&t.ID, &t.Title, &t.Status, &t.BlockerReason, &t.Important, &t.CreatedAt, &t.UpdatedAt, &t.ClosedAt)
	return t, err
}

// Create inserts the task and its "created" event in one transaction.
func (s *Store) Create(ctx context.Context, title string, important bool) (Task, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return Task{}, fmt.Errorf("begin: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }() // no-op once committed

	t, err := scanTask(tx.QueryRow(ctx,
		`INSERT INTO tasks (title, important) VALUES ($1, $2) RETURNING `+taskColumns,
		title, important))
	if err != nil {
		return Task{}, fmt.Errorf("insert task: %w", err)
	}

	if _, err := tx.Exec(ctx,
		`INSERT INTO events (task_id, type, text) VALUES ($1, 'created', $2)`,
		t.ID, title); err != nil {
		return Task{}, fmt.Errorf("insert event: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return Task{}, fmt.Errorf("commit: %w", err)
	}
	return t, nil
}

// ListActive returns open and blocked tasks, oldest first.
func (s *Store) ListActive(ctx context.Context) ([]Task, error) {
	rows, err := s.pool.Query(ctx,
		`SELECT `+taskColumns+` FROM tasks WHERE status IN ('open', 'blocked') ORDER BY id`)
	if err != nil {
		return nil, fmt.Errorf("query tasks: %w", err)
	}
	defer rows.Close()

	tasks := []Task{}
	for rows.Next() {
		t, err := scanTask(rows)
		if err != nil {
			return nil, fmt.Errorf("scan task: %w", err)
		}
		tasks = append(tasks, t)
	}
	return tasks, rows.Err()
}

// Apply performs a lifecycle action: updates the task and writes its event in one transaction.
func (s *Store) Apply(ctx context.Context, id int64, action Action, text string) (Task, error) {
	r, ok := rules[action]
	if !ok {
		return Task{}, fmt.Errorf("unknown action %q", action)
	}

	var reason *string
	if action == Block {
		reason = &text
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return Task{}, fmt.Errorf("begin: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	t, err := scanTask(tx.QueryRow(ctx,
		`UPDATE tasks SET
			status         = CASE WHEN $2::text = '' THEN status ELSE $2::text END,
			blocker_reason = CASE WHEN $2::text = '' THEN blocker_reason ELSE $3::text END,
			updated_at     = now(),
			closed_at      = CASE WHEN $2::text IN ('done', 'dropped') THEN now() ELSE NULL END
		WHERE id = $1 AND status = ANY($4::text[])
		RETURNING `+taskColumns,
		id, r.to, reason, r.from))
	if errors.Is(err, pgx.ErrNoRows) {
		return Task{}, diagnose(ctx, tx, id)
	}
	if err != nil {
		return Task{}, fmt.Errorf("update task: %w", err)
	}

	if _, err := tx.Exec(ctx,
		`INSERT INTO events (task_id, type, text) VALUES ($1, $2, $3)`,
		id, r.event, text); err != nil {
		return Task{}, fmt.Errorf("insert event: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return Task{}, fmt.Errorf("commit: %w", err)
	}
	return t, nil
}

// diagnose explains why an update matched no rows: missing task or wrong state.
func diagnose(ctx context.Context, tx pgx.Tx, id int64) error {
	var status string
	err := tx.QueryRow(ctx, `SELECT status FROM tasks WHERE id = $1`, id).Scan(&status)
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}
	if err != nil {
		return fmt.Errorf("check task: %w", err)
	}
	return fmt.Errorf("%w (task is %s)", ErrInvalidState, status)
}
