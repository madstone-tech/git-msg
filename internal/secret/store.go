package secret

// SecretStore manages API credentials. Keys are NEVER persisted to config files.
type SecretStore interface {
	Get(provider string) (string, error)
	Set(provider, key string) error
	Delete(provider string) error
}
