package azure

import (
	"context"
	"fmt"

	"chat-backend/internal/chat"
)

type AzureChatProvider struct {
	client AzureQuestionAnsweringClient
}

func NewAzureChatProvider(endpoint, apiKey, projectName, deploymentName string) *AzureChatProvider {
	return &AzureChatProvider{
		client: NewClient(endpoint, apiKey, projectName, deploymentName),
	}
}

func (p *AzureChatProvider) Chat(ctx context.Context, req *chat.ChatRequest) (*chat.ChatResponse, error) {
	if len(req.Messages) == 0 {
		return nil, fmt.Errorf("no messages provided")
	}

	question := req.Messages[len(req.Messages)-1].Content

	queryResp, err := p.client.Query(ctx, question)
	if err != nil {
		return nil, err
	}

	if len(queryResp.Answers) == 0 {
		return &chat.ChatResponse{
			Content: "I don't have an answer for that question.",
		}, nil
	}

	return &chat.ChatResponse{
		Content: queryResp.Answers[0].Answer,
	}, nil
}