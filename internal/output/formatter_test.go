package output

import (
	"strings"
	"testing"
)

type sampleTask struct {
	ID       string
	Title    string
	Assignee string
	Status   string
	Priority string
}

func TestTextFormatter_SingleItem(t *testing.T) {
	task := sampleTask{
		ID:       "abc123",
		Title:    "Fix login bug",
		Assignee: "John Doe",
		Status:   "in progress",
		Priority: "high",
	}
	formatter := NewFormatter("text")

	output, err := formatter.Format(task)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(output, "abc123") {
		t.Error("output should contain task ID")
	}
	if !strings.Contains(output, "Fix login bug") {
		t.Error("output should contain task title")
	}
}

func TestJSONFormatter_SingleItem(t *testing.T) {
	task := sampleTask{
		ID:       "abc123",
		Title:    "Fix login bug",
		Assignee: "John Doe",
		Status:   "in progress",
		Priority: "high",
	}
	formatter := NewFormatter("json")

	output, err := formatter.Format(task)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(output, `"ID": "abc123"`) {
		t.Error("JSON output should contain task ID field")
	}
	if !strings.Contains(output, `"Title": "Fix login bug"`) {
		t.Error("JSON output should contain task title field")
	}
}

func TestTextFormatter_List(t *testing.T) {
	tasks := []sampleTask{
		{ID: "abc123", Title: "Task 1", Status: "open"},
		{ID: "def456", Title: "Task 2", Status: "closed"},
	}
	formatter := NewFormatter("text")

	output, err := formatter.Format(tasks)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(output, "abc123") {
		t.Error("output should contain first task ID")
	}
	if !strings.Contains(output, "def456") {
		t.Error("output should contain second task ID")
	}
}

func TestJSONFormatter_List(t *testing.T) {
	tasks := []sampleTask{
		{ID: "abc123", Title: "Task 1", Status: "open"},
		{ID: "def456", Title: "Task 2", Status: "closed"},
	}
	formatter := NewFormatter("json")

	output, err := formatter.Format(tasks)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.HasPrefix(output, "[") {
		t.Error("JSON list output should start with '['")
	}
	if !strings.Contains(output, `"abc123"`) {
		t.Error("JSON output should contain first task ID")
	}
}

func TestFormatter_DefaultsToText(t *testing.T) {
	formatter := NewFormatter("")
	task := sampleTask{ID: "abc123", Title: "Test"}

	output, err := formatter.Format(task)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if output == "" {
		t.Error("formatter should produce output")
	}
}

func TestFormatter_InvalidFormat(t *testing.T) {
	formatter := NewFormatter("xml")

	task := sampleTask{ID: "abc123", Title: "Test"}
	output, err := formatter.Format(task)

	if err != nil {
		t.Fatalf("should not error, got: %v", err)
	}
	if output == "" {
		t.Error("should fall back to text format")
	}
}

func TestJSONFormatter_PrettyPrint(t *testing.T) {
	task := sampleTask{
		ID:    "abc123",
		Title: "Test",
	}
	formatter := NewFormatter("json")

	output, err := formatter.Format(task)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(output, "\n") {
		t.Error("JSON output should be pretty-printed with newlines")
	}
}
