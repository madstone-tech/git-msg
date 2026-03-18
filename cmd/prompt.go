package cmd

import (
	"context"
	"errors"
	"fmt"

	"github.com/madstone0-0/git-msg/internal/prompt"
	"github.com/madstone0-0/git-msg/internal/ui"
)

// ListPrompts prints all available templates with their source.
func ListPrompts(ctx context.Context, store prompt.TemplateStore) error {
	entries, err := store.List()
	if err != nil {
		return err
	}
	for _, e := range entries {
		fmt.Printf("  %-20s [%s]\n", e.Name, e.Source)
	}
	return nil
}

// ShowPrompt prints the effective template content.
func ShowPrompt(ctx context.Context, store prompt.TemplateStore, name string) error {
	t, err := store.Get(name)
	if err != nil {
		return fmt.Errorf("template %q not found\n  → run: git-msg prompt list", name)
	}
	fmt.Printf("# Template: %s\n# %s\n\n## System\n%s\n\n## User\n%s\n",
		t.Name, t.Description, t.System, t.User)
	return nil
}

// EditPrompt opens the template in $EDITOR and saves the result.
// I/O (editor launch, temp file) is delegated to ui.OpenInEditor.
// Serialisation is delegated to the store's Marshal/Unmarshal methods.
func EditPrompt(ctx context.Context, store prompt.TemplateStore, name string) error {
	t, err := store.Get(name)
	if err != nil {
		return fmt.Errorf("template %q not found\n  → run: git-msg prompt list", name)
	}

	data, err := store.Marshal(t)
	if err != nil {
		return err
	}

	edited, err := ui.OpenInEditor(string(data))
	if err != nil {
		return fmt.Errorf("editor: %w", err)
	}

	updated, err := store.Unmarshal([]byte(edited))
	if err != nil {
		return fmt.Errorf("could not parse edited template: %w", err)
	}
	return store.Save(updated)
}

// ResetPrompt deletes the user override for a template.
func ResetPrompt(ctx context.Context, store prompt.TemplateStore, name string) error {
	err := store.Delete(name)
	if errors.Is(err, prompt.ErrTemplateNotFound) {
		fmt.Printf("Template %q has no user override (already using embedded default).\n", name)
		return nil
	}
	if err != nil {
		return err
	}
	fmt.Printf("Template %q reset to embedded default.\n", name)
	return nil
}
