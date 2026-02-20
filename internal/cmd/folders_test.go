package cmd

import (
	"testing"

	"github.com/Fire-Dragon-DoL/clickup-cli/internal/config"
)

func TestFoldersListCommand(t *testing.T) {
	cfg = &config.Config{
		SpaceID:       "space123",
		OutputFormat:  "text",
		StrictResolve: false,
	}

	authSource = &mockAuthSource{apiKey: "test-key"}

	cmd := foldersListCmd
	if cmd == nil {
		t.Fatal("foldersListCmd is nil")
	}
	if cmd.Use != "list" {
		t.Errorf("expected Use 'list', got '%s'", cmd.Use)
	}
	if cmd.Short == "" {
		t.Error("expected non-empty Short description")
	}
}

func TestFoldersCommand(t *testing.T) {
	cmd := foldersCmd
	if cmd == nil {
		t.Fatal("foldersCmd is nil")
	}
	if cmd.Use != "folders" {
		t.Errorf("expected Use 'folders', got '%s'", cmd.Use)
	}
	if cmd.Short == "" {
		t.Error("expected non-empty Short description")
	}
}

type mockAuthSource struct {
	apiKey string
	err    error
}

func (m *mockAuthSource) GetAPIKey() (string, error) {
	if m.err != nil {
		return "", m.err
	}
	return m.apiKey, nil
}
