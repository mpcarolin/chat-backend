package azure

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"

	"chat-backend/internal/chat"
)

type AzureChatProvider struct {
	endpoint       string
	apiKey         string
	projectName    string
	deploymentName string
	httpClient     *http.Client
}

type azureRequest struct {
	Question                 string  `json:"question"`
	ConfidenceScoreThreshold float64 `json:"confidenceScoreThreshold"`
	Top                      int     `json:"top"`
}

type azureResponse struct {
	Answers []azureAnswer `json:"answers"`
}

type azureAnswer struct {
	Answer          string            `json:"answer"`
	ConfidenceScore float64           `json:"confidenceScore"`
	Source          string            `json:"source"`
	Metadata        map[string]string `json:"metadata"`
}

func NewAzureChatProvider(endpoint, apiKey, projectName, deploymentName string) *AzureChatProvider {
	return &AzureChatProvider{
		endpoint:       endpoint,
		apiKey:         apiKey,
		projectName:    projectName,
		deploymentName: deploymentName,
		httpClient:     &http.Client{},
	}
}

func (p *AzureChatProvider) Chat(ctx context.Context, req *chat.ChatRequest) (*chat.ChatResponse, error) {
	if len(req.Messages) == 0 {
		return nil, fmt.Errorf("no messages provided")
	}

	question := req.Messages[len(req.Messages)-1].Content

	azureReq := azureRequest{
		Question:                 question,
		ConfidenceScoreThreshold: 0.2,
		Top:                      1,
	}

	jsonData, err := json.Marshal(azureReq)
	if err != nil {
		slog.Error("Failed to marshal Azure request", "error", err)
		return nil, fmt.Errorf("failed to prepare request")
	}

	apiURL := fmt.Sprintf("%s/language/:query-knowledgebases?api-version=2023-04-01&projectName=%s&deploymentName=%s",
		p.endpoint, url.QueryEscape(p.projectName), url.QueryEscape(p.deploymentName))

	httpReq, err := http.NewRequestWithContext(ctx, "POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		slog.Error("Failed to create HTTP request", "error", err)
		return nil, fmt.Errorf("failed to create request")
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Ocp-Apim-Subscription-Key", p.apiKey)

	resp, err := p.httpClient.Do(httpReq)
	if err != nil {
		slog.Error("Failed to make request to Azure", "error", err, "url", apiURL)
		return nil, fmt.Errorf("failed to connect to Azure service")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("Failed to read Azure response body", "error", err)
		return nil, fmt.Errorf("failed to read response")
	}

	if resp.StatusCode != http.StatusOK {
		slog.Error("Azure API returned error", "status", resp.StatusCode, "body", string(body))
		return nil, fmt.Errorf("Azure service error")
	}

	var azureResp azureResponse
	if err := json.Unmarshal(body, &azureResp); err != nil {
		slog.Error("Failed to unmarshal Azure response", "error", err, "body", string(body))
		return nil, fmt.Errorf("failed to parse response")
	}

	if len(azureResp.Answers) == 0 {
		return &chat.ChatResponse{
			Content: "I don't have an answer for that question.",
		}, nil
	}

	return &chat.ChatResponse{
		Content: azureResp.Answers[0].Answer,
	}, nil
}