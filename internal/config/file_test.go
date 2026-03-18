package config_test

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/madstone-tech/git-msg/internal/config"
)

func newTestFileStore(t *testing.T) *config.FileStore {
	t.Helper()
	dir := t.TempDir()
	// Create a FileStore pointing at temp dir using exported constructor
	// We need to expose a way to set path for testing
	// Add a NewFileStoreWithPath constructor
	return config.NewFileStoreWithPath(filepath.Join(dir, "config.toml"))
}

func TestFileStore_LoadSaveRoundtrip(t *testing.T) {
	store := newTestFileStore(t)
	cfg := config.DefaultConfig()
	cfg.Provider.Name = "openai"
	if err := store.Save(cfg); err != nil {
		t.Fatal(err)
	}
	loaded, err := store.Load()
	if err != nil {
		t.Fatal(err)
	}
	if loaded.Provider.Name != "openai" {
		t.Errorf("expected openai, got %q", loaded.Provider.Name)
	}
}

func TestFileStore_Load_ErrNoConfig(t *testing.T) {
	dir := t.TempDir()
	store := config.NewFileStoreWithPath(filepath.Join(dir, "nofile.toml"))
	_, err := store.Load()
	if !errors.Is(err, config.ErrNoConfig) {
		t.Errorf("expected ErrNoConfig, got %v", err)
	}
}
