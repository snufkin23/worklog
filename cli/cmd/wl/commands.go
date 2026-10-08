package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/snufkin23/worklog/cli/internal/client"
	"github.com/snufkin23/worklog/cli/internal/config"
)

func parseID(s string) (int64, error) {
	id, err := strconv.ParseInt(strings.TrimPrefix(strings.TrimSpace(s), "#"), 10, 64)
	if err != nil || id < 1 {
		return 0, fmt.Errorf("invalid task number %q", s)
	}
	return id, nil
}

func showToday(c *client.Client) error {
	t, err := c.Today()
	if err != nil {
		return err
	}
	printToday(os.Stdout, t, appLocation())
	return nil
}

func addTask(c *client.Client, args []string) error {
	fs := flag.NewFlagSet("add", flag.ContinueOnError)
	important := fs.Bool("i", false, "mark the task as important")
	if err := fs.Parse(args); err != nil {
		return err
	}

	title := strings.TrimSpace(strings.Join(fs.Args(), " "))
	if title == "" {
		return errors.New("usage: wl add [-i] <title>")
	}

	t, err := c.AddTask(title, *important)
	if err != nil {
		return err
	}
	fmt.Printf("added #%d %s\n", t.ID, t.Title)
	return nil
}

func addNote(c *client.Client, args []string) error {
	text := strings.TrimSpace(strings.Join(args, " "))
	if text == "" {
		return errors.New("usage: wl note <text>")
	}
	if _, err := c.AddNote(text); err != nil {
		return err
	}
	fmt.Println("noted")
	return nil
}

func act(c *client.Client, action string, args []string, needText bool) error {
	if len(args) == 0 {
		if needText {
			return fmt.Errorf("usage: wl %s <n> <text>", action)
		}
		return fmt.Errorf("usage: wl %s <n>", action)
	}

	id, err := parseID(args[0])
	if err != nil {
		return err
	}

	text := strings.TrimSpace(strings.Join(args[1:], " "))
	if needText && text == "" {
		return fmt.Errorf("usage: wl %s <n> <text>", action)
	}

	t, err := c.Act(id, action, text)
	if err != nil {
		return err
	}
	fmt.Printf("#%d is now %s: %s\n", t.ID, t.Status, t.Title)
	return nil
}

func configure() error {
	reader := bufio.NewReader(os.Stdin)

	url, err := prompt(reader, "Server URL (for example http://localhost:8080): ")
	if err != nil {
		return err
	}
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return errors.New("URL must start with http:// or https://")
	}

	token, err := prompt(reader, "API token: ")
	if err != nil {
		return err
	}

	path, err := config.Save(config.Config{URL: url, Token: token})
	if err != nil {
		return err
	}
	fmt.Println("saved to", path)
	return nil
}

func prompt(r *bufio.Reader, label string) (string, error) {
	fmt.Print(label)
	line, err := r.ReadString('\n')
	if err != nil && (!errors.Is(err, io.EOF) || line == "") {
		return "", err
	}
	v := strings.TrimSpace(line)
	if v == "" {
		return "", errors.New("a value is required")
	}
	return v, nil
}
