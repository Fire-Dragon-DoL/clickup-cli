package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUpdateTask(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT request, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"id": "task123",
			"name": "Updated Task"
		}`))
	}))
	defer server.Close()
	client := NewClient("key", server.URL)

	payload := map[string]any{
		"name":     "Updated Task",
		"assignee": "user456",
	}

	result, err := UpdateTask(client, "task123", payload)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ID != "task123" {
		t.Errorf("expected ID 'task123', got '%s'", result.ID)
	}
	if result.Name != "Updated Task" {
		t.Errorf("expected name 'Updated Task', got '%s'", result.Name)
	}
}

func TestUpdateTaskPartialUpdate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT request, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"id": "task456",
			"name": "Original Title"
		}`))
	}))
	defer server.Close()
	client := NewClient("key", server.URL)

	payload := map[string]any{
		"status": "completed",
	}

	result, err := UpdateTask(client, "task456", payload)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ID != "task456" {
		t.Errorf("expected ID 'task456', got '%s'", result.ID)
	}
}

func TestUpdateTaskNotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"err": "Task not found", "ECODE": "ITEM_015"}`))
	}))
	defer server.Close()
	client := NewClient("key", server.URL)

	payload := map[string]any{"name": "Updated"}

	_, err := UpdateTask(client, "notfound", payload)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	apiErr, ok := err.(*Error)
	if !ok {
		t.Fatalf("expected *Error, got %T", err)
	}
	if apiErr.StatusCode != http.StatusNotFound {
		t.Errorf("expected status code 404, got %d", apiErr.StatusCode)
	}
}

func TestUpdateTaskAllFields(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"id": "task789",
			"name": "Complete Update"
		}`))
	}))
	defer server.Close()
	client := NewClient("key", server.URL)

	payload := map[string]any{
		"name":        "Complete Update",
		"assignee":    "user789",
		"status":      "in review",
		"priority":    5,
		"due_date":    "2025-12-31",
		"description": "Updated description",
	}

	result, err := UpdateTask(client, "task789", payload)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Name != "Complete Update" {
		t.Errorf("expected name 'Complete Update', got '%s'", result.Name)
	}
}

func TestUpdateTaskBuildsCorrectPath(t *testing.T) {
	var capturedPath string
	var capturedMethod string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedPath = r.RequestURI
		capturedMethod = r.Method
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id": "task123", "name": "Updated"}`))
	}))
	defer server.Close()
	client := NewClient("key", server.URL)

	payload := map[string]any{"name": "Updated"}
	UpdateTask(client, "task123", payload)

	if capturedMethod != "PUT" {
		t.Errorf("expected method PUT, got %s", capturedMethod)
	}
	if capturedPath != "/task/task123" {
		t.Errorf("expected path '/task/task123', got '%s'", capturedPath)
	}
}
