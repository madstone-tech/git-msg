package secret

// FakeSecretStore is an in-memory SecretStore for testing.
type FakeSecretStore struct {
	Keys   map[string]string
	GetErr error
	SetErr error
}

func NewFakeSecretStore() *FakeSecretStore {
	return &FakeSecretStore{Keys: make(map[string]string)}
}

func (f *FakeSecretStore) Get(provider string) (string, error) {
	if f.GetErr != nil {
		return "", f.GetErr
	}
	v, ok := f.Keys[provider]
	if !ok {
		return "", ErrNoCredential
	}
	return v, nil
}

func (f *FakeSecretStore) Set(provider, key string) error {
	if f.SetErr != nil {
		return f.SetErr
	}
	f.Keys[provider] = key
	return nil
}

func (f *FakeSecretStore) Delete(provider string) error {
	delete(f.Keys, provider)
	return nil
}
