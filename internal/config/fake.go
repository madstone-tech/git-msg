package config

// FakeStore is an in-memory Store for testing.
type FakeStore struct {
	Cfg     Config
	LoadErr error
	SaveErr error
	Saved   *Config
}

func (f *FakeStore) Load() (Config, error) {
	return f.Cfg, f.LoadErr
}

func (f *FakeStore) Save(c Config) error {
	f.Saved = &c
	return f.SaveErr
}

func (f *FakeStore) Format() ([]byte, error) {
	return []byte("[provider]\n  name = \"" + f.Cfg.Provider.Name + "\"\n"), nil
}
