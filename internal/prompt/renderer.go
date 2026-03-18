package prompt

import (
	"fmt"

	ason "github.com/madstone-tech/ason/pkg"
)

// Renderer renders Template text using ason Jinja2-style engine.
type Renderer struct {
	engine ason.Engine
}

// NewRenderer returns a Renderer with the default ason engine.
func NewRenderer() *Renderer {
	return &Renderer{engine: ason.NewDefaultEngine()}
}

// Render returns the rendered system and user strings with vars injected.
func (r *Renderer) Render(t Template, vars TemplateVars) (system, user string, err error) {
	m := map[string]interface{}{
		"diff":   vars.Diff,
		"branch": vars.Branch,
		"log":    vars.Log,
	}
	system, err = r.engine.Render(t.System, m)
	if err != nil {
		return "", "", fmt.Errorf("rendering system prompt: %w", err)
	}
	user, err = r.engine.Render(t.User, m)
	if err != nil {
		return "", "", fmt.Errorf("rendering user prompt: %w", err)
	}
	return system, user, nil
}
