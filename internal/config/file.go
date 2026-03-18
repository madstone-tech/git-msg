package config

import (
	"errors"
	"os"
	"path/filepath"

	toml "github.com/pelletier/go-toml/v2"

	"github.com/madstone0-0/git-msg/internal/dirs"
)

// ErrNoConfig is returned by Load when no config file exists (triggers first-run wizard).
var ErrNoConfig = errors.New("no configuration file found")

// FileStore implements Store using a TOML file resolved via internal/dirs
// (~/.config/mdstn/git-msg/config.toml by default).
type FileStore struct {
	path string
}

// NewFileStore returns a FileStore using the XDG config path.
func NewFileStore() (*FileStore, error) {
	path, err := dirs.ConfigFile()
	if err != nil {
		return nil, err
	}
	return &FileStore{path: path}, nil
}

// NewFileStoreWithPath returns a FileStore with a custom path (for testing).
func NewFileStoreWithPath(path string) *FileStore {
	return &FileStore{path: path}
}

func (s *FileStore) Load() (Config, error) {
	data, err := os.ReadFile(s.path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return Config{}, ErrNoConfig
		}
		return Config{}, err
	}
	var cfg Config
	if err := toml.Unmarshal(data, &cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func (s *FileStore) Save(cfg Config) error {
	if err := os.MkdirAll(filepath.Dir(s.path), 0750); err != nil {
		return err
	}
	data, err := toml.Marshal(cfg)
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0644)
}

func (s *FileStore) Format() ([]byte, error) {
	cfg, err := s.Load()
	if err != nil {
		return nil, err
	}
	return toml.Marshal(cfg)
}
