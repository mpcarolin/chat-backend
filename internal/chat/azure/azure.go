package azure

import (
	"context"
	"fmt"

	"chat-backend/internal/chat"
)

type AzureChatProvider struct {
	endpoint string
	apiKey   string
}

func NewAzureChatProvider(endpoint, apiKey string) *AzureChatProvider {
	return &AzureChatProvider{
		endpoint: endpoint,
		apiKey:   apiKey,
	}
}

func (p *AzureChatProvider) Chat(ctx context.Context, req *chat.ChatRequest) (*chat.ChatResponse, error) {
	return nil, fmt.Errorf("azure provider not implemented yet")
}