package mock

import (
	"context"
	"fmt"

	"chat-backend/internal/chat"
)

type MockChatProvider struct{}

func NewMockChatProvider() *MockChatProvider {
	return &MockChatProvider{}
}

func (m *MockChatProvider) Chat(ctx context.Context, req *chat.ChatRequest) (*chat.ChatResponse, error) {
	if len(req.Messages) == 0 {
		return nil, fmt.Errorf("no messages provided")
	}

	question := req.Messages[len(req.Messages)-1].Content

	var answer string
	switch {
	case question == "":
		answer = "Please provide a question."
	case len(question) < 3:
		answer = "Your question seems too short. Could you provide more details?"
	default:
		answer = fmt.Sprintf("This is a mock response to your question: '%s'. In a real implementation, this would come from a chat service.", question)
	}

	return &chat.ChatResponse{
		Content: answer,
	}, nil
}

