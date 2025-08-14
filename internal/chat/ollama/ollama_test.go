package ollama

import (
	"context"
	"errors"
	"testing"

	"chat-backend/internal/chat"
)

type mockOllamaClient struct {
	chatFunc func(ctx context.Context, req *ChatRequest) (*ChatResponse, error)
}

func (m *mockOllamaClient) Chat(ctx context.Context, req *ChatRequest) (*ChatResponse, error) {
	if m.chatFunc != nil {
		return m.chatFunc(ctx, req)
	}
	return nil, errors.New("not implemented")
}

func TestOllamaChatProvider_Chat_Success(t *testing.T) {
	mockClient := &mockOllamaClient{
		chatFunc: func(ctx context.Context, req *ChatRequest) (*ChatResponse, error) {
			if len(req.Messages) != 1 {
				t.Errorf("expected 1 message, got %d", len(req.Messages))
			}
			if req.Messages[0].Role != "user" {
				t.Errorf("expected role 'user', got %s", req.Messages[0].Role)
			}
			if req.Messages[0].Content != "Hello" {
				t.Errorf("expected content 'Hello', got %s", req.Messages[0].Content)
			}
			if req.Stream != true {
				t.Error("expected Stream to be true")
			}

			return &ChatResponse{
				Message: OllamaMessage{
					Role:    "assistant",
					Content: "Hello there!",
				},
				Done: true,
			}, nil
		},
	}

	provider := &OllamaChatProvider{client: mockClient}

	req := &chat.ChatRequest{
		Messages: []chat.Message{
			{Role: "user", Content: "Hello"},
		},
		Streaming: true,
	}

	resp, err := provider.Chat(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.Content != "Hello there!" {
		t.Errorf("expected content 'Hello there!', got %s", resp.Content)
	}
}

func TestOllamaChatProvider_Chat_NoMessages(t *testing.T) {
	mockClient := &mockOllamaClient{}
	provider := &OllamaChatProvider{client: mockClient}

	req := &chat.ChatRequest{
		Messages: []chat.Message{},
	}

	_, err := provider.Chat(context.Background(), req)
	if err == nil {
		t.Fatal("expected error for empty messages")
	}

	if err.Error() != "no messages provided" {
		t.Errorf("expected 'no messages provided' error, got: %v", err)
	}
}

func TestOllamaChatProvider_Chat_ClientError(t *testing.T) {
	expectedError := errors.New("client error")
	mockClient := &mockOllamaClient{
		chatFunc: func(ctx context.Context, req *ChatRequest) (*ChatResponse, error) {
			return nil, expectedError
		},
	}

	provider := &OllamaChatProvider{client: mockClient}

	req := &chat.ChatRequest{
		Messages: []chat.Message{
			{Role: "user", Content: "Hello"},
		},
	}

	_, err := provider.Chat(context.Background(), req)
	if err != expectedError {
		t.Errorf("expected client error to be propagated, got: %v", err)
	}
}

func TestOllamaChatProvider_Chat_MultipleMessages(t *testing.T) {
	mockClient := &mockOllamaClient{
		chatFunc: func(ctx context.Context, req *ChatRequest) (*ChatResponse, error) {
			if len(req.Messages) != 3 {
				t.Errorf("expected 3 messages, got %d", len(req.Messages))
			}

			expectedMessages := []OllamaMessage{
				{Role: "system", Content: "You are a helpful assistant"},
				{Role: "user", Content: "Hello"},
				{Role: "assistant", Content: "Hi there!"},
			}

			for i, expected := range expectedMessages {
				if req.Messages[i].Role != expected.Role {
					t.Errorf("message %d: expected role '%s', got '%s'", i, expected.Role, req.Messages[i].Role)
				}
				if req.Messages[i].Content != expected.Content {
					t.Errorf("message %d: expected content '%s', got '%s'", i, expected.Content, req.Messages[i].Content)
				}
			}

			return &ChatResponse{
				Message: OllamaMessage{
					Role:    "assistant",
					Content: "Response to conversation",
				},
				Done: true,
			}, nil
		},
	}

	provider := &OllamaChatProvider{client: mockClient}

	req := &chat.ChatRequest{
		Messages: []chat.Message{
			{Role: "system", Content: "You are a helpful assistant"},
			{Role: "user", Content: "Hello"},
			{Role: "assistant", Content: "Hi there!"},
		},
		Streaming: false,
	}

	resp, err := provider.Chat(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.Content != "Response to conversation" {
		t.Errorf("expected content 'Response to conversation', got %s", resp.Content)
	}
}

func TestNewOllamaChatProvider(t *testing.T) {
	provider := NewOllamaChatProvider("http://test:8080", "test-model")
	
	if provider == nil {
		t.Fatal("expected provider to be created")
	}

	if provider.client == nil {
		t.Fatal("expected client to be initialized")
	}
}