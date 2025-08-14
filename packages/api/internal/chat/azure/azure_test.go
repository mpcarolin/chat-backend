package azure

import (
	"context"
	"testing"

	"chat-backend/internal/chat"
)

type mockAzureClient struct {
	response *QueryResponse
	err      error
	lastQuery string
}

func (m *mockAzureClient) Query(ctx context.Context, query string) (*QueryResponse, error) {
	m.lastQuery = query
	if m.err != nil {
		return nil, m.err
	}
	return m.response, nil
}

func TestAzureChatProvider_ExtractsLastUserMessage(t *testing.T) {
	mockClient := &mockAzureClient{
		response: &QueryResponse{
			Answers: []QueryAnswer{
				{
					Answer: "This is the answer from Azure QnA",
				},
			},
		},
	}

	provider := &AzureChatProvider{
		client: mockClient,
	}

	req := &chat.ChatRequest{
		Messages: []chat.Message{
			{Role: "user", Content: "First question"},
			{Role: "assistant", Content: "First response"},
			{Role: "user", Content: "This should be extracted"},
			{Role: "assistant", Content: "Another response"},
		},
	}

	_, err := provider.Chat(context.Background(), req)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	expectedQuery := "This should be extracted"
	if mockClient.lastQuery != expectedQuery {
		t.Errorf("Expected query '%s', got '%s'", expectedQuery, mockClient.lastQuery)
	}
}

func TestAzureChatProvider_HandlesOnlyAssistantMessages(t *testing.T) {
	mockClient := &mockAzureClient{}

	provider := &AzureChatProvider{
		client: mockClient,
	}

	req := &chat.ChatRequest{
		Messages: []chat.Message{
			{Role: "system", Content: "System message"},
			{Role: "assistant", Content: "Assistant response"},
		},
	}

	_, err := provider.Chat(context.Background(), req)
	if err == nil || err.Error() != "no user message found in conversation" {
		t.Errorf("Expected 'no user message found in conversation' error, got: %v", err)
	}
}

func TestAzureChatProvider_FindsLastUserMessage(t *testing.T) {
	mockClient := &mockAzureClient{
		response: &QueryResponse{
			Answers: []QueryAnswer{
				{
					Answer: "Response to last question",
				},
			},
		},
	}

	provider := &AzureChatProvider{
		client: mockClient,
	}

	req := &chat.ChatRequest{
		Messages: []chat.Message{
			{Role: "user", Content: "First user message"},
			{Role: "assistant", Content: "Response"},
			{Role: "system", Content: "System instruction"},
			{Role: "user", Content: "Last user message"},
			{Role: "assistant", Content: "Another response"},
			{Role: "system", Content: "Final system message"},
		},
	}

	_, err := provider.Chat(context.Background(), req)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	expectedQuery := "Last user message"
	if mockClient.lastQuery != expectedQuery {
		t.Errorf("Expected query '%s', got '%s'", expectedQuery, mockClient.lastQuery)
	}
}