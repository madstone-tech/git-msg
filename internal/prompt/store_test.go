package prompt_test

import (
	"errors"
	"testing"

	"github.com/madstone-tech/git-msg/internal/prompt"
)

func TestFakeStore_UserOverridesEmbedded(t *testing.T) {
	store := prompt.NewFakeTemplateStore()
	store.Templates["conventional"] = prompt.Template{Name: "conventional", System: "user-system"}

	tmpl, err := store.Get("conventional")
	if err != nil {
		t.Fatal(err)
	}
	if tmpl.System != "user-system" {
		t.Errorf("expected user-system, got %q", tmpl.System)
	}
}

func TestFakeStore_ErrTemplateNotFound(t *testing.T) {
	store := prompt.NewFakeTemplateStore()
	_, err := store.Get("nonexistent")
	if !errors.Is(err, prompt.ErrTemplateNotFound) {
		t.Errorf("expected ErrTemplateNotFound, got %v", err)
	}
}

func TestFakeStore_List(t *testing.T) {
	store := prompt.NewFakeTemplateStore()
	store.Templates["a"] = prompt.Template{Name: "a"}
	store.Templates["b"] = prompt.Template{Name: "b"}
	entries, err := store.List()
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(entries))
	}
}
