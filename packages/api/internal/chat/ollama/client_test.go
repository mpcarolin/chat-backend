package ollama

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name     string
		baseURL  string
		model    string
		expected struct {
			baseURL string
			model   string
		}
	}{
		{
			name:    "with custom values",
			baseURL: "http://custom:8080",
			model:   "llama2",
			expected: struct {
				baseURL string
				model   string
			}{
				baseURL: "http://custom:8080",
				model:   "llama2",
			},
		},
		{
			name:    "with empty values - should use defaults",
			baseURL: "",
			model:   "",
			expected: struct {
				baseURL string
				model   string
			}{
				baseURL: "http://localhost:11434",
				model:   "mistral",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient(tt.baseURL, tt.model).(*ollamaHttpClient)
			
			if client.baseURL != tt.expected.baseURL {
				t.Errorf("expected baseURL %s, got %s", tt.expected.baseURL, client.baseURL)
			}
			
			if client.model != tt.expected.model {
				t.Errorf("expected model %s, got %s", tt.expected.model, client.model)
			}
		})
	}
}

func TestOllamaHttpClient_Chat_NonStreaming(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST request, got %s", r.Method)
		}
		
		if r.URL.Path != "/api/chat" {
			t.Errorf("expected /api/chat path, got %s", r.URL.Path)
		}
		
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("expected application/json content type, got %s", r.Header.Get("Content-Type"))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message":{"role":"assistant","content":"Hello there!"},"done":true}`))
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-model")
	
	req := &ChatRequest{
		Messages: []OllamaMessage{
			{Role: "user", Content: "Hello"},
		},
		Stream: false,
	}

	resp, err := client.Chat(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.Message.Role != "assistant" {
		t.Errorf("expected role 'assistant', got %s", resp.Message.Role)
	}

	if resp.Message.Content != "Hello there!" {
		t.Errorf("expected content 'Hello there!', got %s", resp.Message.Content)
	}

	if !resp.Done {
		t.Error("expected Done to be true")
	}
}

func TestOllamaHttpClient_Chat_Streaming(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST request, got %s", r.Method)
		}
		
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		
		responses := []string{
			`{"message":{"role":"assistant","content":"Hello"},"done":false}`,
			`{"message":{"role":"assistant","content":" there"},"done":false}`,
			`{"message":{"role":"assistant","content":"!"},"done":true}`,
		}
		
		for _, resp := range responses {
			w.Write([]byte(resp + "\n"))
		}
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-model")
	
	req := &ChatRequest{
		Messages: []OllamaMessage{
			{Role: "user", Content: "Hello"},
		},
		Stream: true,
	}

	resp, err := client.Chat(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.Message.Role != "assistant" {
		t.Errorf("expected role 'assistant', got %s", resp.Message.Role)
	}

	if resp.Message.Content != "Hello there!" {
		t.Errorf("expected content 'Hello there!', got %s", resp.Message.Content)
	}

	if !resp.Done {
		t.Error("expected Done to be true")
	}
}

func TestOllamaHttpClient_Chat_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-model")
	
	req := &ChatRequest{
		Messages: []OllamaMessage{
			{Role: "user", Content: "Hello"},
		},
		Stream: false,
	}

	_, err := client.Chat(context.Background(), req)
	if err == nil {
		t.Fatal("expected error for server error response")
	}

	if !strings.Contains(err.Error(), "ollama API returned status 500") {
		t.Errorf("expected error message about status 500, got: %v", err)
	}
}

func TestOllamaHttpClient_Chat_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`invalid json`))
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-model")
	
	req := &ChatRequest{
		Messages: []OllamaMessage{
			{Role: "user", Content: "Hello"},
		},
		Stream: false,
	}

	_, err := client.Chat(context.Background(), req)
	if err == nil {
		t.Fatal("expected error for invalid JSON response")
	}

	if !strings.Contains(err.Error(), "failed to decode response") {
		t.Errorf("expected error message about decode failure, got: %v", err)
	}
}