package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDeleteTaskSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected method DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/task/abc123" {
			t.Errorf("expected path /task/abc123, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	}))
	defer server.Close()
	client := NewClient("key", server.URL, "")

	err := DeleteTask(client, "abc123")

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestDeleteTaskNotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"err": "Task not found", "ECODE": "ITEM_015"}`))
	}))
	defer server.Close()
	client := NewClient("key", server.URL, "")

	err := DeleteTask(client, "notfound")

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	apiErr, ok := err.(*Error)
	if !ok {
		t.Fatalf("expected *Error, got %T", err)
	}
	if apiErr.StatusCode != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", apiErr.StatusCode)
	}
	if apiErr.Message != "Task not found" {
		t.Errorf("expected message 'Task not found', got '%s'", apiErr.Message)
	}
}

func TestArchiveTaskSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected method PUT, got %s", r.Method)
		}
		if r.URL.Path != "/task/abc123/archive" {
			t.Errorf("expected path /task/abc123/archive, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	}))
	defer server.Close()
	client := NewClient("key", server.URL, "")

	err := ArchiveTask(client, "abc123")

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestArchiveTaskNotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"err": "Task not found", "ECODE": "ITEM_015"}`))
	}))
	defer server.Close()
	client := NewClient("key", server.URL, "")

	err := ArchiveTask(client, "notfound")

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	apiErr, ok := err.(*Error)
	if !ok {
		t.Fatalf("expected *Error, got %T", err)
	}
	if apiErr.StatusCode != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", apiErr.StatusCode)
	}
	if apiErr.Message != "Task not found" {
		t.Errorf("expected message 'Task not found', got '%s'", apiErr.Message)
	}
}
