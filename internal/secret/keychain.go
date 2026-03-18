package secret

import (
	"fmt"
	"os"
	"strings"

	"github.com/99designs/keyring"
)

const serviceName = "git-msg"

// KeychainStore implements SecretStore via the system keychain with env var fallback.
type KeychainStore struct{}

func (k KeychainStore) Get(provider string) (string, error) {
	ring, err := keyring.Open(keyring.Config{
		ServiceName: serviceName,
	})
	if err == nil {
		item, kerr := ring.Get(provider)
		if kerr == nil && len(item.Data) > 0 {
			return string(item.Data), nil
		}
	}

	// Fallback: environment variable GIT_MSG_<PROVIDER>_API_KEY
	envKey := fmt.Sprintf("GIT_MSG_%s_API_KEY", strings.ToUpper(provider))
	if val := os.Getenv(envKey); val != "" {
		return val, nil
	}

	return "", fmt.Errorf("%w for provider %q\n  → set via: git-msg config set provider.name %s && git-msg hook install\n  → or export %s=<your-key>",
		ErrNoCredential, provider, provider, envKey)
}

func (k KeychainStore) Set(provider, key string) error {
	ring, err := keyring.Open(keyring.Config{
		ServiceName: serviceName,
	})
	if err != nil {
		return err
	}
	return ring.Set(keyring.Item{
		Key:  provider,
		Data: []byte(key),
	})
}

func (k KeychainStore) Delete(provider string) error {
	ring, err := keyring.Open(keyring.Config{
		ServiceName: serviceName,
	})
	if err != nil {
		return err
	}
	return ring.Remove(provider)
}
