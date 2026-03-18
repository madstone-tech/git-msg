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

const anthropicEndpoint = "https://api.anthropic.com/v1/messages"

type AnthropicProvider struct {
	model    string
	apiKey   string
	client   *http.Client
	endpoint string
}

func NewAnthropicProvider(model, apiKey string) *AnthropicProvider {
	return &AnthropicProvider{
		model:    model,
		apiKey:   apiKey,
		client:   &http.Client{Timeout: 30 * time.Second},
		endpoint: anthropicEndpoint,
	}
}

// NewAnthropicProviderWithEndpoint is for testing only.
func NewAnthropicProviderWithEndpoint(model, apiKey, endpoint string) *AnthropicProvider {
	return &AnthropicProvider{
		model:    model,
		apiKey:   apiKey,
		client:   &http.Client{Timeout: 30 * time.Second},
		endpoint: endpoint,
	}
}

func (p *AnthropicProvider) Generate(ctx context.Context, system, user string) (string, error) {
	payload := map[string]interface{}{
		"model":      p.model,
		"max_tokens": 1024,
		"system":     system,
		"messages": []map[string]string{
			{"role": "user", "content": user},
		},
	}
	body, _ := json.Marshal(payload)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.endpoint, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", p.apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

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
		Content []struct {
			Text string `json:"text"`
		} `json:"content"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", fmt.Errorf("parsing anthropic response: %w", err)
	}
	if len(result.Content) == 0 {
		return "", fmt.Errorf("anthropic returned no content")
	}
	return result.Content[0].Text, nil
}
