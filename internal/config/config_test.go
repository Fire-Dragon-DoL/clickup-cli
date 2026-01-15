package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfig_DefaultValues(t *testing.T) {
	cfg := Load()

	if cfg.SpaceID != "" {
		t.Errorf("expected empty SpaceID, got %q", cfg.SpaceID)
	}
	if cfg.OutputFormat != "text" {
		t.Errorf("expected output format 'text', got %q", cfg.OutputFormat)
	}
	if cfg.StrictResolve != false {
		t.Errorf("expected strict_resolve false, got %v", cfg.StrictResolve)
	}
}

func TestLoadConfig_FromFile(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.json")
	configContent := `{"space_id": "file_space", "output_format": "json", "strict_resolve": true}`
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatal(err)
	}

	cfg := LoadFromFile(configPath)

	if cfg.SpaceID != "file_space" {
		t.Errorf("expected SpaceID 'file_space', got %q", cfg.SpaceID)
	}
	if cfg.OutputFormat != "json" {
		t.Errorf("expected output format 'json', got %q", cfg.OutputFormat)
	}
	if cfg.StrictResolve != true {
		t.Errorf("expected strict_resolve true, got %v", cfg.StrictResolve)
	}
}

func TestLoadConfig_EnvOverridesFile(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.json")
	configContent := `{"space_id": "file_space", "output_format": "text"}`
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatal(err)
	}

	t.Setenv("CLICKUP_SPACE_ID", "env_space")
	t.Setenv("CLICKUP_OUTPUT_FORMAT", "json")

	cfg := LoadFromFile(configPath)

	if cfg.SpaceID != "env_space" {
		t.Errorf("expected SpaceID 'env_space' from env, got %q", cfg.SpaceID)
	}
	if cfg.OutputFormat != "json" {
		t.Errorf("expected output format 'json' from env, got %q", cfg.OutputFormat)
	}
}

func TestLoadConfig_CLIOverridesEnv(t *testing.T) {
	t.Setenv("CLICKUP_SPACE_ID", "env_space")
	t.Setenv("CLICKUP_OUTPUT_FORMAT", "text")

	cfg := Load()
	cfg.ApplyCLIOverrides("cli_space", "json", true)

	if cfg.SpaceID != "cli_space" {
		t.Errorf("expected SpaceID 'cli_space' from CLI, got %q", cfg.SpaceID)
	}
	if cfg.OutputFormat != "json" {
		t.Errorf("expected output format 'json' from CLI, got %q", cfg.OutputFormat)
	}
	if cfg.StrictResolve != true {
		t.Errorf("expected strict_resolve true from CLI, got %v", cfg.StrictResolve)
	}
}

func TestLoadConfig_PriorityChain(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.json")
	configContent := `{"space_id": "file_space", "output_format": "text", "strict_resolve": false}`
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatal(err)
	}

	t.Setenv("CLICKUP_SPACE_ID", "env_space")

	cfg := LoadFromFile(configPath)
	cfg.ApplyCLIOverrides("cli_space", "", false)

	if cfg.SpaceID != "cli_space" {
		t.Errorf("expected SpaceID 'cli_space' from CLI (highest priority), got %q", cfg.SpaceID)
	}
	if cfg.OutputFormat != "text" {
		t.Errorf("expected output format 'text' from file (env not set for this field), got %q", cfg.OutputFormat)
	}
}

func TestLoadConfig_PartialCLIOverrides(t *testing.T) {
	t.Setenv("CLICKUP_SPACE_ID", "env_space")

	cfg := Load()
	cfg.ApplyCLIOverrides("", "json", false)

	if cfg.SpaceID != "env_space" {
		t.Errorf("expected SpaceID 'env_space' from env (CLI not provided), got %q", cfg.SpaceID)
	}
	if cfg.OutputFormat != "json" {
		t.Errorf("expected output format 'json' from CLI, got %q", cfg.OutputFormat)
	}
}
