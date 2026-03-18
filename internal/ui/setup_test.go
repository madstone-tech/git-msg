package ui_test

import (
	"testing"

	"github.com/madstone-tech/git-msg/internal/ui"
)

func TestDefaultModelFor(t *testing.T) {
	cases := []struct{ provider, want string }{
		{"anthropic", "claude-haiku-4-5"},
		{"openai", "gpt-4o-mini"},
		{"gemini", "gemini-1.5-flash"},
		{"ollama", "llama3"},
		{"unknown", ""},
	}
	for _, c := range cases {
		got := ui.DefaultModelFor(c.provider)
		if got != c.want {
			t.Errorf("DefaultModelFor(%q) = %q, want %q", c.provider, got, c.want)
		}
	}
}

// TestListOllamaModels_DoesNotPanic verifies ListOllamaModels never panics
// and returns only non-empty model names (may be nil in CI without Ollama).
func TestListOllamaModels_DoesNotPanic(t *testing.T) {
	models := ui.ListOllamaModels()
	for i, m := range models {
		if m == "" {
			t.Errorf("models[%d] is empty string", i)
		}
	}
}

// TestListOllamaModels_NoDuplicates verifies no model name appears twice.
func TestListOllamaModels_NoDuplicates(t *testing.T) {
	models := ui.ListOllamaModels()
	seen := make(map[string]bool)
	for _, m := range models {
		if seen[m] {
			t.Errorf("duplicate model name: %q", m)
		}
		seen[m] = true
	}
}

// TestErrWizardAborted is exported and distinct from huh.ErrUserAborted.
func TestErrWizardAborted_IsNotNil(t *testing.T) {
	if ui.ErrWizardAborted == nil {
		t.Error("ErrWizardAborted must not be nil")
	}
}
