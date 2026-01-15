package output

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/Fire-Dragon-DoL/clickup-cli/internal/api"
)

type Formatter struct {
	format string
}

func NewFormatter(format string) *Formatter {
	if format == "" || (format != "text" && format != "json") {
		format = "text"
	}
	return &Formatter{format: format}
}

func (f *Formatter) Format(data any) (string, error) {
	if f.format == "json" {
		return f.formatJSON(data)
	}
	return f.formatText(data)
}

func (f *Formatter) formatJSON(data any) (string, error) {
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (f *Formatter) formatText(data any) (string, error) {
	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Slice {
		return f.formatSlice(v)
	}
	return f.formatStruct(v)
}

func (f *Formatter) formatSlice(v reflect.Value) (string, error) {
	var lines []string
	for i := 0; i < v.Len(); i++ {
		item := v.Index(i)
		line, err := f.formatStruct(item)
		if err != nil {
			return "", err
		}
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n"), nil
}

func (f *Formatter) formatStruct(v reflect.Value) (string, error) {
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return fmt.Sprintf("%v", v.Interface()), nil
	}

	t := v.Type()
	var parts []string
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)
		if !field.IsExported() {
			continue
		}
		strVal := fmt.Sprintf("%v", value.Interface())
		if strVal != "" {
			parts = append(parts, fmt.Sprintf("%s: %s", field.Name, strVal))
		}
	}
	return strings.Join(parts, " | "), nil
}

func (f *Formatter) FormatTaskList(tasks []api.Task, recursive bool) (string, error) {
	if f.format == "json" {
		b, err := json.MarshalIndent(tasks, "", "  ")
		if err != nil {
			return "", err
		}
		return string(b), nil
	}
	result, err := f.formatTaskListText(tasks, recursive)
	return result, err
}

func (f *Formatter) formatTaskListText(tasks []api.Task, recursive bool) (string, error) {
	var lines []string
	for _, task := range tasks {
		lines = append(lines, f.formatTaskWithIndent(task, 0)...)
		if recursive {
			lines = append(lines, f.formatSubtasksText(task.Subtasks, 1)...)
		}
	}
	return strings.Join(lines, "\n"), nil
}

func (f *Formatter) formatTaskWithIndent(task api.Task, indent int) []string {
	prefix := strings.Repeat("  ", indent)
	status := ""
	if task.Status != nil {
		status = task.Status.Status
	}
	priority := ""
	if task.Priority != nil {
		priority = task.Priority.Priority
	}
	assignee := ""
	if task.Assignee != nil {
		assignee = task.Assignee.Username
	}

	line := fmt.Sprintf("%s%s | %s | %s | %s | %s", prefix, task.ID, task.Name, assignee, status, priority)
	return []string{line}
}

func (f *Formatter) formatSubtasksText(tasks []api.Task, indent int) []string {
	var lines []string
	for _, task := range tasks {
		lines = append(lines, f.formatTaskWithIndent(task, indent)...)
		lines = append(lines, f.formatSubtasksText(task.Subtasks, indent+1)...)
	}
	return lines
}
