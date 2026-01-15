package cmd

import (
	"testing"

	"github.com/Fire-Dragon-DoL/clickup-cli/internal/api"
	"github.com/Fire-Dragon-DoL/clickup-cli/internal/config"
	"github.com/Fire-Dragon-DoL/clickup-cli/internal/output"
)

func TestTasksCmd(t *testing.T) {
	cmd := tasksCmd
	if cmd == nil {
		t.Fatal("tasksCmd is nil")
	}
	if cmd.Use != "tasks" {
		t.Errorf("expected Use 'tasks', got '%s'", cmd.Use)
	}
	if cmd.Short == "" {
		t.Error("expected non-empty Short description")
	}
}

func TestTasksListCmd(t *testing.T) {
	cmd := tasksListCmd
	if cmd == nil {
		t.Fatal("tasksListCmd is nil")
	}
	if cmd.Use != "list" {
		t.Errorf("expected Use 'list', got '%s'", cmd.Use)
	}
	if cmd.Short == "" {
		t.Error("expected non-empty Short description")
	}
}

func TestTasksListCmdHasListFlag(t *testing.T) {
	cmd := tasksListCmd
	listFlag := cmd.Flags().Lookup("list")
	if listFlag == nil {
		t.Error("expected 'list' flag to exist")
	}
	if listFlag.Shorthand != "l" {
		t.Errorf("expected list flag shorthand 'l', got '%s'", listFlag.Shorthand)
	}
}

func TestTasksListCmdHasRecursiveFlag(t *testing.T) {
	cmd := tasksListCmd
	recursiveFlag := cmd.Flags().Lookup("recursive")
	if recursiveFlag == nil {
		t.Error("expected 'recursive' flag to exist")
	}
	if recursiveFlag.Shorthand != "r" {
		t.Errorf("expected recursive flag shorthand 'r', got '%s'", recursiveFlag.Shorthand)
	}
}

func TestFormatTasksListView(t *testing.T) {
	tasks := []api.Task{
		{
			ID:   "task1",
			Name: "Test Task 1",
			Assignee: &api.User{
				ID:       "user1",
				Username: "john",
			},
			Status: &struct {
				ID      string `json:"id"`
				Status  string `json:"status"`
				Color   string `json:"color"`
				OrderBy int    `json:"orderby"`
			}{
				ID:     "status1",
				Status: "open",
			},
			Priority: &struct {
				ID       int    `json:"id"`
				Priority string `json:"priority"`
				Color    string `json:"color"`
				OrderBy  int    `json:"orderby"`
			}{
				ID:       1,
				Priority: "high",
			},
		},
	}

	cfg = &config.Config{OutputFormat: "text"}
	formatter = output.NewFormatter("text")

	formatted, err := formatTasksListView(tasks)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if formatted == "" {
		t.Error("expected formatted output to be non-empty")
	}
}

func TestFormatTasksListViewWithNilFields(t *testing.T) {
	tasks := []api.Task{
		{
			ID:       "task1",
			Name:     "Test Task",
			Assignee: nil,
			Status:   nil,
			Priority: nil,
		},
	}

	cfg = &config.Config{OutputFormat: "text"}
	formatter = output.NewFormatter("text")

	formatted, err := formatTasksListView(tasks)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if formatted == "" {
		t.Error("expected formatted output to be non-empty")
	}
}

func TestTasksShowCmd(t *testing.T) {
	cmd := tasksShowCmd
	if cmd == nil {
		t.Fatal("tasksShowCmd is nil")
	}
	if cmd.Use != "show <task-id|name|url>" {
		t.Errorf("expected Use 'show <task-id|name|url>', got '%s'", cmd.Use)
	}
	if cmd.Short == "" {
		t.Error("expected non-empty Short description")
	}
}

func TestFormatTaskDetailsView(t *testing.T) {
	task := api.Task{
		ID:   "task123",
		Name: "Test Task",
		Status: &struct {
			ID      string `json:"id"`
			Status  string `json:"status"`
			Color   string `json:"color"`
			OrderBy int    `json:"orderby"`
		}{
			ID:     "status1",
			Status: "open",
		},
		Priority: &struct {
			ID       int    `json:"id"`
			Priority string `json:"priority"`
			Color    string `json:"color"`
			OrderBy  int    `json:"orderby"`
		}{
			ID:       1,
			Priority: "high",
		},
		Assignee: &api.User{
			ID:       "user1",
			Username: "john",
		},
		DueDate:     "2025-12-31",
		Description: "Test description",
	}
	comments := []api.Comment{
		{
			ID:          "comment1",
			TextContent: "First comment",
			User: api.User{
				ID:       "user1",
				Username: "john",
			},
			DateCreated: "2025-01-14",
		},
	}

	cfg = &config.Config{OutputFormat: "text"}
	formatter = output.NewFormatter("text")

	formatted, err := formatTaskDetailsView(task, comments...)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if formatted == "" {
		t.Error("expected formatted output to be non-empty")
	}
}

func TestFormatTaskDetailsViewWithoutComments(t *testing.T) {
	task := api.Task{
		ID:   "task123",
		Name: "Test Task",
	}

	cfg = &config.Config{OutputFormat: "text"}
	formatter = output.NewFormatter("text")

	formatted, err := formatTaskDetailsView(task)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if formatted == "" {
		t.Error("expected formatted output to be non-empty")
	}
}

func TestTasksDeleteCmd(t *testing.T) {
	cmd := tasksDeleteCmd

	if cmd == nil {
		t.Fatal("tasksDeleteCmd is nil")
	}
	if cmd.Use != "delete <task-id|name|url>" {
		t.Errorf("expected Use 'delete <task-id|name|url>', got '%s'", cmd.Use)
	}
	if cmd.Short == "" {
		t.Error("expected non-empty Short description")
	}
}

func TestTasksArchiveCmd(t *testing.T) {
	cmd := tasksArchiveCmd

	if cmd == nil {
		t.Fatal("tasksArchiveCmd is nil")
	}
	if cmd.Use != "archive <task-id|name|url>" {
		t.Errorf("expected Use 'archive <task-id|name|url>', got '%s'", cmd.Use)
	}
	if cmd.Short == "" {
		t.Error("expected non-empty Short description")
	}
}
