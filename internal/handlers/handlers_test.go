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

func TestFAQHandler_Success(t *testing.T) {
	mockProvider := &mockChatProvider{
		response: &chat.ChatResponse{
			Content: "This is the answer to your question.",
		},
	}
	appCtx := app.NewAppContext(mockProvider)

	reqBody := FAQRequest{Question: "What is the answer?"}
	jsonBody, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/api/faq", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	handler := FAQHandler(appCtx)
	handler(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, recorder.Code)
	}

	var response FAQResponse
	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	expectedAnswer := "This is the answer to your question."
	if response.Answer != expectedAnswer {
		t.Errorf("Expected answer '%s', got '%s'", expectedAnswer, response.Answer)
	}
}

func TestFAQHandler_EmptyQuestion(t *testing.T) {
	mockProvider := &mockChatProvider{}
	appCtx := app.NewAppContext(mockProvider)

	reqBody := FAQRequest{Question: ""}
	jsonBody, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/api/faq", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	handler := FAQHandler(appCtx)
	handler(recorder, req)

	if recorder.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, recorder.Code)
	}
}

func TestFAQHandler_InvalidJSON(t *testing.T) {
	mockProvider := &mockChatProvider{}
	appCtx := app.NewAppContext(mockProvider)

	req := httptest.NewRequest("POST", "/api/faq", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	handler := FAQHandler(appCtx)
	handler(recorder, req)

	if recorder.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, recorder.Code)
	}
}

func TestFAQHandler_ChatProviderError(t *testing.T) {
	mockProvider := &mockChatProvider{
		err: chat.ErrProviderUnavailable,
	}
	appCtx := app.NewAppContext(mockProvider)

	reqBody := FAQRequest{Question: "What is the answer?"}
	jsonBody, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/api/faq", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	handler := FAQHandler(appCtx)
	handler(recorder, req)

	if recorder.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, recorder.Code)
	}
}

func TestStatusHandler(t *testing.T) {
	mockProvider := &mockChatProvider{}
	appCtx := app.NewAppContext(mockProvider)

	req := httptest.NewRequest("GET", "/status", nil)
	recorder := httptest.NewRecorder()

	handler := StatusHandler(appCtx)
	handler(recorder, req)

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

