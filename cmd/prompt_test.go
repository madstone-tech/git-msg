package cmd_test

import (
	"context"
	"testing"

	"github.com/madstone-tech/git-msg/cmd"
	"github.com/madstone-tech/git-msg/internal/prompt"
)

func TestListPrompts(t *testing.T) {
	store := prompt.NewFakeTemplateStore()
	store.Templates["conventional"] = prompt.Template{Name: "conventional"}
	store.Templates["custom"] = prompt.Template{Name: "custom"}
	// Should not error
	if err := cmd.ListPrompts(context.Background(), store); err != nil {
		t.Fatal(err)
	}
}

func TestShowPrompt(t *testing.T) {
	store := prompt.NewFakeTemplateStore()
	store.Templates["conventional"] = prompt.Template{
		Name: "conventional", System: "sys", User: "usr",
	}
	if err := cmd.ShowPrompt(context.Background(), store, "conventional"); err != nil {
		t.Fatal(err)
	}
}

func TestShowPrompt_NotFound(t *testing.T) {
	store := prompt.NewFakeTemplateStore()
	err := cmd.ShowPrompt(context.Background(), store, "nonexistent")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestResetPrompt_Exists(t *testing.T) {
	store := prompt.NewFakeTemplateStore()
	store.Templates["conventional"] = prompt.Template{Name: "conventional"}
	if err := cmd.ResetPrompt(context.Background(), store, "conventional"); err != nil {
		t.Fatal(err)
	}
	if _, exists := store.Templates["conventional"]; exists {
		t.Error("template should be deleted after reset")
	}
}

func TestResetPrompt_NotFound(t *testing.T) {
	store := prompt.NewFakeTemplateStore()
	// Should exit 0 with informational message, not error
	err := cmd.ResetPrompt(context.Background(), store, "nonexistent")
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}
