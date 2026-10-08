package main

import (
	"fmt"
	"io"
	"time"
	_ "time/tzdata" // embed zone data so Nepal time works on any machine

	"github.com/snufkin23/worklog/cli/internal/client"
)

func appLocation() *time.Location {
	loc, err := time.LoadLocation("Asia/Kathmandu")
	if err != nil {
		return time.Local
	}
	return loc
}

func printToday(w io.Writer, t client.Today, loc *time.Location) {
	fmt.Fprintf(w, "Today · %s\n", t.Date)

	fmt.Fprintln(w, "\nActive")
	if len(t.ActiveTasks) == 0 {
		fmt.Fprintln(w, "  nothing open")
	}
	for _, task := range t.ActiveTasks {
		mark := " "
		if task.Important {
			mark = "!"
		}
		status := task.Status
		if task.Status == "blocked" && task.BlockerReason != nil {
			status += ": " + *task.BlockerReason
		}
		fmt.Fprintf(w, "  #%-3d %s %s  [%s]\n", task.ID, mark, task.Title, status)
	}

	fmt.Fprintln(w, "\nLog")
	if len(t.Events) == 0 {
		fmt.Fprintln(w, "  nothing logged yet")
	}
	for _, e := range t.Events {
		fmt.Fprintf(w, "  %s  %-9s %s\n", e.CreatedAt.In(loc).Format("15:04"), e.Type, describe(e))
	}
}

// describe renders what an event is about: a task reference and any extra text.
func describe(e client.Event) string {
	if e.TaskID == nil {
		return e.Text
	}
	label := fmt.Sprintf("#%d", *e.TaskID)
	if e.TaskTitle != nil {
		label += " " + *e.TaskTitle
	}
	if e.Text != "" && (e.TaskTitle == nil || e.Text != *e.TaskTitle) {
		label += " — " + e.Text
	}
	return label
}
