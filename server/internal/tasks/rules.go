package tasks

import "slices"

type Action string

const (
	Progress Action = "progress"
	Block    Action = "block"
	Unblock  Action = "unblock"
	Done     Action = "done"
	Drop     Action = "drop"
)

type rule struct {
	from  []string // statuses the action is allowed from
	to    string   // new status; empty means unchanged
	event string   // event type written to the history
}

var rules = map[Action]rule{
	Progress: {from: []string{"open", "blocked"}, to: "", event: "progress"},
	Block:    {from: []string{"open"}, to: "blocked", event: "blocked"},
	Unblock:  {from: []string{"blocked"}, to: "open", event: "unblocked"},
	Done:     {from: []string{"open", "blocked"}, to: "done", event: "done"},
	Drop:     {from: []string{"open", "blocked"}, to: "dropped", event: "dropped"},
}

// Allowed reports whether an action may be applied to a task in the given status.
func Allowed(action Action, status string) bool {
	r, ok := rules[action]
	if !ok {
		return false
	}
	return slices.Contains(r.from, status)
}
