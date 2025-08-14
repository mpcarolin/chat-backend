package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/labstack/echo/v4"

	"chat-backend/internal/app"
	"chat-backend/internal/chat"
)

type mockChatProvider struct {
	response *chat.ChatResponse
	err      error
}

func (m *mockChatProvider) Chat(ctx context.Context, req *chat.ChatRequest) (*chat.ChatResponse, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.response, nil
}

func init() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelWarn,
	})))
}

func TestChatHandler_Success(t *testing.T) {
	mockProvider := &mockChatProvider{
		response: &chat.ChatResponse{
			Content: "This is the answer to your question.",
		},
	}
	appCtx := app.NewAppContext(mockProvider)

	reqBody := ChatRequest{
		Messages: []Message{
			{Role: "user", Content: "What is the answer?"},
		},
	}
	jsonBody, _ := json.Marshal(reqBody)

	e := echo.New()
	req := httptest.NewRequest("POST", "/api/chat", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	c := e.NewContext(req, recorder)

	handler := ChatHandler(appCtx)
	err := handler(c)
	if err != nil {
		t.Fatalf("Handler returned error: %v", err)
	}

	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, recorder.Code)
	}

	var response ChatResponse
	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	expectedResponse := "This is the answer to your question."
	if response.Response != expectedResponse {
		t.Errorf("Expected response '%s', got '%s'", expectedResponse, response.Response)
	}
}

func TestChatHandler_EmptyMessages(t *testing.T) {
	mockProvider := &mockChatProvider{}
	appCtx := app.NewAppContext(mockProvider)

	reqBody := ChatRequest{
		Messages: []Message{},
	}
	jsonBody, _ := json.Marshal(reqBody)

	e := echo.New()
	req := httptest.NewRequest("POST", "/api/chat", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	c := e.NewContext(req, recorder)

	handler := ChatHandler(appCtx)
	err := handler(c)
	if err != nil {
		t.Fatalf("Handler returned error: %v", err)
	}

	if recorder.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, recorder.Code)
	}
}

func TestChatHandler_InvalidJSON(t *testing.T) {
	mockProvider := &mockChatProvider{}
	appCtx := app.NewAppContext(mockProvider)

	e := echo.New()
	req := httptest.NewRequest("POST", "/api/chat", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	c := e.NewContext(req, recorder)

	handler := ChatHandler(appCtx)
	err := handler(c)
	if err != nil {
		t.Fatalf("Handler returned error: %v", err)
	}

	if recorder.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, recorder.Code)
	}
}

func TestChatHandler_ChatProviderError(t *testing.T) {
	mockProvider := &mockChatProvider{
		err: chat.ErrProviderUnavailable,
	}
	appCtx := app.NewAppContext(mockProvider)

	reqBody := ChatRequest{
		Messages: []Message{
			{Role: "user", Content: "What is the answer?"},
		},
	}
	jsonBody, _ := json.Marshal(reqBody)

	e := echo.New()
	req := httptest.NewRequest("POST", "/api/chat", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	c := e.NewContext(req, recorder)

	handler := ChatHandler(appCtx)
	err := handler(c)
	if err != nil {
		t.Fatalf("Handler returned error: %v", err)
	}

	if recorder.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, recorder.Code)
	}
}

func TestStatusHandler(t *testing.T) {
	mockProvider := &mockChatProvider{}
	appCtx := app.NewAppContext(mockProvider)

	e := echo.New()
	req := httptest.NewRequest("GET", "/status", nil)
	recorder := httptest.NewRecorder()
	c := e.NewContext(req, recorder)

	handler := StatusHandler(appCtx)
	err := handler(c)
	if err != nil {
		t.Fatalf("Handler returned error: %v", err)
	}

	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, recorder.Code)
	}

	var response Status
	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response.Status != "Running" {
		t.Errorf("Expected status 'Running', got '%s'", response.Status)
	}

	if response.Date == "" {
		t.Error("Expected date to be set")
	}
}

func TestChatHandler_MultipleMessages(t *testing.T) {
	mockProvider := &mockChatProvider{
		response: &chat.ChatResponse{
			Content: "Based on our conversation, here's my response.",
		},
	}
	appCtx := app.NewAppContext(mockProvider)

	reqBody := ChatRequest{
		Messages: []Message{
			{Role: "user", Content: "Hello"},
			{Role: "assistant", Content: "Hi there!"},
			{Role: "user", Content: "How are you?"},
		},
		Streaming: true,
	}
	jsonBody, _ := json.Marshal(reqBody)

	e := echo.New()
	req := httptest.NewRequest("POST", "/api/chat", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	c := e.NewContext(req, recorder)

	handler := ChatHandler(appCtx)
	err := handler(c)
	if err != nil {
		t.Fatalf("Handler returned error: %v", err)
	}

	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, recorder.Code)
	}

	var response ChatResponse
	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	expectedResponse := "Based on our conversation, here's my response."
	if response.Response != expectedResponse {
		t.Errorf("Expected response '%s', got '%s'", expectedResponse, response.Response)
	}
}

func TestChatHandler_WithStreaming(t *testing.T) {
	mockProvider := &mockChatProvider{
		response: &chat.ChatResponse{
			Content: "Streaming response content.",
		},
	}
	appCtx := app.NewAppContext(mockProvider)

	reqBody := ChatRequest{
		Messages: []Message{
			{Role: "user", Content: "Tell me a story"},
		},
		Streaming: true,
	}
	jsonBody, _ := json.Marshal(reqBody)

	e := echo.New()
	req := httptest.NewRequest("POST", "/api/chat", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	c := e.NewContext(req, recorder)

	handler := ChatHandler(appCtx)
	err := handler(c)
	if err != nil {
		t.Fatalf("Handler returned error: %v", err)
	}

	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, recorder.Code)
	}

	var response ChatResponse
	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	expectedResponse := "Streaming response content."
	if response.Response != expectedResponse {
		t.Errorf("Expected response '%s', got '%s'", expectedResponse, response.Response)
	}
}