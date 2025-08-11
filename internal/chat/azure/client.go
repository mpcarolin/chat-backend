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
)

type AzureQuestionAnsweringClient interface {
	Query(ctx context.Context, question string) (*QueryResponse, error)
}

type HTTPClient struct {
	endpoint       string
	apiKey         string
	projectName    string
	deploymentName string
	httpClient     *http.Client
}

type QueryRequest struct {
	Question                 string  `json:"question"`
	ConfidenceScoreThreshold float64 `json:"confidenceScoreThreshold"`
	Top                      int     `json:"top"`
}

type QueryResponse struct {
	Answers []QueryAnswer `json:"answers"`
}

type QueryAnswer struct {
	Answer          string            `json:"answer"`
	ConfidenceScore float64           `json:"confidenceScore"`
	Source          string            `json:"source"`
	Metadata        map[string]string `json:"metadata"`
}

func NewClient(endpoint, apiKey, projectName, deploymentName string) AzureQuestionAnsweringClient {
	return &HTTPClient{
		endpoint:       endpoint,
		apiKey:         apiKey,
		projectName:    projectName,
		deploymentName: deploymentName,
		httpClient:     &http.Client{},
	}
}

func (c *HTTPClient) Query(ctx context.Context, question string) (*QueryResponse, error) {
	queryReq := QueryRequest{
		Question:                 question,
		ConfidenceScoreThreshold: 0.2,
		Top:                      1,
	}

	jsonData, err := json.Marshal(queryReq)
	if err != nil {
		slog.Error("Failed to marshal Azure query request", "error", err)
		return nil, fmt.Errorf("failed to prepare request")
	}

	apiURL := fmt.Sprintf("%s/language/:query-knowledgebases?api-version=2021-10-01&projectName=%s&deploymentName=%s",
		c.endpoint, url.QueryEscape(c.projectName), url.QueryEscape(c.deploymentName))

	httpReq, err := http.NewRequestWithContext(ctx, "POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		slog.Error("Failed to create HTTP request", "error", err)
		return nil, fmt.Errorf("failed to create request")
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Ocp-Apim-Subscription-Key", c.apiKey)

	resp, err := c.httpClient.Do(httpReq)
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

	var queryResp QueryResponse
	if err := json.Unmarshal(body, &queryResp); err != nil {
		slog.Error("Failed to unmarshal Azure response", "error", err, "body", string(body))
		return nil, fmt.Errorf("failed to parse response")
	}

	return &queryResp, nil
}

