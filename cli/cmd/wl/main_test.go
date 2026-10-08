package main

import (
	"testing"

	"github.com/snufkin23/worklog/cli/internal/client"
)

func TestParseID(t *testing.T) {
	tests := []struct {
		in      string
		want    int64
		wantErr bool
	}{
		{"12", 12, false},
		{"#12", 12, false},
		{" 7 ", 7, false},
		{"0", 0, true},
		{"-3", 0, true},
		{"abc", 0, true},
		{"", 0, true},
	}

	for _, tt := range tests {
		got, err := parseID(tt.in)
		if (err != nil) != tt.wantErr {
			t.Errorf("parseID(%q) error = %v, wantErr %v", tt.in, err, tt.wantErr)
			continue
		}
		if got != tt.want {
			t.Errorf("parseID(%q) = %d, want %d", tt.in, got, tt.want)
		}
	}
}

func TestDescribe(t *testing.T) {
	id := int64(2)
	title := "write rules"

	tests := []struct {
		name string
		e    client.Event
		want string
	}{
		{"standalone note", client.Event{Type: "note", Text: "standup at 10:30"}, "standup at 10:30"},
		{"created repeats the title once", client.Event{TaskID: &id, TaskTitle: &title, Type: "created", Text: title}, "#2 write rules"},
		{"progress adds its text", client.Event{TaskID: &id, TaskTitle: &title, Type: "progress", Text: "merged"}, "#2 write rules — merged"},
		{"done has no text", client.Event{TaskID: &id, TaskTitle: &title, Type: "done"}, "#2 write rules"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := describe(tt.e); got != tt.want {
				t.Errorf("describe() = %q, want %q", got, tt.want)
			}
		})
	}
}
