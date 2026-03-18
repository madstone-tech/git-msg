package prompt

// TemplateStore manages prompt templates (embedded defaults + user overrides).
type TemplateStore interface {
	Get(name string) (Template, error)
	List() ([]TemplateEntry, error)
	Save(Template) error
	Delete(name string) error
	// Marshal serialises a Template to bytes (TOML). Used by editor workflows
	// in the application layer so cmd/ need not import a serialisation library.
	Marshal(t Template) ([]byte, error)
	// Unmarshal deserialises bytes into a Template.
	Unmarshal(data []byte) (Template, error)
}
