package keyring

import "github.com/zalando/go-keyring"

type SystemProvider struct{}

func NewSystemProvider() *SystemProvider {
	return &SystemProvider{}
}

func (s *SystemProvider) Get(service, user string) (string, error) {
	return keyring.Get(service, user)
}

func (s *SystemProvider) Set(service, user, password string) error {
	return keyring.Set(service, user, password)
}

func (s *SystemProvider) Delete(service, user string) error {
	return keyring.Delete(service, user)
}
