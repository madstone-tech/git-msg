package llm_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/madstone0-0/git-msg/internal/llm"
)

func TestOpenAIProvider_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") == "" {
			t.Error("missing Authorization header")
		}
		resp := map[string]interface{}{
			"choices": []map[string]interface{}{
				{"message": map[string]string{"content": "feat: add thing"}},
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer srv.Close()

	p := llm.NewOpenAIProviderWithEndpoint("gpt-4o-mini", "test-key", srv.URL)
	msg, err := p.Generate(context.Background(), "system prompt", "user prompt")
	if err != nil {
		t.Fatal(err)
	}
	if msg != "feat: add thing" {
		t.Errorf("unexpected: %q", msg)
	}
}

func TestOpenAIProvider_NonTwoXX(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error":{"message":"invalid_api_key"}}`))
	}))
	defer srv.Close()

	p := llm.NewOpenAIProviderWithEndpoint("gpt-4o-mini", "bad-key", srv.URL)
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
