package cmd_test

import (
	"errors"
	"testing"

	"github.com/madstone-tech/git-msg/cmd"
	"github.com/madstone-tech/git-msg/internal/config"
	"github.com/madstone-tech/git-msg/internal/secret"
)

func TestNewLLMProvider_Anthropic(t *testing.T) {
	cfg := config.Config{Provider: config.ProviderConfig{Name: "anthropic", Model: "claude-haiku-4-5"}}
	secrets := secret.NewFakeSecretStore()
	secrets.Keys["anthropic"] = "test-key"
	p, err := cmd.NewLLMProvider(cfg, secrets)
	if err != nil {
		t.Fatal(err)
	}
	if p == nil {
		t.Fatal("expected provider")
	}
}

func TestNewLLMProvider_Ollama_NoKey(t *testing.T) {
	cfg := config.Config{
		Provider: config.ProviderConfig{Name: "ollama", Model: "llama3"},
		Ollama:   config.OllamaConfig{Host: "http://localhost:11434"},
	}
	p, err := cmd.NewLLMProvider(cfg, secret.NewFakeSecretStore())
	if err != nil {
		t.Fatal(err)
	}
	if p == nil {
		t.Fatal("expected provider")
	}
}

func TestNewLLMProvider_Unknown(t *testing.T) {
	cfg := config.Config{Provider: config.ProviderConfig{Name: "unknown"}}
	_, err := cmd.NewLLMProvider(cfg, secret.NewFakeSecretStore())
	if err == nil {
		t.Fatal("expected error for unknown provider")
	}
}

func TestNewLLMProvider_MissingKey(t *testing.T) {
	cfg := config.Config{Provider: config.ProviderConfig{Name: "openai", Model: "gpt-4o-mini"}}
	_, err := cmd.NewLLMProvider(cfg, secret.NewFakeSecretStore())
	if err == nil {
		t.Fatal("expected credential error")
	}
	if !errors.Is(err, secret.ErrNoCredential) {
		t.Errorf("expected ErrNoCredential, got %v", err)
	}
}
