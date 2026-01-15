package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetTasks(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"tasks": [
				{"id": "task1", "name": "Task 1"},
				{"id": "task2", "name": "Task 2"}
			]
		}`))
	}))
	defer server.Close()
	client := NewClient("key", server.URL, "")

	result, err := GetTasks(client, "list123", false)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Tasks) != 2 {
		t.Errorf("expected 2 tasks, got %d", len(result.Tasks))
	}
	if result.Tasks[0].ID != "task1" {
		t.Errorf("expected first task ID 'task1', got '%s'", result.Tasks[0].ID)
	}
	if result.Tasks[0].Name != "Task 1" {
		t.Errorf("expected first task name 'Task 1', got '%s'", result.Tasks[0].Name)
	}
}

func TestGetTasksBuildsCorrectPath(t *testing.T) {
	var capturedPath string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedPath = r.RequestURI
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"tasks": []}`))
	}))
	defer server.Close()
	client := NewClient("key", server.URL, "")

	GetTasks(client, "list123", false)

	expectedPath := "/list/list123/task?archived=false"
	if capturedPath != expectedPath {
		t.Errorf("expected path '%s', got '%s'", expectedPath, capturedPath)
	}
}

func TestGetTasksWithRecursiveFlag(t *testing.T) {
	var capturedPath string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedPath = r.RequestURI
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"tasks": []}`))
	}))
	defer server.Close()
	client := NewClient("key", server.URL, "")

	GetTasks(client, "list456", true)

	if capturedPath != "/list/list456/task?archived=false&subtasks=true" {
		t.Errorf("expected path with subtasks=true, got '%s'", capturedPath)
	}
}

func TestGetTasksError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"err": "List not found", "ECODE": "LIST_001"}`))
	}))
	defer server.Close()
	client := NewClient("key", server.URL, "")

	_, err := GetTasks(client, "invalid", false)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestCreateTask(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"id": "task-123",
			"name": "New Task",
			"list": "list-456"
		}`))
	}))
	defer server.Close()
	client := NewClient("key", server.URL, "")

	payload := map[string]any{
		"name":    "New Task",
		"list_id": "list-456",
	}
	result, err := CreateTask(client, payload)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ID != "task-123" {
		t.Errorf("expected ID 'task-123', got '%s'", result.ID)
	}
	if result.Name != "New Task" {
		t.Errorf("expected Name 'New Task', got '%s'", result.Name)
	}
	if result.ListID != "list-456" {
		t.Errorf("expected ListID 'list-456', got '%s'", result.ListID)
	}
}

func TestCreateTaskCallsCorrectPath(t *testing.T) {
	var capturedPath string
	var capturedMethod string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedPath = r.RequestURI
		capturedMethod = r.Method
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id": "task-123", "name": "New Task", "list": "list-456"}`))
	}))
	defer server.Close()
	client := NewClient("key", server.URL, "")

	payload := map[string]any{
		"name":    "New Task",
		"list_id": "list-456",
	}
	CreateTask(client, payload)

	if capturedMethod != "POST" {
		t.Errorf("expected method POST, got %s", capturedMethod)
	}
	if capturedPath != "/list/list-456/task" {
		t.Errorf("expected path '/list/list-456/task', got '%s'", capturedPath)
	}
}

func TestCreateTaskWithOptionalFields(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"id": "task-123",
			"name": "New Task with Details",
			"list": "list-456",
			"description": "Test description",
			"parent": "parent-789"
		}`))
	}))
	defer server.Close()
	client := NewClient("key", server.URL, "")

	payload := map[string]any{
		"name":        "New Task with Details",
		"list_id":     "list-456",
		"description": "Test description",
		"priority":    1,
		"status":      "to do",
		"due_date":    1234567890000,
		"parent":      "parent-789",
	}
	result, err := CreateTask(client, payload)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Name != "New Task with Details" {
		t.Errorf("expected Name 'New Task with Details', got '%s'", result.Name)
	}
	if result.Description != "Test description" {
		t.Errorf("expected Description 'Test description', got '%s'", result.Description)
	}
	if result.ParentID != "parent-789" {
		t.Errorf("expected ParentID 'parent-789', got '%s'", result.ParentID)
	}
}

func TestCreateTaskMissingListID(t *testing.T) {
	client := NewClient("key", "http://localhost", "")

	payload := map[string]any{
		"name": "Task without list",
	}
	_, err := CreateTask(client, payload)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "list_id is required" {
		t.Errorf("expected error 'list_id is required', got '%s'", err.Error())
	}
}

func TestCreateTaskMissingNameField(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id": "task-123", "name": "", "list": "list-456"}`))
	}))
	defer server.Close()
	client := NewClient("key", server.URL, "")

	payload := map[string]any{
		"list_id": "list-456",
	}
	result, err := CreateTask(client, payload)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if result.ID != "task-123" {
		t.Errorf("expected ID 'task-123', got '%s'", result.ID)
	}
}

func TestGetTask(t *testing.T) {
	task := Task{
		ID:   "task123",
		Name: "Test Task",
		Status: &struct {
			ID      string `json:"id"`
			Status  string `json:"status"`
			Color   string `json:"color"`
			OrderBy int    `json:"orderby"`
		}{
			ID:     "status1",
			Status: "open",
			Color:  "#FF0000",
		},
		Priority: &struct {
			ID       int    `json:"id"`
			Priority string `json:"priority"`
			Color    string `json:"color"`
			OrderBy  int    `json:"orderby"`
		}{
			ID:       1,
			Priority: "high",
			Color:    "#FF0000",
		},
		DueDate:     "2025-12-31",
		Description: "Test task description",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(task)
	}))
	defer server.Close()
	client := NewClient("key", server.URL, "")

	result, err := GetTask(client, "task123")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ID != "task123" {
		t.Errorf("expected task ID 'task123', got '%s'", result.ID)
	}
	if result.Name != "Test Task" {
		t.Errorf("expected task name 'Test Task', got '%s'", result.Name)
	}
	if result.Status.Status != "open" {
		t.Errorf("expected status 'open', got '%s'", result.Status.Status)
	}
	if result.Description != "Test task description" {
		t.Errorf("expected description 'Test task description', got '%s'", result.Description)
	}
}

func TestGetTaskNotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"err":   "Task not found",
			"ECODE": "ITEM_015",
		})
	}))
	defer server.Close()
	client := NewClient("key", server.URL, "")

	_, err := GetTask(client, "notfound")

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	apiErr, ok := err.(*Error)
	if !ok {
		t.Fatalf("expected *Error, got %T", err)
	}
	if apiErr.StatusCode != http.StatusNotFound {
		t.Errorf("expected status code %d, got %d", http.StatusNotFound, apiErr.StatusCode)
	}
}

func TestGetTaskComments(t *testing.T) {
	commentsResp := CommentsResponse{
		Comments: []Comment{
			{
				ID:          "comment1",
				TextContent: "First comment",
				User: User{
					ID:       "user1",
					Username: "john",
					Email:    "john@example.com",
				},
				DateCreated: "2025-01-14",
			},
			{
				ID:          "comment2",
				TextContent: "Second comment",
				User: User{
					ID:       "user2",
					Username: "jane",
					Email:    "jane@example.com",
				},
				DateCreated: "2025-01-15",
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(commentsResp)
	}))
	defer server.Close()
	client := NewClient("key", server.URL, "")

	comments, err := GetTaskComments(client, "task123")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(comments) != 2 {
		t.Errorf("expected 2 comments, got %d", len(comments))
	}
	if comments[0].ID != "comment1" {
		t.Errorf("expected first comment ID 'comment1', got '%s'", comments[0].ID)
	}
	if comments[0].TextContent != "First comment" {
		t.Errorf("expected first comment text 'First comment', got '%s'", comments[0].TextContent)
	}
	if comments[1].ID != "comment2" {
		t.Errorf("expected second comment ID 'comment2', got '%s'", comments[1].ID)
	}
}

func TestGetTaskCommentsEmpty(t *testing.T) {
	commentsResp := CommentsResponse{
		Comments: []Comment{},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(commentsResp)
	}))
	defer server.Close()
	client := NewClient("key", server.URL, "")

	comments, err := GetTaskComments(client, "task123")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(comments) != 0 {
		t.Errorf("expected 0 comments, got %d", len(comments))
	}
}

func TestGetTasksWithSubtasksRecursive(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"tasks": [
				{
					"id": "parent1",
					"name": "Parent Task",
					"subtasks": [
						{
							"id": "child1",
							"name": "Child Task",
							"subtasks": [
								{"id": "grandchild1", "name": "Grandchild Task"}
							]
						}
					]
				}
			]
		}`))
	}))
	defer server.Close()
	client := NewClient("key", server.URL, "")

	result, err := GetTasks(client, "list123", true)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Tasks) != 1 {
		t.Errorf("expected 1 task, got %d", len(result.Tasks))
	}
	parent := result.Tasks[0]
	if len(parent.Subtasks) != 1 {
		t.Errorf("expected 1 child task, got %d", len(parent.Subtasks))
	}
	child := parent.Subtasks[0]
	if len(child.Subtasks) != 1 {
		t.Errorf("expected 1 grandchild task, got %d", len(child.Subtasks))
	}
	grandchild := child.Subtasks[0]
	if grandchild.ID != "grandchild1" {
		t.Errorf("expected grandchild ID 'grandchild1', got '%s'", grandchild.ID)
	}
}

func TestGetTasksWithMultipleRootsAndChildren(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"tasks": [
				{
					"id": "root1",
					"name": "Root 1",
					"subtasks": [
						{"id": "child1.1", "name": "Child 1.1"},
						{"id": "child1.2", "name": "Child 1.2"}
					]
				},
				{
					"id": "root2",
					"name": "Root 2",
					"subtasks": []
				}
			]
		}`))
	}))
	defer server.Close()
	client := NewClient("key", server.URL, "")

	result, err := GetTasks(client, "list123", true)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Tasks) != 2 {
		t.Errorf("expected 2 root tasks, got %d", len(result.Tasks))
	}
	root1 := result.Tasks[0]
	if len(root1.Subtasks) != 2 {
		t.Errorf("expected 2 subtasks for root1, got %d", len(root1.Subtasks))
	}
	root2 := result.Tasks[1]
	if len(root2.Subtasks) != 0 {
		t.Errorf("expected 0 subtasks for root2, got %d", len(root2.Subtasks))
	}
}
