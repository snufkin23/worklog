package tasks

import "time"

type Task struct {
	ID            int64      `json:"id"`
	Title         string     `json:"title"`
	Status        string     `json:"status"`
	BlockerReason *string    `json:"blocker_reason,omitempty"`
	Important     bool       `json:"important"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	ClosedAt      *time.Time `json:"closed_at,omitempty"`
}
