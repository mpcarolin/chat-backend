package chat

import (
	"context"
	"errors"
)

var (
	ErrProviderUnavailable = errors.New("chat provider unavailable")
	ErrInvalidRequest      = errors.New("invalid request")
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Messages  []Message `json:"messages"`
	Streaming bool      `json:"streaming,omitempty"`
}

type ChatResponse struct {
	Content string `json:"content"`
	Usage   *Usage `json:"usage,omitempty"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type StreamCallback func(chunk *ChatResponse) error

type ChatProvider interface {
	Chat(ctx context.Context, req *ChatRequest) (*ChatResponse, error)
	ChatStream(ctx context.Context, req *ChatRequest, callback StreamCallback) error
}