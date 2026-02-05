package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// StringOrNumber handles JSON fields that can be either string or number
type StringOrNumber string

func (s *StringOrNumber) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err == nil {
		*s = StringOrNumber(str)
		return nil
	}

	var num json.Number
	if err := json.Unmarshal(data, &num); err == nil {
		*s = StringOrNumber(num.String())
		return nil
	}

	return fmt.Errorf("cannot unmarshal %s into StringOrNumber", data)
}

type User struct {
	ID       StringOrNumber `json:"id"`
	Username string         `json:"username"`
	Email    string         `json:"email"`
	Color    string         `json:"color"`
	Initials string         `json:"initials"`
	Avatar   string         `json:"avatar"`
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
		ID       string `json:"id"`
		Priority string `json:"priority"`
		Color    string `json:"color"`
		OrderBy  int    `json:"orderby"`
	} `json:"priority"`
	Assignee  *User  `json:"assignee"`
	Assignees []User `json:"assignees"`
	ParentID  string `json:"parent"`
	List      *struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"list"`
	Subtasks []Task `json:"subtasks"`
}

type CommentSegment struct {
	Text       string `json:"text"`
	Attributes struct {
		Link string `json:"link"`
	} `json:"attributes"`
}

type Comment struct {
	ID          string           `json:"id"`
	HistoryID   string           `json:"hist_id"`
	CommentText string           `json:"comment_text"`
	Content     []CommentSegment `json:"comment"`
	User        User             `json:"user"`
	Resolved    bool             `json:"resolved"`
	Date        string           `json:"date"`
}

func (c Comment) RenderText() string {
	if len(c.Content) == 0 {
		return c.CommentText
	}
	var result string
	for _, seg := range c.Content {
		result += seg.Text
		if seg.Attributes.Link != "" {
			result += " (" + seg.Attributes.Link + ")"
		}
	}
	return result
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
	resp, err := Do[any, TaskListResponse](c, http.MethodGet, path, nil)
	if err != nil {
		return resp, err
	}
	if recursive {
		resp.Tasks = buildTaskTree(resp.Tasks)
	}
	return resp, nil
}

// buildTaskTree converts a flat list of tasks into a tree structure
// by nesting subtasks under their parents based on ParentID
func buildTaskTree(tasks []Task) []Task {
	taskMap := make(map[string]*Task)
	for i := range tasks {
		task := tasks[i]
		task.Subtasks = []Task{}
		taskMap[task.ID] = &task
	}

	var rootIDs []string
	for _, task := range tasks {
		if task.ParentID == "" {
			rootIDs = append(rootIDs, task.ID)
		} else if parent, ok := taskMap[task.ParentID]; ok {
			child := taskMap[task.ID]
			parent.Subtasks = append(parent.Subtasks, *child)
		} else {
			rootIDs = append(rootIDs, task.ID)
		}
	}

	// Rebuild tree recursively to get nested subtasks
	var buildTree func(id string) Task
	buildTree = func(id string) Task {
		task := *taskMap[id]
		task.Subtasks = []Task{}
		for _, child := range taskMap[id].Subtasks {
			task.Subtasks = append(task.Subtasks, buildTree(child.ID))
		}
		return task
	}

	var roots []Task
	for _, id := range rootIDs {
		roots = append(roots, buildTree(id))
	}

	return roots
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
