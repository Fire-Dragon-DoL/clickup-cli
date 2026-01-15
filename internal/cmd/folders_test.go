package cmd

import (
	"testing"

	"github.com/Fire-Dragon-DoL/clickup-cli/internal/config"
	"github.com/Fire-Dragon-DoL/clickup-cli/internal/keyring"
)

func TestFoldersListCommand(t *testing.T) {
	cfg = &config.Config{
		SpaceID:       "space123",
		OutputFormat:  "text",
		StrictResolve: false,
	}

	mockProvider := &mockKeyringProvider{
		apiKey: "test-key",
	}
	kr = keyring.New(mockProvider)

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

type mockKeyringProvider struct {
	apiKey string
	err    error
}

func (m *mockKeyringProvider) Get(service, user string) (string, error) {
	if m.err != nil {
		return "", m.err
	}
	return m.apiKey, nil
}

func (m *mockKeyringProvider) Set(service, user, password string) error {
	return nil
}

func (m *mockKeyringProvider) Delete(service, user string) error {
	return nil
}
