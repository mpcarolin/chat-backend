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

	// Find the last user message to use as the question
	var question string
	for i := len(req.Messages) - 1; i >= 0; i-- {
		if req.Messages[i].Role == "user" {
			question = req.Messages[i].Content
			break
		}
	}

	if question == "" {
		return nil, fmt.Errorf("no user message found in conversation")
	}

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

func (p *AzureChatProvider) ChatStream(ctx context.Context, req *chat.ChatRequest, callback chat.StreamCallback) error {
	return fmt.Errorf("streaming not supported by azure provider")
}

