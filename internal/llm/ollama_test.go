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

func TestOllamaProvider_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Ollama MUST NOT have an Authorization header
		if r.Header.Get("Authorization") != "" {
			t.Errorf("unexpected Authorization header on Ollama request: %q", r.Header.Get("Authorization"))
		}
		resp := map[string]interface{}{
			"message": map[string]string{
				"role":    "assistant",
				"content": "feat: add ollama support",
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer srv.Close()

	p := llm.NewOllamaProviderWithEndpoint("llama3", srv.URL)
	msg, err := p.Generate(context.Background(), "system prompt", "user prompt")
	if err != nil {
		t.Fatal(err)
	}
	if msg != "feat: add ollama support" {
		t.Errorf("unexpected: %q", msg)
	}
}

func TestOllamaProvider_NonTwoXX(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"model not found"}`))
	}))
	defer srv.Close()

	p := llm.NewOllamaProviderWithEndpoint("llama3", srv.URL)
	_, err := p.Generate(context.Background(), "sys", "user")
	if err == nil {
		t.Fatal("expected error")
	}
	var provErr llm.ErrProviderRequest
	if !errors.As(err, &provErr) {
		t.Errorf("expected ErrProviderRequest, got %T: %v", err, err)
	}
	if provErr.StatusCode != 500 {
		t.Errorf("expected 500, got %d", provErr.StatusCode)
	}
}
