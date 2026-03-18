package llm

import "context"

// FakeProvider is a test double for Provider.
type FakeProvider struct {
	Response   string
	Err        error
	Called     bool
	LastSystem string
	LastUser   string
}

func (f *FakeProvider) Generate(ctx context.Context, system, user string) (string, error) {
	f.Called = true
	f.LastSystem = system
	f.LastUser = user
	return f.Response, f.Err
}
