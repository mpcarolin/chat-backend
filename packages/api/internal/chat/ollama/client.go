package ollama

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type OllamaClient interface {
	Chat(ctx context.Context, req *ChatRequest) (*ChatResponse, error)
}

type ollamaHttpClient struct {
	baseURL    string
	model      string
	httpClient *http.Client
}

type OllamaMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Model    string          `json:"model"`
	Messages []OllamaMessage `json:"messages"`
	Stream   bool            `json:"stream"`
}

type ChatResponse struct {
	Message OllamaMessage `json:"message"`
	Done    bool          `json:"done"`
}

func NewClient(baseURL, model string) OllamaClient {
	if baseURL == "" {
		baseURL = "http://localhost:11434"
	}
	if model == "" {
		model = "mistral"
	}

	return &ollamaHttpClient{
		baseURL:    baseURL,
		model:      model,
		httpClient: &http.Client{},
	}
}

func (c *ollamaHttpClient) Chat(ctx context.Context, req *ChatRequest) (*ChatResponse, error) {
	req.Model = c.model

	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/api/chat", c.baseURL)
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ollama API returned status %d", resp.StatusCode)
	}

	if req.Stream {
		return c.handleStreamingResponse(resp.Body)
	}

	return c.handleNonStreamingResponse(resp.Body)
}

// TODO: this isn't right, not actually streaming. Let's fix.
func (c *ollamaHttpClient) handleStreamingResponse(body io.Reader) (*ChatResponse, error) {
	var fullContent strings.Builder
	decoder := json.NewDecoder(body)

	for {
		var ollamaResp ChatResponse
		if err := decoder.Decode(&ollamaResp); err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("failed to decode streaming response: %w", err)
		}

		fullContent.WriteString(ollamaResp.Message.Content)

		if ollamaResp.Done {
			break
		}
	}

	return &ChatResponse{
		Message: OllamaMessage{
			Role:    "assistant",
			Content: fullContent.String(),
		},
		Done: true,
	}, nil
}

func (c *ollamaHttpClient) handleNonStreamingResponse(body io.Reader) (*ChatResponse, error) {
	var ollamaResp ChatResponse
	if err := json.NewDecoder(body).Decode(&ollamaResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &ollamaResp, nil
}

