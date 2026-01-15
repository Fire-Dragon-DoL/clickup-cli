package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetFolders(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected method GET, got %s", r.Method)
		}
		if r.URL.Path != "/space/123/folder" {
			t.Errorf("expected path /space/123/folder, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(FoldersResponse{
			Folders: []Folder{
				{ID: "456", Name: "My Folder"},
				{ID: "789", Name: "Another Folder"},
			},
		})
	}))
	defer server.Close()

	client := NewClient("test-key", server.URL, "")

	arrange := struct {
		spaceID string
	}{
		spaceID: "123",
	}

	result, err := GetFolders(client, arrange.spaceID)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("expected 2 folders, got %d", len(result))
	}
	if result[0].ID != "456" {
		t.Errorf("expected first folder ID '456', got '%s'", result[0].ID)
	}
	if result[0].Name != "My Folder" {
		t.Errorf("expected first folder name 'My Folder', got '%s'", result[0].Name)
	}
	if result[1].ID != "789" {
		t.Errorf("expected second folder ID '789', got '%s'", result[1].ID)
	}
	if result[1].Name != "Another Folder" {
		t.Errorf("expected second folder name 'Another Folder', got '%s'", result[1].Name)
	}
}

func TestGetFoldersEmpty(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(FoldersResponse{Folders: []Folder{}})
	}))
	defer server.Close()

	client := NewClient("test-key", server.URL, "")

	result, err := GetFolders(client, "123")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected 0 folders, got %d", len(result))
	}
}

func TestGetFoldersError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"err":   "Unauthorized",
			"ECODE": "OAUTH_017",
		})
	}))
	defer server.Close()

	client := NewClient("bad-key", server.URL, "")

	_, err := GetFolders(client, "123")

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
