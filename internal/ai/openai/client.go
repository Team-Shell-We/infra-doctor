// Package openai implements ai.Client against OpenAI's Chat Completions
// API. It is a direct net/http integration rather than the official SDK —
// this is one stable, non-streaming endpoint, not worth the dependency
// weight of a full client library.
package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Team-Shell-We/infra-doctor/internal/ai"
)

const (
	defaultBaseURL = "https://api.openai.com/v1"
	defaultModel   = "gpt-4o-mini"
)

type Client struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

func New(apiKey string) *Client {
	return &Client{
		apiKey:     apiKey,
		baseURL:    defaultBaseURL,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

var _ ai.Client = (*Client)(nil)

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type responseFormat struct {
	Type string `json:"type"`
}

type chatRequest struct {
	Model          string          `json:"model"`
	Messages       []chatMessage   `json:"messages"`
	ResponseFormat *responseFormat `json:"response_format,omitempty"`
	MaxTokens      int             `json:"max_tokens,omitempty"`
}

type chatChoice struct {
	Message chatMessage `json:"message"`
}

type chatResponse struct {
	Choices []chatChoice `json:"choices"`
}

func (c *Client) Complete(ctx context.Context, req ai.CompletionRequest) (*ai.CompletionResponse, error) {

	model := req.Model
	if model == "" {
		model = defaultModel
	}

	body := chatRequest{
		Model: model,
		Messages: []chatMessage{
			{Role: "system", Content: req.SystemPrompt},
			{Role: "user", Content: req.UserPrompt},
		},
		MaxTokens: req.MaxTokens,
	}

	if req.JSONMode {
		body.ResponseFormat = &responseFormat{Type: "json_object"}
	}

	var result chatResponse

	if err := c.post(ctx, "/chat/completions", body, &result); err != nil {
		return nil, err
	}

	if len(result.Choices) == 0 {
		return nil, fmt.Errorf("openai: empty response")
	}

	return &ai.CompletionResponse{Content: result.Choices[0].Message.Content}, nil
}

// VerifyCredentials makes the cheapest possible authenticated call
// (list models) purely to confirm the key is accepted, without spending
// completion tokens.
func (c *Client) VerifyCredentials(ctx context.Context) error {

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/models", nil)
	if err != nil {
		return err
	}

	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return ai.ErrInvalidAPIKey
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("openai: unexpected status verifying credentials: %s", resp.Status)
	}

	return nil
}

func (c *Client) post(ctx context.Context, path string, body, out any) error {

	payload, err := json.Marshal(body)
	if err != nil {
		return err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+path, bytes.NewReader(payload))
	if err != nil {
		return err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode == http.StatusUnauthorized {
		return ai.ErrInvalidAPIKey
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("openai: request failed with status %s: %s", resp.Status, string(data))
	}

	return json.Unmarshal(data, out)
}
