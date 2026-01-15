package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewClient(t *testing.T) {
	client := NewClient("test-api-key", "", "")

	if client.apiKey != "test-api-key" {
		t.Errorf("expected apiKey 'test-api-key', got '%s'", client.apiKey)
	}
	if client.baseURL != DefaultBaseURL {
		t.Errorf("expected baseURL '%s', got '%s'", DefaultBaseURL, client.baseURL)
	}
}

func TestNewClientCustomBaseURL(t *testing.T) {
	customURL := "https://custom.clickup.com/api/v2"

	client := NewClient("test-api-key", customURL, "")

	if client.baseURL != customURL {
		t.Errorf("expected baseURL '%s', got '%s'", customURL, client.baseURL)
	}
}

func TestDoInjectsAuthHeader(t *testing.T) {
	var receivedAuth string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedAuth = r.Header.Get("Authorization")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}))
	defer server.Close()
	client := NewClient("my-secret-key", server.URL, "")

	_, err := Do[any, map[string]string](client, http.MethodGet, "/test", nil)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if receivedAuth != "my-secret-key" {
		t.Errorf("expected Authorization header 'my-secret-key', got '%s'", receivedAuth)
	}
}

func TestDoSendsRequestBody(t *testing.T) {
	type RequestBody struct {
		Name string `json:"name"`
	}
	var receivedBody RequestBody
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&receivedBody)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}))
	defer server.Close()
	client := NewClient("key", server.URL, "")

	_, err := Do[RequestBody, map[string]string](client, http.MethodPost, "/test", &RequestBody{Name: "test-name"})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if receivedBody.Name != "test-name" {
		t.Errorf("expected body name 'test-name', got '%s'", receivedBody.Name)
	}
}

func TestDoReturnsResponse(t *testing.T) {
	type Response struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Response{ID: "123", Name: "Test Task"})
	}))
	defer server.Close()
	client := NewClient("key", server.URL, "")

	result, err := Do[any, Response](client, http.MethodGet, "/task/123", nil)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ID != "123" {
		t.Errorf("expected ID '123', got '%s'", result.ID)
	}
	if result.Name != "Test Task" {
		t.Errorf("expected Name 'Test Task', got '%s'", result.Name)
	}
}

func TestDoReturnsErrorOnHTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"err":   "Invalid request",
			"ECODE": "BAD_REQUEST",
		})
	}))
	defer server.Close()
	client := NewClient("key", server.URL, "")

	_, err := Do[any, map[string]string](client, http.MethodGet, "/test", nil)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	apiErr, ok := err.(*Error)
	if !ok {
		t.Fatalf("expected *Error, got %T", err)
	}
	if apiErr.Message != "Invalid request" {
		t.Errorf("expected message 'Invalid request', got '%s'", apiErr.Message)
	}
	if apiErr.Code != "BAD_REQUEST" {
		t.Errorf("expected code 'BAD_REQUEST', got '%s'", apiErr.Code)
	}
	if apiErr.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status code %d, got %d", http.StatusBadRequest, apiErr.StatusCode)
	}
}

func TestDoReturnsErrorOnUnauthorized(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"err":   "Token invalid",
			"ECODE": "OAUTH_017",
		})
	}))
	defer server.Close()
	client := NewClient("bad-key", server.URL, "")

	_, err := Do[any, map[string]string](client, http.MethodGet, "/test", nil)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	apiErr, ok := err.(*Error)
	if !ok {
		t.Fatalf("expected *Error, got %T", err)
	}
	if apiErr.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected status code %d, got %d", http.StatusUnauthorized, apiErr.StatusCode)
	}
}

func TestDoReturnsErrorOnNotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"err":   "Task not found",
			"ECODE": "ITEM_015",
		})
	}))
	defer server.Close()
	client := NewClient("key", server.URL, "")

	_, err := Do[any, map[string]string](client, http.MethodGet, "/task/notfound", nil)

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

func TestDoHandlesNetworkError(t *testing.T) {
	client := NewClient("key", "http://localhost:1", "")

	_, err := Do[any, map[string]string](client, http.MethodGet, "/test", nil)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestDoSetsContentTypeHeader(t *testing.T) {
	var receivedContentType string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedContentType = r.Header.Get("Content-Type")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}))
	defer server.Close()
	client := NewClient("key", server.URL, "")

	_, err := Do[map[string]string, map[string]string](client, http.MethodPost, "/test", &map[string]string{"key": "value"})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if receivedContentType != "application/json" {
		t.Errorf("expected Content-Type 'application/json', got '%s'", receivedContentType)
	}
}

func TestErrorImplementsError(t *testing.T) {
	err := &Error{
		StatusCode: 400,
		Code:       "BAD_REQUEST",
		Message:    "Invalid input",
	}

	expected := "clickup api error (400): Invalid input [BAD_REQUEST]"
	if err.Error() != expected {
		t.Errorf("expected '%s', got '%s'", expected, err.Error())
	}
}
