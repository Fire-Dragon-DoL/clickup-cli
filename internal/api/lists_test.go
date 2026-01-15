package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetLists(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected method GET, got %s", r.Method)
		}
		if r.URL.Path != "/folder/456/list" {
			t.Errorf("expected path /folder/456/list, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(ListsResponse{
			Lists: []List{
				{ID: "list1", Name: "Backlog"},
				{ID: "list2", Name: "In Progress"},
			},
		})
	}))
	defer server.Close()

	client := NewClient("test-key", server.URL)

	arrange := struct {
		folderID string
	}{
		folderID: "456",
	}

	result, err := GetLists(client, arrange.folderID)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("expected 2 lists, got %d", len(result))
	}
	if result[0].ID != "list1" {
		t.Errorf("expected first list ID 'list1', got '%s'", result[0].ID)
	}
	if result[0].Name != "Backlog" {
		t.Errorf("expected first list name 'Backlog', got '%s'", result[0].Name)
	}
	if result[1].ID != "list2" {
		t.Errorf("expected second list ID 'list2', got '%s'", result[1].ID)
	}
	if result[1].Name != "In Progress" {
		t.Errorf("expected second list name 'In Progress', got '%s'", result[1].Name)
	}
}

func TestGetListsEmpty(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(ListsResponse{Lists: []List{}})
	}))
	defer server.Close()

	client := NewClient("test-key", server.URL)

	result, err := GetLists(client, "456")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected 0 lists, got %d", len(result))
	}
}

func TestGetListsError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"err":   "Folder not found",
			"ECODE": "ITEM_015",
		})
	}))
	defer server.Close()

	client := NewClient("test-key", server.URL)

	_, err := GetLists(client, "notfound")

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
