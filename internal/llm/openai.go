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

const openAIEndpoint = "https://api.openai.com/v1/chat/completions"

type OpenAIProvider struct {
	model    string
	apiKey   string
	client   *http.Client
	endpoint string
}

func NewOpenAIProvider(model, apiKey string) *OpenAIProvider {
	return &OpenAIProvider{
		model:    model,
		apiKey:   apiKey,
		client:   &http.Client{Timeout: 30 * time.Second},
		endpoint: openAIEndpoint,
	}
}

// NewOpenAIProviderWithEndpoint is for testing only.
func NewOpenAIProviderWithEndpoint(model, apiKey, endpoint string) *OpenAIProvider {
	return &OpenAIProvider{
		model:    model,
		apiKey:   apiKey,
		client:   &http.Client{Timeout: 30 * time.Second},
		endpoint: endpoint,
	}
}

func (p *OpenAIProvider) Generate(ctx context.Context, system, user string) (string, error) {
	payload := map[string]interface{}{
		"model": p.model,
		"temperature": 0,
		"messages": []map[string]string{
			{"role": "system", "content": system},
			{"role": "user", "content": user},
		},
	}
	body, _ := json.Marshal(payload)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.endpoint, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+p.apiKey)

	resp, err := p.client.Do(req)
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()
	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", ErrProviderRequest{StatusCode: resp.StatusCode, Body: string(respBody)}
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", fmt.Errorf("parsing openai response: %w", err)
	}
	if len(result.Choices) == 0 {
		return "", fmt.Errorf("openai returned no choices")
	}
	return result.Choices[0].Message.Content, nil
}
