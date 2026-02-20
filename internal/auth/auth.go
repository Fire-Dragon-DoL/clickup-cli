package auth

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os/exec"
	"strings"

	"github.com/Fire-Dragon-DoL/clickup-cli/internal/keyring"
)

const defaultChezmoiKey = "secrets.clickup.key"

type Source interface {
	GetAPIKey() (string, error)
}

func ParseSource(rawURL string) (Source, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("invalid auth source URL: %w", err)
	}

	switch u.Scheme {
	case "keyring":
		return keyring.New(keyring.NewSystemProvider()), nil
	case "chezmoi":
		keyPath := u.Host
		if u.Path != "" {
			keyPath = keyPath + u.Path
		}
		if keyPath == "" {
			keyPath = defaultChezmoiKey
		}
		return &ChezmoiSource{KeyPath: keyPath, execCommand: defaultExecCommand}, nil
	default:
		return nil, fmt.Errorf("unknown auth source scheme: %q", u.Scheme)
	}
}

type execCommandFunc func(name string, args ...string) ([]byte, error)

func defaultExecCommand(name string, args ...string) ([]byte, error) {
	return exec.Command(name, args...).Output()
}

type ChezmoiSource struct {
	KeyPath     string
	execCommand execCommandFunc
}

func (c *ChezmoiSource) GetAPIKey() (string, error) {
	output, err := c.execCommand("chezmoi", "data", "--format", "json")
	if err != nil {
		return "", fmt.Errorf("failed to run chezmoi data: %w", err)
	}

	var data map[string]any
	if err := json.Unmarshal(output, &data); err != nil {
		return "", fmt.Errorf("failed to parse chezmoi data: %w", err)
	}

	parts := strings.Split(c.KeyPath, ".")
	current := any(data)
	for _, part := range parts {
		m, ok := current.(map[string]any)
		if !ok {
			return "", fmt.Errorf("key %q not found in chezmoi data", c.KeyPath)
		}
		current, ok = m[part]
		if !ok {
			return "", fmt.Errorf("key %q not found in chezmoi data", c.KeyPath)
		}
	}

	key, ok := current.(string)
	if !ok {
		return "", fmt.Errorf("key %q is not a string in chezmoi data", c.KeyPath)
	}

	return key, nil
}
