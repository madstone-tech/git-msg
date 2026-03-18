package prompt

import "errors"

// ErrTemplateNotFound is returned when a named template does not exist.
var ErrTemplateNotFound = errors.New("template not found")

// TemplateSource indicates whether a template is embedded or user-defined.
type TemplateSource int

const (
	SourceEmbedded TemplateSource = iota
	SourceUser
)

func (s TemplateSource) String() string {
	if s == SourceUser {
		return "user"
	}
	return "embedded"
}

// Template is a named prompt template with system and user sections.
type Template struct {
	Name        string `toml:"name"`
	Description string `toml:"description"`
	System      string `toml:"system"`
	User        string `toml:"user"`
}

// TemplateEntry is a list item returned by TemplateStore.List().
type TemplateEntry struct {
	Name   string
	Source TemplateSource
}

// TemplateVars holds the runtime variables injected into prompt templates.
type TemplateVars struct {
	Diff   string
	Branch string
	Log    string
}
