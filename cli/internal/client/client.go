package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

type Task struct {
	ID            int64     `json:"id"`
	Title         string    `json:"title"`
	Status        string    `json:"status"`
	BlockerReason *string   `json:"blocker_reason"`
	Important     bool      `json:"important"`
	CreatedAt     time.Time `json:"created_at"`
}

type Event struct {
	ID        int64     `json:"id"`
	TaskID    *int64    `json:"task_id"`
	TaskTitle *string   `json:"task_title"`
	Type      string    `json:"type"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
}

type Today struct {
	Date        string  `json:"date"`
	ActiveTasks []Task  `json:"active_tasks"`
	Events      []Event `json:"events"`
}

type Client struct {
	baseURL string
	token   string
	http    *http.Client
}

func New(baseURL, token string) *Client {
	return &Client{
		baseURL: baseURL,
		token:   token,
		http:    &http.Client{Timeout: 15 * time.Second},
	}
}

func (c *Client) AddTask(title string, important bool) (Task, error) {
	var t Task
	err := c.do(http.MethodPost, "/tasks", map[string]any{"title": title, "important": important}, &t)
	return t, err
}

// Act applies a lifecycle action: progress, block, unblock, done or drop.
func (c *Client) Act(id int64, action, text string) (Task, error) {
	var t Task
	path := "/tasks/" + strconv.FormatInt(id, 10) + "/" + action
	err := c.do(http.MethodPost, path, map[string]string{"text": text}, &t)
	return t, err
}

func (c *Client) AddNote(text string) (Event, error) {
	var e Event
	err := c.do(http.MethodPost, "/notes", map[string]string{"text": text}, &e)
	return e, err
}

func (c *Client) Today() (Today, error) {
	var t Today
	err := c.do(http.MethodGet, "/today", nil, &t)
	return t, err
}

func (c *Client) do(method, path string, in, out any) error {
	var body io.Reader
	if in != nil {
		data, err := json.Marshal(in)
		if err != nil {
			return fmt.Errorf("encode request: %w", err)
		}
		body = bytes.NewReader(data)
	}

	req, err := http.NewRequest(method, c.baseURL+path, body)
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	if in != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(io.LimitReader(resp.Body, 4<<20))
	if err != nil {
		return fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode >= 300 {
		var apiErr struct {
			Error string `json:"error"`
		}
		if json.Unmarshal(data, &apiErr) == nil && apiErr.Error != "" {
			return errors.New(apiErr.Error)
		}
		return fmt.Errorf("server returned %s", resp.Status)
	}

	if out == nil {
		return nil
	}
	if err := json.Unmarshal(data, out); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}
	return nil
}
