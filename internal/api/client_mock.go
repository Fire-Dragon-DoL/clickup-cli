package api

import "github.com/Fire-Dragon-DoL/clickup-cli/internal/resolver"

// MockCall records a single method invocation
type MockCall struct {
	Method string
	Path   string
	Body   any
}

// MockClient is a test double for API operations
type MockClient struct {
	// Response is returned by MockDo (must match expected type)
	Response any
	// Error is returned by MockDo when set
	Error error
	// Calls records all MockDo invocations
	Calls []MockCall

	// Searcher interface responses
	TasksResponse   []resolver.SearchResult
	TasksError      error
	ListsResponse   []resolver.SearchResult
	ListsError      error
	FoldersResponse []resolver.SearchResult
	FoldersError    error
	UsersResponse   []resolver.SearchResult
	UsersError      error
}

// MockDo simulates the Do function with controllable responses
func MockDo[Req any, Res any](m *MockClient, method, path string, body *Req) (Res, error) {
	var zero Res

	m.Calls = append(m.Calls, MockCall{
		Method: method,
		Path:   path,
		Body:   body,
	})

	if m.Error != nil {
		return zero, m.Error
	}

	if m.Response != nil {
		if res, ok := m.Response.(Res); ok {
			return res, nil
		}
	}

	return zero, nil
}

// SearchTasks implements resolver.Searcher
func (m *MockClient) SearchTasks(query string) ([]resolver.SearchResult, error) {
	if m.TasksError != nil {
		return nil, m.TasksError
	}
	return m.TasksResponse, nil
}

// SearchLists implements resolver.Searcher
func (m *MockClient) SearchLists(query string) ([]resolver.SearchResult, error) {
	if m.ListsError != nil {
		return nil, m.ListsError
	}
	return m.ListsResponse, nil
}

// SearchFolders implements resolver.Searcher
func (m *MockClient) SearchFolders(query string) ([]resolver.SearchResult, error) {
	if m.FoldersError != nil {
		return nil, m.FoldersError
	}
	return m.FoldersResponse, nil
}

// SearchUsers implements resolver.Searcher
func (m *MockClient) SearchUsers(query string) ([]resolver.SearchResult, error) {
	if m.UsersError != nil {
		return nil, m.UsersError
	}
	return m.UsersResponse, nil
}

// Reset clears all recorded calls and responses
func (m *MockClient) Reset() {
	m.Calls = nil
	m.Response = nil
	m.Error = nil
	m.TasksResponse = nil
	m.TasksError = nil
	m.ListsResponse = nil
	m.ListsError = nil
	m.FoldersResponse = nil
	m.FoldersError = nil
	m.UsersResponse = nil
	m.UsersError = nil
}
