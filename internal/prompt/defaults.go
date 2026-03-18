package prompt

import (
	"embed"
	"io/fs"
	"path/filepath"
	"strings"

	toml "github.com/pelletier/go-toml/v2"
)

//go:embed embedded/*.toml
var embeddedFS embed.FS

// loadEmbedded returns all embedded templates keyed by name.
func loadEmbedded() (map[string]Template, error) {
	templates := make(map[string]Template)
	entries, err := fs.ReadDir(embeddedFS, "embedded")
	if err != nil {
		return nil, err
	}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".toml") {
			continue
		}
		data, err := embeddedFS.ReadFile(filepath.Join("embedded", e.Name()))
		if err != nil {
			return nil, err
		}
		var t Template
		if err := toml.Unmarshal(data, &t); err != nil {
			return nil, err
		}
		templates[t.Name] = t
	}
	return templates, nil
}
