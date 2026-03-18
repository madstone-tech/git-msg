package cmd

import (
	"fmt"

	"github.com/madstone-tech/git-msg/internal/config"
	"github.com/madstone-tech/git-msg/internal/llm"
	"github.com/madstone-tech/git-msg/internal/secret"
)

// NewLLMProvider constructs the appropriate llm.Provider from config and
// credentials. This is application-layer wiring — it does not belong inside
// the llm infrastructure package.
func NewLLMProvider(cfg config.Config, secrets secret.SecretStore) (llm.Provider, error) {
	name := cfg.Provider.Name
	model := cfg.Provider.Model

	switch name {
	case "openai":
		key, err := secrets.Get(name)
		if err != nil {
			return nil, err
		}
		return llm.NewOpenAIProvider(model, key), nil
	case "anthropic":
		key, err := secrets.Get(name)
		if err != nil {
			return nil, err
		}
		return llm.NewAnthropicProvider(model, key), nil
	case "gemini":
		key, err := secrets.Get(name)
		if err != nil {
			return nil, err
		}
		return llm.NewGeminiProvider(model, key), nil
	case "ollama":
		host := cfg.Ollama.Host
		if host == "" {
			host = "http://localhost:11434"
		}
		return llm.NewOllamaProvider(model, host), nil
	default:
		return nil, fmt.Errorf("unknown provider %q; supported: openai, anthropic, gemini, ollama", name)
	}
}
