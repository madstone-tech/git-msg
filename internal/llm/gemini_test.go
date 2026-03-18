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

func TestGeminiProvider_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Gemini uses API key in query string, not header
		if r.URL.Query().Get("key") == "" {
			t.Error("missing API key in query string")
		}
		resp := map[string]interface{}{
			"candidates": []map[string]interface{}{
				{
					"content": map[string]interface{}{
						"parts": []map[string]string{
							{"text": "feat: add gemini support"},
						},
					},
				},
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer srv.Close()

	// Gemini's baseURL override: append query params
	baseURL := srv.URL + "?key=test-key"
	p := llm.NewGeminiProviderWithEndpoint("gemini-1.5-flash", "test-key", baseURL)
	msg, err := p.Generate(context.Background(), "system prompt", "user prompt")
	if err != nil {
		t.Fatal(err)
	}
	if msg != "feat: add gemini support" {
		t.Errorf("unexpected: %q", msg)
	}
}

func TestGeminiProvider_NonTwoXX(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(`{"error":{"code":403,"message":"API_KEY_INVALID"}}`))
	}))
	defer srv.Close()

	p := llm.NewGeminiProviderWithEndpoint("gemini-1.5-flash", "bad-key", srv.URL)
	_, err := p.Generate(context.Background(), "sys", "user")
	if err == nil {
		t.Fatal("expected error")
	}
	var provErr llm.ErrProviderRequest
	if !errors.As(err, &provErr) {
		t.Errorf("expected ErrProviderRequest, got %T: %v", err, err)
	}
	if provErr.StatusCode != 403 {
		t.Errorf("expected 403, got %d", provErr.StatusCode)
	}
}
