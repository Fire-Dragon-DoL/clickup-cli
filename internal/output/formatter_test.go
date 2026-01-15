package output

import (
	"strings"
	"testing"

	"github.com/Fire-Dragon-DoL/clickup-cli/internal/api"
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

func TestFormatTaskList_Flat(t *testing.T) {
	tasks := []api.Task{
		{
			ID:   "task1",
			Name: "Task 1",
		},
		{
			ID:   "task2",
			Name: "Task 2",
		},
	}
	formatter := NewFormatter("text")

	output, err := formatter.FormatTaskList(tasks, false)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(output, "task1") {
		t.Error("output should contain first task ID")
	}
	if !strings.Contains(output, "task2") {
		t.Error("output should contain second task ID")
	}
	if !strings.Contains(output, "Task 1") {
		t.Error("output should contain first task name")
	}
}

func TestFormatTaskList_Recursive(t *testing.T) {
	tasks := []api.Task{
		{
			ID:   "parent",
			Name: "Parent Task",
			Subtasks: []api.Task{
				{
					ID:   "child",
					Name: "Child Task",
				},
			},
		},
	}
	formatter := NewFormatter("text")

	output, err := formatter.FormatTaskList(tasks, true)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(output, "parent") {
		t.Error("output should contain parent task ID")
	}
	if !strings.Contains(output, "child") {
		t.Error("output should contain child task ID")
	}
	lines := strings.Split(output, "\n")
	if len(lines) < 2 {
		t.Error("output should have at least 2 lines for parent and child")
	}
	parentIdx := -1
	childIdx := -1
	for i, line := range lines {
		if strings.Contains(line, "parent") {
			parentIdx = i
		}
		if strings.Contains(line, "child") {
			childIdx = i
		}
	}
	if parentIdx >= childIdx {
		t.Error("parent task should appear before child task")
	}
	if childIdx > -1 && !strings.HasPrefix(lines[childIdx], "  ") {
		t.Error("child task should be indented")
	}
}

func TestFormatTaskList_JSON(t *testing.T) {
	tasks := []api.Task{
		{
			ID:   "task1",
			Name: "Task 1",
		},
	}
	formatter := NewFormatter("json")

	output, err := formatter.FormatTaskList(tasks, false)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.HasPrefix(output, "[") {
		t.Error("JSON output should start with '['")
	}
	if !strings.Contains(output, `"id": "task1"`) {
		t.Error("JSON output should contain task ID field")
	}
}

func TestFormatTaskList_HierarchicalIndentation(t *testing.T) {
	tasks := []api.Task{
		{
			ID:   "root",
			Name: "Root",
			Subtasks: []api.Task{
				{
					ID:   "child",
					Name: "Child",
					Subtasks: []api.Task{
						{
							ID:   "grandchild",
							Name: "Grandchild",
						},
					},
				},
			},
		},
	}
	formatter := NewFormatter("text")

	output, err := formatter.FormatTaskList(tasks, true)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(output, "\n")

	rootLine := ""
	childLine := ""
	grandchildLine := ""
	for _, line := range lines {
		if strings.Contains(line, "root") && !strings.HasPrefix(line, " ") {
			rootLine = line
		} else if strings.Contains(line, "child") && strings.HasPrefix(line, "  ") && !strings.HasPrefix(line, "    ") {
			childLine = line
		} else if strings.Contains(line, "grandchild") && strings.HasPrefix(line, "    ") {
			grandchildLine = line
		}
	}

	if rootLine == "" {
		t.Error("should have non-indented root task")
	}
	if childLine == "" {
		t.Error("should have 2-space indented child task")
	}
	if grandchildLine == "" {
		t.Error("should have 4-space indented grandchild task")
	}
}
