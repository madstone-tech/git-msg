package prompt

// FakeTemplateStore is an in-memory TemplateStore for testing.
type FakeTemplateStore struct {
	Templates map[string]Template
	GetErr    error
	ListErr   error
	SaveErr   error
	DeleteErr error
}

func NewFakeTemplateStore() *FakeTemplateStore {
	return &FakeTemplateStore{Templates: make(map[string]Template)}
}

func (f *FakeTemplateStore) Get(name string) (Template, error) {
	if f.GetErr != nil {
		return Template{}, f.GetErr
	}
	t, ok := f.Templates[name]
	if !ok {
		return Template{}, ErrTemplateNotFound
	}
	return t, nil
}

func (f *FakeTemplateStore) List() ([]TemplateEntry, error) {
	if f.ListErr != nil {
		return nil, f.ListErr
	}
	var entries []TemplateEntry
	for name := range f.Templates {
		entries = append(entries, TemplateEntry{Name: name, Source: SourceUser})
	}
	return entries, nil
}

func (f *FakeTemplateStore) Save(t Template) error {
	if f.SaveErr != nil {
		return f.SaveErr
	}
	f.Templates[t.Name] = t
	return nil
}

func (f *FakeTemplateStore) Delete(name string) error {
	if f.DeleteErr != nil {
		return f.DeleteErr
	}
	delete(f.Templates, name)
	return nil
}

func (f *FakeTemplateStore) Marshal(t Template) ([]byte, error) {
	// Minimal serialisation for tests: "name=<name>" is sufficient for round-trip.
	return []byte("name = \"" + t.Name + "\"\nsystem = \"" + t.System + "\"\nuser = \"" + t.User + "\"\n"), nil
}

func (f *FakeTemplateStore) Unmarshal(data []byte) (Template, error) {
	// Delegate to a real TOML unmarshal so tests aren't fragile.
	var t Template
	// Use a simple line parser to avoid importing toml in a fake.
	// This is sufficient for test assertions — production code uses FileStore.
	for _, line := range splitLines(string(data)) {
		if len(line) < 3 {
			continue
		}
		if line[:5] == "name " {
			t.Name = extractTOMLString(line)
		} else if len(line) > 7 && line[:7] == "system " {
			t.System = extractTOMLString(line)
		} else if len(line) > 5 && line[:5] == "user " {
			t.User = extractTOMLString(line)
		}
	}
	return t, nil
}

func splitLines(s string) []string {
	var lines []string
	cur := ""
	for _, c := range s {
		if c == '\n' {
			lines = append(lines, cur)
			cur = ""
		} else {
			cur += string(c)
		}
	}
	if cur != "" {
		lines = append(lines, cur)
	}
	return lines
}

func extractTOMLString(line string) string {
	start := -1
	end := -1
	for i, c := range line {
		if c == '"' {
			if start == -1 {
				start = i + 1
			} else {
				end = i
				break
			}
		}
	}
	if start != -1 && end != -1 && end > start {
		return line[start:end]
	}
	return ""
}
