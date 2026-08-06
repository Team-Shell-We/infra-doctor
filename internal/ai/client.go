package ai

import (
	"context"
	"errors"
)

// CompletionRequest is a single-shot completion request. There is no
// multi-turn conversation concept in this product — every AI command sends
// one pre-built system prompt plus one pre-built user prompt and expects a
// single structured answer back.
type CompletionRequest struct {
	Model        string // empty means the provider's default model
	SystemPrompt string
	UserPrompt   string
	JSONMode     bool
	MaxTokens    int
}

type CompletionResponse struct {
	Content string
}

// Client is implemented once per AI provider (currently only OpenAI).
// VerifyCredentials and Complete are kept separate because login only needs
// to know whether a key works, while explain/recommend need an actual
// completion — the cheapest correct implementation of each differs.
type Client interface {
	Complete(ctx context.Context, req CompletionRequest) (*CompletionResponse, error)
	VerifyCredentials(ctx context.Context) error
}

var (
	ErrNotLoggedIn   = errors.New("not logged in")
	ErrInvalidAPIKey = errors.New("invalid OpenAI API Key")
)
