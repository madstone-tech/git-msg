package llm_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/madstone-tech/git-msg/internal/llm"
)

func TestAnthropicProvider_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Assert required headers
		if r.Header.Get("x-api-key") == "" {
			t.Error("missing x-api-key header")
		}
		if r.Header.Get("anthropic-version") != "2023-06-01" {
			t.Errorf("unexpected anthropic-version: %q", r.Header.Get("anthropic-version"))
		}
		resp := map[string]interface{}{
			"content": []map[string]string{
				{"type": "text", "text": "feat: add anthropic support"},
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer srv.Close()

	p := llm.NewAnthropicProviderWithEndpoint("claude-haiku-4-5", "test-key", srv.URL)
	msg, err := p.Generate(context.Background(), "system prompt", "user prompt")
	if err != nil {
		t.Fatal(err)
	}
	if msg != "feat: add anthropic support" {
		t.Errorf("unexpected message: %q", msg)
	}
}

func TestAnthropicProvider_NonTwoXX(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"type":"error","error":{"type":"authentication_error"}}`))
	}))
	defer srv.Close()

	p := llm.NewAnthropicProviderWithEndpoint("claude-haiku-4-5", "bad-key", srv.URL)
	_, err := p.Generate(context.Background(), "sys", "user")
	if err == nil {
		t.Fatal("expected error")
	}
	var provErr llm.ErrProviderRequest
	if !errors.As(err, &provErr) {
		t.Errorf("expected ErrProviderRequest, got %T: %v", err, err)
	}
	if provErr.StatusCode != 401 {
		t.Errorf("expected 401, got %d", provErr.StatusCode)
	}
}
