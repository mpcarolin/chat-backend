package ollama

import (
	"context"
	"fmt"

	"chat-backend/internal/chat"
)

type OllamaChatProvider struct {
	client OllamaClient
}

func NewOllamaChatProvider(baseURL, model string) *OllamaChatProvider {
	return &OllamaChatProvider{
		client: NewClient(baseURL, model),
	}
}

func (p *OllamaChatProvider) Chat(ctx context.Context, req *chat.ChatRequest) (*chat.ChatResponse, error) {
	if len(req.Messages) == 0 {
		return nil, fmt.Errorf("no messages provided")
	}

	ollamaMessages := make([]OllamaMessage, len(req.Messages))
	for i, msg := range req.Messages {
		ollamaMessages[i] = OllamaMessage{
			Role:    msg.Role,
			Content: msg.Content,
		}
	}

	ollamaReq := &ChatRequest{
		Messages: ollamaMessages,
		Stream:   req.Streaming,
	}

	ollamaResp, err := p.client.Chat(ctx, ollamaReq)
	if err != nil {
		return nil, err
	}

	return &chat.ChatResponse{
		Content: ollamaResp.Message.Content,
	}, nil
}

func (p *OllamaChatProvider) ChatStream(ctx context.Context, req *chat.ChatRequest, callback chat.StreamCallback) error {
	if len(req.Messages) == 0 {
		return fmt.Errorf("no messages provided")
	}

	ollamaMessages := make([]OllamaMessage, len(req.Messages))
	for i, msg := range req.Messages {
		ollamaMessages[i] = OllamaMessage{
			Role:    msg.Role,
			Content: msg.Content,
		}
	}

	ollamaReq := &ChatRequest{
		Messages: ollamaMessages,
		Stream:   true,
	}

	// Create a callback that converts ollama responses to chat responses
	ollamaCallback := func(ollamaResp *ChatResponse) error {
		chatResp := &chat.ChatResponse{
			Content: ollamaResp.Message.Content,
		}
		return callback(chatResp)
	}

	return p.client.ChatStream(ctx, ollamaReq, ollamaCallback)
}