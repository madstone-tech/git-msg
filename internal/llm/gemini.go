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

type GeminiProvider struct {
	model   string
	apiKey  string
	client  *http.Client
	baseURL string // override for testing
}

func NewGeminiProvider(model, apiKey string) *GeminiProvider {
	return &GeminiProvider{
		model:  model,
		apiKey: apiKey,
		client: &http.Client{Timeout: 30 * time.Second},
	}
}

// NewGeminiProviderWithEndpoint is for testing only.
func NewGeminiProviderWithEndpoint(model, apiKey, baseURL string) *GeminiProvider {
	return &GeminiProvider{
		model:   model,
		apiKey:  apiKey,
		client:  &http.Client{Timeout: 30 * time.Second},
		baseURL: baseURL,
	}
}

func (p *GeminiProvider) Generate(ctx context.Context, system, user string) (string, error) {
	var endpoint string
	if p.baseURL != "" {
		endpoint = p.baseURL
	} else {
		endpoint = fmt.Sprintf(
			"https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s",
			p.model, p.apiKey,
		)
	}
	payload := map[string]interface{}{
		"system_instruction": map[string]interface{}{
			"parts": []map[string]string{{"text": system}},
		},
		"contents": []map[string]interface{}{
			{"parts": []map[string]string{{"text": user}}},
		},
	}
	body, _ := json.Marshal(payload)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

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
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", fmt.Errorf("parsing gemini response: %w", err)
	}
	if len(result.Candidates) == 0 || len(result.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("gemini returned no candidates")
	}
	return result.Candidates[0].Content.Parts[0].Text, nil
}
