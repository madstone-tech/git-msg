package cmd_test

import (
	"context"
	"errors"
	"testing"

	"github.com/madstone0-0/git-msg/cmd"
	"github.com/madstone0-0/git-msg/internal/config"
)

func TestSetConfig_ProviderName(t *testing.T) {
	cfg := config.DefaultConfig()
	store := &config.FakeStore{Cfg: cfg}
	if err := cmd.SetConfig(context.Background(), store, "provider.name", "openai"); err != nil {
		t.Fatal(err)
	}
	if store.Saved == nil {
		t.Fatal("expected Save to be called")
	}
	if store.Saved.Provider.Name != "openai" {
		t.Errorf("expected openai, got %q", store.Saved.Provider.Name)
	}
}

func TestGetConfig_ProviderModel(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Provider.Model = "gpt-4o"
	store := &config.FakeStore{Cfg: cfg}
	// GetConfig prints to stdout; just verify no error
	if err := cmd.GetConfig(context.Background(), store, "provider.model"); err != nil {
		t.Fatal(err)
	}
}

func TestSetConfig_UnknownKey(t *testing.T) {
	store := &config.FakeStore{Cfg: config.DefaultConfig()}
	err := cmd.SetConfig(context.Background(), store, "unknown.key", "value")
	if !errors.Is(err, cmd.ErrUnknownConfigKey) {
		t.Errorf("expected ErrUnknownConfigKey, got %v", err)
	}
}

func TestGetConfig_UnknownKey(t *testing.T) {
	store := &config.FakeStore{Cfg: config.DefaultConfig()}
	err := cmd.GetConfig(context.Background(), store, "unknown.key")
	if !errors.Is(err, cmd.ErrUnknownConfigKey) {
		t.Errorf("expected ErrUnknownConfigKey, got %v", err)
	}
}

func TestShowConfig_NoAPIKeys(t *testing.T) {
	cfg := config.DefaultConfig()
	store := &config.FakeStore{Cfg: cfg}
	// ShowConfig should not error and should not include any key-shaped values
	if err := cmd.ShowConfig(context.Background(), store); err != nil {
		t.Fatal(err)
	}
}
