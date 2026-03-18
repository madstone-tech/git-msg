package cmd

import (
	"context"
	"errors"
	"fmt"

	"github.com/madstone-tech/git-msg/internal/config"
	"github.com/madstone-tech/git-msg/internal/secret"
	"github.com/madstone-tech/git-msg/internal/ui"
)

// Version is set at build time via -ldflags "-X main.Version=vX.Y.Z".
var Version = "dev"

// contextKey is an unexported type for context keys in this package.
type contextKey int

const (
	configKey  contextKey = iota
	secretsKey contextKey = iota
)

// WithConfig stores a loaded Config in the context.
func WithConfig(ctx context.Context, cfg config.Config) context.Context {
	return context.WithValue(ctx, configKey, cfg)
}

// ConfigFromContext retrieves the Config stored by WithConfig.
func ConfigFromContext(ctx context.Context) (config.Config, bool) {
	cfg, ok := ctx.Value(configKey).(config.Config)
	return cfg, ok
}

// WithSecrets stores a SecretStore in the context.
func WithSecrets(ctx context.Context, s secret.SecretStore) context.Context {
	return context.WithValue(ctx, secretsKey, s)
}

// SecretsFromContext retrieves the SecretStore stored by WithSecrets.
func SecretsFromContext(ctx context.Context) (secret.SecretStore, bool) {
	s, ok := ctx.Value(secretsKey).(secret.SecretStore)
	return s, ok
}

// EnsureConfig loads config, running the first-run wizard if none exists.
// When the wizard completes, EnsureConfig persists the config and stores the
// credential — the ui layer returns plain values and has no side-effects.
func EnsureConfig(ctx context.Context, cfgStore config.Store, secrets secret.SecretStore) (config.Config, error) {
	cfg, err := cfgStore.Load()
	if err == nil {
		return cfg, nil
	}
	if !errors.Is(err, config.ErrNoConfig) {
		return config.Config{}, err
	}

	// First run: collect values via wizard (no I/O in ui layer).
	result, err := ui.RunSetupWizard()
	if err != nil {
		if errors.Is(err, ui.ErrWizardAborted) {
			return config.Config{}, fmt.Errorf("setup cancelled — run git-msg again to configure")
		}
		return config.Config{}, err
	}

	// Persist credential (application layer responsibility).
	if result.APIKey != "" {
		if err := secrets.Set(result.Provider, result.APIKey); err != nil {
			return config.Config{}, fmt.Errorf("storing API key: %w", err)
		}
	}

	// Persist config (application layer responsibility).
	cfg = config.DefaultConfig()
	cfg.Provider.Name = result.Provider
	cfg.Provider.Model = result.Model
	if err := cfgStore.Save(cfg); err != nil {
		return config.Config{}, fmt.Errorf("saving config: %w", err)
	}

	fmt.Printf("\n  Config saved. Provider: %s  Model: %s\n\n", result.Provider, result.Model)
	return cfg, nil
}
