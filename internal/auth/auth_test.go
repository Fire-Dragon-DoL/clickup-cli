package auth

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestParseSource_KeyringDefault(t *testing.T) {
	src, err := ParseSource("keyring://")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if src == nil {
		t.Fatal("expected non-nil source")
	}
}

func TestParseSource_ChezmoiWithKey(t *testing.T) {
	src, err := ParseSource("chezmoi://secrets.clickup.key")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	cs, ok := src.(*ChezmoiSource)
	if !ok {
		t.Fatal("expected *ChezmoiSource")
	}
	if cs.KeyPath != "secrets.clickup.key" {
		t.Errorf("expected key path 'secrets.clickup.key', got %q", cs.KeyPath)
	}
}

func TestParseSource_ChezmoiDefaultKey(t *testing.T) {
	src, err := ParseSource("chezmoi://")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	cs, ok := src.(*ChezmoiSource)
	if !ok {
		t.Fatal("expected *ChezmoiSource")
	}
	if cs.KeyPath != "secrets.clickup.key" {
		t.Errorf("expected default key path 'secrets.clickup.key', got %q", cs.KeyPath)
	}
}

func TestParseSource_ChezmoiNestedKey(t *testing.T) {
	src, err := ParseSource("chezmoi://deep/nested.key")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	cs, ok := src.(*ChezmoiSource)
	if !ok {
		t.Fatal("expected *ChezmoiSource")
	}
	if cs.KeyPath != "deep/nested.key" {
		t.Errorf("expected key path 'deep/nested.key', got %q", cs.KeyPath)
	}
}

func TestParseSource_UnknownScheme(t *testing.T) {
	_, err := ParseSource("vault://secret/key")

	if err == nil {
		t.Fatal("expected error for unknown scheme")
	}
}

func TestParseSource_InvalidURL(t *testing.T) {
	_, err := ParseSource("://bad")

	if err == nil {
		t.Fatal("expected error for invalid URL")
	}
}

func TestChezmoiSource_GetAPIKey(t *testing.T) {
	data := map[string]any{
		"secrets": map[string]any{
			"clickup": map[string]any{
				"key": "test-api-key-123",
			},
		},
	}
	jsonBytes, _ := json.Marshal(data)

	src := &ChezmoiSource{
		KeyPath: "secrets.clickup.key",
		execCommand: func(name string, args ...string) ([]byte, error) {
			return jsonBytes, nil
		},
	}

	apiKey, err := src.GetAPIKey()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if apiKey != "test-api-key-123" {
		t.Errorf("expected 'test-api-key-123', got %q", apiKey)
	}
}

func TestChezmoiSource_GetAPIKey_MissingKey(t *testing.T) {
	data := map[string]any{
		"other": "value",
	}
	jsonBytes, _ := json.Marshal(data)

	src := &ChezmoiSource{
		KeyPath: "secrets.clickup.key",
		execCommand: func(name string, args ...string) ([]byte, error) {
			return jsonBytes, nil
		},
	}

	_, err := src.GetAPIKey()

	if err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestChezmoiSource_GetAPIKey_ExecError(t *testing.T) {
	src := &ChezmoiSource{
		KeyPath: "secrets.clickup.key",
		execCommand: func(name string, args ...string) ([]byte, error) {
			return nil, fmt.Errorf("chezmoi not found")
		},
	}

	_, err := src.GetAPIKey()

	if err == nil {
		t.Fatal("expected error when exec fails")
	}
}

func TestChezmoiSource_GetAPIKey_NotAString(t *testing.T) {
	data := map[string]any{
		"secrets": map[string]any{
			"clickup": map[string]any{
				"key": 12345,
			},
		},
	}
	jsonBytes, _ := json.Marshal(data)

	src := &ChezmoiSource{
		KeyPath: "secrets.clickup.key",
		execCommand: func(name string, args ...string) ([]byte, error) {
			return jsonBytes, nil
		},
	}

	_, err := src.GetAPIKey()

	if err == nil {
		t.Fatal("expected error when key is not a string")
	}
}
