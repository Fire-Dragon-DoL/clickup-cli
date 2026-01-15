package keyring

const (
	serviceName = "clickup-cli"
	keyUser     = "api_key"
)

type Provider interface {
	Get(service, user string) (string, error)
	Set(service, user, password string) error
	Delete(service, user string) error
}

type Keyring struct {
	provider Provider
}

func New(provider Provider) *Keyring {
	return &Keyring{provider: provider}
}

func (k *Keyring) GetAPIKey() (string, error) {
	return k.provider.Get(serviceName, keyUser)
}

func (k *Keyring) SetAPIKey(apiKey string) error {
	return k.provider.Set(serviceName, keyUser, apiKey)
}

func (k *Keyring) DeleteAPIKey() error {
	return k.provider.Delete(serviceName, keyUser)
}
