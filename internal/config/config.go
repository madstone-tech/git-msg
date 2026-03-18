package config

// Config is the full application configuration stored at
// ~/.config/mdstn/git-msg/config.toml
// ($XDG_CONFIG_HOME/mdstn/git-msg/config.toml when XDG_CONFIG_HOME is set).
type Config struct {
	Provider ProviderConfig `toml:"provider"`
	Ollama   OllamaConfig   `toml:"ollama"`
	Prompt   PromptConfig   `toml:"prompt"`
	Hook     HookConfig     `toml:"hook"`
}

type ProviderConfig struct {
	Name  string `toml:"name"`
	Model string `toml:"model"`
}

type OllamaConfig struct {
	Host string `toml:"host"`
}

type PromptConfig struct {
	Default string `toml:"default"`
}

type HookConfig struct {
	Global bool `toml:"global"`
}

// DefaultConfig returns the factory defaults used when no config file exists.
func DefaultConfig() Config {
	return Config{
		Provider: ProviderConfig{
			Name:  "anthropic",
			Model: "claude-haiku-4-5",
		},
		Ollama: OllamaConfig{
			Host: "http://localhost:11434",
		},
		Prompt: PromptConfig{
			Default: "conventional",
		},
		Hook: HookConfig{
			Global: false,
		},
	}
}
