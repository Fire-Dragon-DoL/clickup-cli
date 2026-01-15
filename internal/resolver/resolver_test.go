package resolver

import (
	"errors"
	"testing"
)

func TestDetectIdentifierType(t *testing.T) {
	tests := []struct {
		input    string
		expected IdentifierType
	}{
		{"abc123", TypeID},
		{"86abc123", TypeID},
		{"Fix login bug", TypeName},
		{"My Task", TypeName},
		{"https://app.clickup.com/t/abc123", TypeURL},
		{"https://app.clickup.com/t/86abc123", TypeURL},
		{"https://app.clickup.com/123/v/li/456", TypeURL},
		{"https://app.clickup.com/123/v/f/789/456", TypeURL},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := DetectIdentifierType(tt.input)

			if result != tt.expected {
				t.Errorf("DetectIdentifierType(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestParseTaskURL(t *testing.T) {
	tests := []struct {
		name      string
		url       string
		wantID    string
		wantError bool
	}{
		{"standard task url", "https://app.clickup.com/t/abc123", "abc123", false},
		{"task url with custom id", "https://app.clickup.com/t/86abc123", "86abc123", false},
		{"task url with workspace", "https://app.clickup.com/t/123456/abc123", "abc123", false},
		{"invalid url", "https://example.com/task/123", "", true},
		{"not a url", "abc123", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := ParseTaskURL(tt.url)

			if tt.wantError && err == nil {
				t.Error("expected error, got nil")
			}
			if !tt.wantError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if id != tt.wantID {
				t.Errorf("ParseTaskURL(%q) = %q, want %q", tt.url, id, tt.wantID)
			}
		})
	}
}

func TestParseListURL(t *testing.T) {
	tests := []struct {
		name      string
		url       string
		wantID    string
		wantError bool
	}{
		{"list url", "https://app.clickup.com/123/v/li/456", "456", false},
		{"list url with longer id", "https://app.clickup.com/123456/v/li/789012", "789012", false},
		{"invalid url", "https://example.com/list/123", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := ParseListURL(tt.url)

			if tt.wantError && err == nil {
				t.Error("expected error, got nil")
			}
			if !tt.wantError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if id != tt.wantID {
				t.Errorf("ParseListURL(%q) = %q, want %q", tt.url, id, tt.wantID)
			}
		})
	}
}

func TestParseFolderURL(t *testing.T) {
	tests := []struct {
		name      string
		url       string
		wantID    string
		wantError bool
	}{
		{"folder url", "https://app.clickup.com/123/v/f/456/789", "456", false},
		{"folder url longer ids", "https://app.clickup.com/123456/v/f/789012/345678", "789012", false},
		{"invalid url", "https://example.com/folder/123", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := ParseFolderURL(tt.url)

			if tt.wantError && err == nil {
				t.Error("expected error, got nil")
			}
			if !tt.wantError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if id != tt.wantID {
				t.Errorf("ParseFolderURL(%q) = %q, want %q", tt.url, id, tt.wantID)
			}
		})
	}
}

func TestResolverResolveTask_ByID(t *testing.T) {
	mock := &MockSearcher{}
	r := New(mock, false)

	taskID, err := r.ResolveTask("abc123")

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if taskID != "abc123" {
		t.Errorf("ResolveTask() = %q, want %q", taskID, "abc123")
	}
}

func TestResolverResolveTask_ByURL(t *testing.T) {
	mock := &MockSearcher{}
	r := New(mock, false)

	taskID, err := r.ResolveTask("https://app.clickup.com/t/xyz789")

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if taskID != "xyz789" {
		t.Errorf("ResolveTask() = %q, want %q", taskID, "xyz789")
	}
}

func TestResolverResolveTask_ByName(t *testing.T) {
	mock := &MockSearcher{
		SearchTasksResult: []SearchResult{{ID: "found123", Name: "Fix login bug"}},
	}
	r := New(mock, false)

	taskID, err := r.ResolveTask("Fix login bug")

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if taskID != "found123" {
		t.Errorf("ResolveTask() = %q, want %q", taskID, "found123")
	}
}

func TestResolverResolveTask_ByName_NotFound(t *testing.T) {
	mock := &MockSearcher{
		SearchTasksResult: []SearchResult{},
	}
	r := New(mock, false)

	_, err := r.ResolveTask("Nonexistent task")

	if err == nil {
		t.Error("expected error, got nil")
	}
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestResolverResolveTask_ByName_Ambiguous_StrictFalse(t *testing.T) {
	mock := &MockSearcher{
		SearchTasksResult: []SearchResult{
			{ID: "first123", Name: "Bug fix"},
			{ID: "second456", Name: "Bug fix v2"},
		},
	}
	r := New(mock, false)

	taskID, err := r.ResolveTask("Bug fix")

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if taskID != "first123" {
		t.Errorf("ResolveTask() = %q, want %q (first match)", taskID, "first123")
	}
}

func TestResolverResolveTask_ByName_Ambiguous_StrictTrue(t *testing.T) {
	mock := &MockSearcher{
		SearchTasksResult: []SearchResult{
			{ID: "first123", Name: "Bug fix"},
			{ID: "second456", Name: "Bug fix v2"},
		},
	}
	r := New(mock, true)

	_, err := r.ResolveTask("Bug fix")

	if err == nil {
		t.Error("expected error, got nil")
	}
	var ambiguousErr *AmbiguousError
	if !errors.As(err, &ambiguousErr) {
		t.Errorf("expected AmbiguousError, got %v", err)
	}
	if len(ambiguousErr.Matches) != 2 {
		t.Errorf("expected 2 matches in error, got %d", len(ambiguousErr.Matches))
	}
}

func TestResolverResolveList_ByID(t *testing.T) {
	mock := &MockSearcher{}
	r := New(mock, false)

	listID, err := r.ResolveList("123456")

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if listID != "123456" {
		t.Errorf("ResolveList() = %q, want %q", listID, "123456")
	}
}

func TestResolverResolveList_ByURL(t *testing.T) {
	mock := &MockSearcher{}
	r := New(mock, false)

	listID, err := r.ResolveList("https://app.clickup.com/123/v/li/456")

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if listID != "456" {
		t.Errorf("ResolveList() = %q, want %q", listID, "456")
	}
}

func TestResolverResolveList_ByName(t *testing.T) {
	mock := &MockSearcher{
		SearchListsResult: []SearchResult{{ID: "list789", Name: "Backlog"}},
	}
	r := New(mock, false)

	listID, err := r.ResolveList("Backlog")

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if listID != "list789" {
		t.Errorf("ResolveList() = %q, want %q", listID, "list789")
	}
}

func TestResolverResolveFolder_ByID(t *testing.T) {
	mock := &MockSearcher{}
	r := New(mock, false)

	folderID, err := r.ResolveFolder("folder123")

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if folderID != "folder123" {
		t.Errorf("ResolveFolder() = %q, want %q", folderID, "folder123")
	}
}

func TestResolverResolveFolder_ByURL(t *testing.T) {
	mock := &MockSearcher{}
	r := New(mock, false)

	folderID, err := r.ResolveFolder("https://app.clickup.com/123/v/f/456/789")

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if folderID != "456" {
		t.Errorf("ResolveFolder() = %q, want %q", folderID, "456")
	}
}

func TestResolverResolveFolder_ByName(t *testing.T) {
	mock := &MockSearcher{
		SearchFoldersResult: []SearchResult{{ID: "folder456", Name: "Engineering"}},
	}
	r := New(mock, false)

	folderID, err := r.ResolveFolder("Engineering")

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if folderID != "folder456" {
		t.Errorf("ResolveFolder() = %q, want %q", folderID, "folder456")
	}
}

func TestResolverResolveUser_ByID(t *testing.T) {
	mock := &MockSearcher{}
	r := New(mock, false)

	userID, err := r.ResolveUser("user123")

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if userID != "user123" {
		t.Errorf("ResolveUser() = %q, want %q", userID, "user123")
	}
}

func TestResolverResolveUser_ByName(t *testing.T) {
	mock := &MockSearcher{
		SearchUsersResult: []SearchResult{{ID: "user456", Name: "John Doe"}},
	}
	r := New(mock, false)

	userID, err := r.ResolveUser("John Doe")

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if userID != "user456" {
		t.Errorf("ResolveUser() = %q, want %q", userID, "user456")
	}
}

func TestResolverResolveUser_ByName_StrictAmbiguous(t *testing.T) {
	mock := &MockSearcher{
		SearchUsersResult: []SearchResult{
			{ID: "user1", Name: "John"},
			{ID: "user2", Name: "John Smith"},
		},
	}
	r := New(mock, true)

	_, err := r.ResolveUser("John")

	if err == nil {
		t.Error("expected error, got nil")
	}
	var ambiguousErr *AmbiguousError
	if !errors.As(err, &ambiguousErr) {
		t.Errorf("expected AmbiguousError, got %v", err)
	}
}

func TestAmbiguousErrorMessage(t *testing.T) {
	err := &AmbiguousError{
		Query: "Bug fix",
		Matches: []SearchResult{
			{ID: "id1", Name: "Bug fix 1"},
			{ID: "id2", Name: "Bug fix 2"},
		},
	}

	msg := err.Error()

	if msg == "" {
		t.Error("expected non-empty error message")
	}
}

func TestMockSearcher_SearchError(t *testing.T) {
	mock := &MockSearcher{
		SearchTasksError: errors.New("api error"),
	}
	r := New(mock, false)

	_, err := r.ResolveTask("Some task")

	if err == nil {
		t.Error("expected error, got nil")
	}
}

// MockSearcher implements Searcher for testing
type MockSearcher struct {
	SearchTasksResult   []SearchResult
	SearchTasksError    error
	SearchListsResult   []SearchResult
	SearchListsError    error
	SearchFoldersResult []SearchResult
	SearchFoldersError  error
	SearchUsersResult   []SearchResult
	SearchUsersError    error
}

func (m *MockSearcher) SearchTasks(query string) ([]SearchResult, error) {
	return m.SearchTasksResult, m.SearchTasksError
}

func (m *MockSearcher) SearchLists(query string) ([]SearchResult, error) {
	return m.SearchListsResult, m.SearchListsError
}

func (m *MockSearcher) SearchFolders(query string) ([]SearchResult, error) {
	return m.SearchFoldersResult, m.SearchFoldersError
}

func (m *MockSearcher) SearchUsers(query string) ([]SearchResult, error) {
	return m.SearchUsersResult, m.SearchUsersError
}
