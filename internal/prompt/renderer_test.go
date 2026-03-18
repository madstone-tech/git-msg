package prompt_test

import (
	"strings"
	"testing"

	"github.com/madstone0-0/git-msg/internal/prompt"
)

func TestRenderer_InjectsVars(t *testing.T) {
	r := prompt.NewRenderer()
	tmpl := prompt.Template{
		System: "system: {{ diff }}",
		User:   "branch={{ branch }} log={{ log }}",
	}
	vars := prompt.TemplateVars{Diff: "mydiff", Branch: "main", Log: "abc123"}
	system, user, err := r.Render(tmpl, vars)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(system, "mydiff") {
		t.Errorf("system missing diff: %q", system)
	}
	if !strings.Contains(user, "main") {
		t.Errorf("user missing branch: %q", user)
	}
	if !strings.Contains(user, "abc123") {
		t.Errorf("user missing log: %q", user)
	}
}
