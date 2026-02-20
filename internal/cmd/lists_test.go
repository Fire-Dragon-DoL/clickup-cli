package cmd

import (
	"testing"

	"github.com/Fire-Dragon-DoL/clickup-cli/internal/config"
)

func TestListsListCommand(t *testing.T) {
	cfg = &config.Config{
		SpaceID:       "space123",
		OutputFormat:  "text",
		StrictResolve: false,
	}

	authSource = &mockAuthSource{apiKey: "test-key"}

	cmd := listsListCmd
	if cmd == nil {
		t.Fatal("listsListCmd is nil")
	}
	if cmd.Use != "list" {
		t.Errorf("expected Use 'list', got '%s'", cmd.Use)
	}
	if cmd.Short == "" {
		t.Error("expected non-empty Short description")
	}
}

func TestListsCommand(t *testing.T) {
	cmd := listsCmd
	if cmd == nil {
		t.Fatal("listsCmd is nil")
	}
	if cmd.Use != "lists" {
		t.Errorf("expected Use 'lists', got '%s'", cmd.Use)
	}
	if cmd.Short == "" {
		t.Error("expected non-empty Short description")
	}
}

func TestListsListCommandHasFolderFlag(t *testing.T) {
	cmd := listsListCmd
	folderFlag := cmd.Flags().Lookup("folder")
	if folderFlag == nil {
		t.Error("expected 'folder' flag to exist")
	}
	if folderFlag.Shorthand != "f" {
		t.Errorf("expected folder flag shorthand 'f', got '%s'", folderFlag.Shorthand)
	}
}
