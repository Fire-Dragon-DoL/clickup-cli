package api

import (
	"errors"
	"testing"

	"github.com/Fire-Dragon-DoL/clickup-cli/internal/resolver"
)

func TestMockDo_ReturnsConfiguredResponse(t *testing.T) {
	type Response struct {
		ID   string
		Name string
	}
	mock := &MockClient{
		Response: Response{ID: "123", Name: "Test"},
	}

	result, err := MockDo[any, Response](mock, "GET", "/test", nil)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ID != "123" {
		t.Errorf("expected ID '123', got '%s'", result.ID)
	}
}

func TestMockDo_ReturnsConfiguredError(t *testing.T) {
	expectedErr := errors.New("api failure")
	mock := &MockClient{
		Error: expectedErr,
	}

	_, err := MockDo[any, map[string]string](mock, "GET", "/test", nil)

	if err != expectedErr {
		t.Errorf("expected error '%v', got '%v'", expectedErr, err)
	}
}

func TestMockDo_RecordsCalls(t *testing.T) {
	type Request struct {
		Name string
	}
	mock := &MockClient{}
	body := &Request{Name: "test"}

	MockDo[Request, any](mock, "POST", "/tasks", body)
	MockDo[any, any](mock, "GET", "/lists", nil)

	if len(mock.Calls) != 2 {
		t.Fatalf("expected 2 calls, got %d", len(mock.Calls))
	}
	if mock.Calls[0].Method != "POST" {
		t.Errorf("expected method 'POST', got '%s'", mock.Calls[0].Method)
	}
	if mock.Calls[0].Path != "/tasks" {
		t.Errorf("expected path '/tasks', got '%s'", mock.Calls[0].Path)
	}
	if mock.Calls[1].Method != "GET" {
		t.Errorf("expected method 'GET', got '%s'", mock.Calls[1].Method)
	}
}

func TestMockClient_SearchTasks(t *testing.T) {
	mock := &MockClient{
		TasksResponse: []resolver.SearchResult{
			{ID: "task1", Name: "Task One"},
		},
	}

	results, err := mock.SearchTasks("Task")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].ID != "task1" {
		t.Errorf("expected ID 'task1', got '%s'", results[0].ID)
	}
}

func TestMockClient_SearchTasksError(t *testing.T) {
	expectedErr := errors.New("search failed")
	mock := &MockClient{
		TasksError: expectedErr,
	}

	_, err := mock.SearchTasks("Task")

	if err != expectedErr {
		t.Errorf("expected error '%v', got '%v'", expectedErr, err)
	}
}

func TestMockClient_SearchLists(t *testing.T) {
	mock := &MockClient{
		ListsResponse: []resolver.SearchResult{
			{ID: "list1", Name: "List One"},
		},
	}

	results, err := mock.SearchLists("List")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
}

func TestMockClient_SearchFolders(t *testing.T) {
	mock := &MockClient{
		FoldersResponse: []resolver.SearchResult{
			{ID: "folder1", Name: "Folder One"},
		},
	}

	results, err := mock.SearchFolders("Folder")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
}

func TestMockClient_SearchUsers(t *testing.T) {
	mock := &MockClient{
		UsersResponse: []resolver.SearchResult{
			{ID: "user1", Name: "User One"},
		},
	}

	results, err := mock.SearchUsers("User")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
}

func TestMockClient_Reset(t *testing.T) {
	mock := &MockClient{
		Response:        "test",
		Error:           errors.New("test"),
		Calls:           []MockCall{{Method: "GET"}},
		TasksResponse:   []resolver.SearchResult{{ID: "1"}},
		TasksError:      errors.New("test"),
		ListsResponse:   []resolver.SearchResult{{ID: "1"}},
		ListsError:      errors.New("test"),
		FoldersResponse: []resolver.SearchResult{{ID: "1"}},
		FoldersError:    errors.New("test"),
		UsersResponse:   []resolver.SearchResult{{ID: "1"}},
		UsersError:      errors.New("test"),
	}

	mock.Reset()

	if mock.Response != nil {
		t.Error("Response should be nil after reset")
	}
	if mock.Error != nil {
		t.Error("Error should be nil after reset")
	}
	if mock.Calls != nil {
		t.Error("Calls should be nil after reset")
	}
	if mock.TasksResponse != nil {
		t.Error("TasksResponse should be nil after reset")
	}
}

func TestMockClient_ImplementsSearcher(t *testing.T) {
	var _ resolver.Searcher = (*MockClient)(nil)
}
