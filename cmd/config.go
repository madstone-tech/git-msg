package cmd

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/madstone-tech/git-msg/internal/config"
)

// ErrUnknownConfigKey is returned for unrecognized dot-delimited config keys.
var ErrUnknownConfigKey = errors.New("unknown config key")

// SetConfig updates a config value by dot-delimited key.
func SetConfig(ctx context.Context, store config.Store, key, value string) error {
	cfg, err := store.Load()
	if err != nil && !errors.Is(err, config.ErrNoConfig) {
		return err
	}
	if errors.Is(err, config.ErrNoConfig) {
		cfg = config.DefaultConfig()
	}
	switch key {
	case "provider.name":
		cfg.Provider.Name = value
	case "provider.model":
		cfg.Provider.Model = value
	case "ollama.host":
		cfg.Ollama.Host = value
	case "prompt.default":
		cfg.Prompt.Default = value
	case "hook.global":
		b, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("hook.global expects true/false, got %q", value)
		}
		cfg.Hook.Global = b
	default:
		return fmt.Errorf("%w: %q\n  → valid keys: provider.name, provider.model, ollama.host, prompt.default, hook.global",
			ErrUnknownConfigKey, key)
	}
	return store.Save(cfg)
}

// GetConfig prints the value of a config key to stdout.
func GetConfig(ctx context.Context, store config.Store, key string) error {
	cfg, err := store.Load()
	if err != nil {
		return err
	}
	switch key {
	case "provider.name":
		fmt.Println(cfg.Provider.Name)
	case "provider.model":
		fmt.Println(cfg.Provider.Model)
	case "ollama.host":
		fmt.Println(cfg.Ollama.Host)
	case "prompt.default":
		fmt.Println(cfg.Prompt.Default)
	case "hook.global":
		fmt.Println(cfg.Hook.Global)
	default:
		return fmt.Errorf("%w: %q", ErrUnknownConfigKey, key)
	}
	return nil
}

// ShowConfig prints the full resolved config.
// Serialisation is delegated to store.Format() — cmd/ has no format coupling.
func ShowConfig(ctx context.Context, store config.Store) error {
	data, err := store.Format()
	if err != nil {
		return err
	}
	fmt.Print(string(data))
	return nil
}
