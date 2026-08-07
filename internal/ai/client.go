package ai

import (
	"context"
	"errors"
)

// CompletionRequest: 단발성 요청. 
// 모든 AI 명령어는 미리 만든 system/user 프롬프트 하나씩 보내고 구조화된 답 하나만 받음
type CompletionRequest struct {
	Model        string // 비어있으면 provider 기본 모델
	SystemPrompt string
	UserPrompt   string
	JSONMode     bool
	MaxTokens    int
}

type CompletionResponse struct {
	Content string
}

// Client는 AI provider마다 하나씩 구현(현재는 OpenAI뿐)
type Client interface {
	Complete(ctx context.Context, req CompletionRequest) (*CompletionResponse, error)
	VerifyCredentials(ctx context.Context) error
}

var (
	ErrNotLoggedIn   = errors.New("not logged in")
	ErrInvalidAPIKey = errors.New("invalid OpenAI API Key")
)
