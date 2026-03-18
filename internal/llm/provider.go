package llm

import "context"

// Provider is a stateless LLM client. All providers use stdlib net/http only.
type Provider interface {
	Generate(ctx context.Context, system, user string) (string, error)
}
