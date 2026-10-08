package main

import (
	"fmt"
	"os"

	"github.com/snufkin23/worklog/cli/internal/client"
	"github.com/snufkin23/worklog/cli/internal/config"
)

const usage = `wl: worklog command line

Usage:
  wl                      show today (default)
  wl today                show today
  wl add [-i] <title>     add a task (-i marks it important)
  wl progress <n> <text>  log progress on task n
  wl block <n> <reason>   mark task n blocked
  wl unblock <n>          clear the blocker on task n
  wl done <n>             mark task n done
  wl drop <n>             drop task n
  wl note <text>          log a standalone note
  wl config               set the server URL and token
`

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "wl:", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	cmd := "today"
	rest := []string{}
	if len(args) > 0 {
		cmd, rest = args[0], args[1:]
	}

	switch cmd {
	case "help", "-h", "--help":
		fmt.Print(usage)
		return nil
	case "config":
		return configure()
	}

	cfg, err := config.Load()
	if err != nil {
		return err
	}
	c := client.New(cfg.URL, cfg.Token)

	switch cmd {
	case "today":
		return showToday(c)
	case "add":
		return addTask(c, rest)
	case "note":
		return addNote(c, rest)
	case "progress", "block":
		return act(c, cmd, rest, true)
	case "unblock", "done", "drop":
		return act(c, cmd, rest, false)
	default:
		return fmt.Errorf("unknown command %q (try: wl help)", cmd)
	}
}
