package keyring

import (
	"errors"
	"testing"
)

type mockKeyringProvider struct {
	secrets map[string]string
	err     error
}

func (m *mockKeyringProvider) Get(service, user string) (string, error) {
	if m.err != nil {
		return "", m.err
	}
	key := service + ":" + user
	if val, ok := m.secrets[key]; ok {
		return val, nil
	}
	return "", errors.New("secret not found in keyring")
}

func (m *mockKeyringProvider) Set(service, user, password string) error {
	if m.err != nil {
		return m.err
	}
	key := service + ":" + user
	m.secrets[key] = password
	return nil
}

func (m *mockKeyringProvider) Delete(service, user string) error {
	if m.err != nil {
		return m.err
	}
	key := service + ":" + user
	delete(m.secrets, key)
	return nil
}

func TestGetAPIKey_Success(t *testing.T) {
	mock := &mockKeyringProvider{
		secrets: map[string]string{
			"clickup-cli:api_key": "pk_test_12345",
		},
	}
	kr := New(mock)

	apiKey, err := kr.GetAPIKey()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if apiKey != "pk_test_12345" {
		t.Errorf("expected 'pk_test_12345', got %q", apiKey)
	}
}

func TestGetAPIKey_NotFound(t *testing.T) {
	mock := &mockKeyringProvider{
		secrets: map[string]string{},
	}
	kr := New(mock)

	_, err := kr.GetAPIKey()

	if err == nil {
		t.Fatal("expected error for missing API key")
	}
}

func TestGetAPIKey_KeyringError(t *testing.T) {
	mock := &mockKeyringProvider{
		err: errors.New("keyring access denied"),
	}
	kr := New(mock)

	_, err := kr.GetAPIKey()

	if err == nil {
		t.Fatal("expected error for keyring failure")
	}
}

func TestSetAPIKey_Success(t *testing.T) {
	mock := &mockKeyringProvider{
		secrets: map[string]string{},
	}
	kr := New(mock)

	err := kr.SetAPIKey("pk_new_key")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if mock.secrets["clickup-cli:api_key"] != "pk_new_key" {
		t.Errorf("API key not stored correctly")
	}
}

func TestDeleteAPIKey_Success(t *testing.T) {
	mock := &mockKeyringProvider{
		secrets: map[string]string{
			"clickup-cli:api_key": "pk_test_12345",
		},
	}
	kr := New(mock)

	err := kr.DeleteAPIKey()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := mock.secrets["clickup-cli:api_key"]; ok {
		t.Error("API key should have been deleted")
	}
}
