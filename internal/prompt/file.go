package prompt

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	toml "github.com/pelletier/go-toml/v2"

	"github.com/madstone-tech/git-msg/internal/dirs"
)

// FileStore implements TemplateStore with user overrides at
// ~/.config/mdstn/git-msg/prompts/ (XDG convention).
type FileStore struct {
	dir string
}

// NewFileStore returns a FileStore using the XDG prompts path.
func NewFileStore() (*FileStore, error) {
	d, err := dirs.PromptsDir()
	if err != nil {
		return nil, err
	}
	return &FileStore{dir: d}, nil
}

func (s *FileStore) Get(name string) (Template, error) {
	// User override takes precedence
	userPath := filepath.Join(s.dir, name+".toml")
	data, err := os.ReadFile(userPath)
	if err == nil {
		var t Template
		if err := toml.Unmarshal(data, &t); err != nil {
			return Template{}, err
		}
		return t, nil
	}
	// Fallback to embedded
	embedded, err := loadEmbedded()
	if err != nil {
		return Template{}, err
	}
	if t, ok := embedded[name]; ok {
		return t, nil
	}
	return Template{}, ErrTemplateNotFound
}

func (s *FileStore) List() ([]TemplateEntry, error) {
	embedded, err := loadEmbedded()
	if err != nil {
		return nil, err
	}
	seen := make(map[string]TemplateSource)
	for name := range embedded {
		seen[name] = SourceEmbedded
	}

	// User files override source tag
	entries, _ := os.ReadDir(s.dir)
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".toml") {
			name := strings.TrimSuffix(e.Name(), ".toml")
			seen[name] = SourceUser
		}
	}

	var result []TemplateEntry
	for name, src := range seen {
		result = append(result, TemplateEntry{Name: name, Source: src})
	}
	return result, nil
}

func (s *FileStore) Save(t Template) error {
	if err := os.MkdirAll(s.dir, 0755); err != nil {
		return err
	}
	data, err := toml.Marshal(t)
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(s.dir, t.Name+".toml"), data, 0644)
}

func (s *FileStore) Delete(name string) error {
	err := os.Remove(filepath.Join(s.dir, name+".toml"))
	if errors.Is(err, os.ErrNotExist) {
		return ErrTemplateNotFound
	}
	return err
}

func (s *FileStore) Marshal(t Template) ([]byte, error) {
	return toml.Marshal(t)
}

func (s *FileStore) Unmarshal(data []byte) (Template, error) {
	var t Template
	if err := toml.Unmarshal(data, &t); err != nil {
		return Template{}, err
	}
	return t, nil
}
