package api

import (
	"fmt"
	"net/http"
)

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Color    string `json:"color"`
	Initials string `json:"initials"`
	Avatar   string `json:"avatar"`
}

type Dependency struct {
	TaskID   string `json:"task_id"`
	DependsOn string `json:"depends_on"`
	Type     int    `json:"type"`
}

type Task struct {
	ID          string         `json:"id"`
	CustomID    string         `json:"custom_id"`
	Name        string         `json:"name"`
	TextContent string         `json:"text_content"`
	Description string         `json:"description"`
	Status      *struct {
		ID      string `json:"id"`
		Status  string `json:"status"`
		Color   string `json:"color"`
		OrderBy int    `json:"orderby"`
	} `json:"status"`
	OrderIndex     string `json:"orderindex"`
	DateCreated    string `json:"date_created"`
	DateUpdated    string `json:"date_updated"`
	DateClosed     string `json:"date_closed"`
	DueDate        string `json:"due_date"`
	StartDate      string `json:"start_date"`
	Priority       *struct {
		ID       int    `json:"id"`
		Priority string `json:"priority"`
		Color    string `json:"color"`
		OrderBy  int    `json:"orderby"`
	} `json:"priority"`
	Assignee *User   `json:"assignee"`
	Assignees []User `json:"assignees"`
	ParentID string `json:"parent"`
	ListID   string `json:"list"`
	Subtasks []Task `json:"subtasks"`
}

type Comment struct {
	ID           string `json:"id"`
	HistoryID    string `json:"history_id"`
	TextContent  string `json:"text_content"`
	User         User   `json:"user"`
	Resolved     bool   `json:"resolved"`
	DateCreated  string `json:"date_created"`
	DateUpdated  string `json:"date_updated"`
}

type CommentsResponse struct {
	Comments []Comment `json:"comments"`
}

type TaskListResponse struct {
	Tasks []Task `json:"tasks"`
}

func GetTasks(c *Client, listID string, recursive bool) (TaskListResponse, error) {
	path := fmt.Sprintf("/list/%s/task?archived=false", listID)
	if recursive {
		path += "&subtasks=true"
	}
	return Do[any, TaskListResponse](c, http.MethodGet, path, nil)
}

func GetTask(c *Client, taskID string) (Task, error) {
	path := fmt.Sprintf("/task/%s", taskID)
	return Do[any, Task](c, http.MethodGet, path, nil)
}

func GetTaskComments(c *Client, taskID string) ([]Comment, error) {
	path := fmt.Sprintf("/task/%s/comment", taskID)
	resp, err := Do[any, CommentsResponse](c, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	return resp.Comments, nil
}

func CreateTask(c *Client, payload map[string]any) (Task, error) {
	var task Task

	listID, ok := payload["list_id"].(string)
	if !ok || listID == "" {
		return task, fmt.Errorf("list_id is required")
	}

	path := fmt.Sprintf("/list/%s/task", listID)
	return Do[map[string]any, Task](c, http.MethodPost, path, &payload)
}

func DeleteTask(c *Client, taskID string) error {
	path := fmt.Sprintf("/task/%s", taskID)
	_, err := Do[any, any](c, http.MethodDelete, path, nil)
	return err
}

func ArchiveTask(c *Client, taskID string) error {
	path := fmt.Sprintf("/task/%s/archive", taskID)
	_, err := Do[any, any](c, http.MethodPut, path, nil)
	return err
}

func UpdateTask(c *Client, taskID string, payload map[string]any) (Task, error) {
	path := fmt.Sprintf("/task/%s", taskID)
	return Do[map[string]any, Task](c, http.MethodPut, path, &payload)
}
