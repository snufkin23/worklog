package tasks

import "testing"

func TestAllowed(t *testing.T) {
	tests := []struct {
		action Action
		status string
		want   bool
	}{
		{Progress, "open", true},
		{Progress, "blocked", true},
		{Progress, "done", false},
		{Block, "open", true},
		{Block, "blocked", false},
		{Block, "done", false},
		{Unblock, "blocked", true},
		{Unblock, "open", false},
		{Done, "open", true},
		{Done, "blocked", true},
		{Done, "done", false},
		{Done, "dropped", false},
		{Drop, "open", true},
		{Drop, "blocked", true},
		{Drop, "dropped", false},
		{Action("nonsense"), "open", false},
	}

	for _, tt := range tests {
		got := Allowed(tt.action, tt.status)
		if got != tt.want {
			t.Errorf("Allowed(%q, %q) = %v, want %v", tt.action, tt.status, got, tt.want)
		}
	}
}
