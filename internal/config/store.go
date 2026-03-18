package config

// Store persists and retrieves application configuration.
type Store interface {
	Load() (Config, error)
	Save(Config) error
	// Format serialises the current config to human-readable bytes (TOML).
	// Callers (e.g. ShowConfig) use this instead of importing a serialisation
	// library directly, keeping the cmd layer free of format coupling.
	Format() ([]byte, error)
}
