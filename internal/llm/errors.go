package llm

import "fmt"

// ErrProviderRequest is returned when an LLM provider returns a non-2xx HTTP status.
// The API key is NEVER included in this error.
type ErrProviderRequest struct {
	StatusCode int
	Body       string
}

func (e ErrProviderRequest) Error() string {
	return fmt.Sprintf("provider request failed with status %d: %s", e.StatusCode, e.Body)
}
