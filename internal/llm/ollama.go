package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type OllamaProvider struct {
	model    string
	host     string
	client   *http.Client
	endpoint string // override for testing
}

func NewOllamaProvider(model, host string) *OllamaProvider {
	return &OllamaProvider{
		model:  model,
		host:   host,
		client: &http.Client{Timeout: 30 * time.Second},
	}
}

// NewOllamaProviderWithEndpoint is for testing only.
func NewOllamaProviderWithEndpoint(model, endpoint string) *OllamaProvider {
	return &OllamaProvider{
		model:    model,
		client:   &http.Client{Timeout: 30 * time.Second},
		endpoint: endpoint,
	}
}

func (p *OllamaProvider) Generate(ctx context.Context, system, user string) (string, error) {
	endpoint := p.endpoint
	if endpoint == "" {
		endpoint = p.host + "/api/chat"
	}
	payload := map[string]interface{}{
		"model": p.model,
		"messages": []map[string]string{
			{"role": "system", "content": system},
			{"role": "user", "content": user},
		},
		"stream": false,
	}
	body, _ := json.Marshal(payload)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	// No Authorization header for Ollama

	resp, err := p.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", ErrProviderRequest{StatusCode: resp.StatusCode, Body: string(respBody)}
	}

	var result struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", fmt.Errorf("parsing ollama response: %w", err)
	}
	return result.Message.Content, nil
}
