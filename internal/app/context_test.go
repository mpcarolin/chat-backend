package app

import (
	"os"
	"testing"

	"chat-backend/internal/chat/azure"
	"chat-backend/internal/chat/mock"
	"chat-backend/internal/chat/ollama"
)

func TestBuildAppContext_Mock(t *testing.T) {
	// Clear any existing env vars
	os.Unsetenv("CHAT_PROVIDER")
	defer os.Unsetenv("CHAT_PROVIDER")

	ctx := BuildAppContext()
	
	if ctx == nil {
		t.Fatal("expected context to be created")
	}

	if ctx.ChatProvider == nil {
		t.Fatal("expected chat provider to be set")
	}

	// Check if it's a mock provider (we can't directly type assert due to interface)
	if _, ok := ctx.ChatProvider.(*mock.MockChatProvider); !ok {
		t.Error("expected mock chat provider for default case")
	}
}

func TestBuildAppContext_MockExplicit(t *testing.T) {
	os.Setenv("CHAT_PROVIDER", "mock")
	defer os.Unsetenv("CHAT_PROVIDER")

	ctx := BuildAppContext()
	
	if ctx == nil {
		t.Fatal("expected context to be created")
	}

	if _, ok := ctx.ChatProvider.(*mock.MockChatProvider); !ok {
		t.Error("expected mock chat provider when CHAT_PROVIDER=mock")
	}
}

func TestBuildAppContext_Ollama(t *testing.T) {
	os.Setenv("CHAT_PROVIDER", "ollama")
	defer os.Unsetenv("CHAT_PROVIDER")

	ctx := BuildAppContext()
	
	if ctx == nil {
		t.Fatal("expected context to be created")
	}

	if _, ok := ctx.ChatProvider.(*ollama.OllamaChatProvider); !ok {
		t.Error("expected ollama chat provider when CHAT_PROVIDER=ollama")
	}
}

func TestBuildAppContext_OllamaWithEnvs(t *testing.T) {
	os.Setenv("CHAT_PROVIDER", "ollama")
	os.Setenv("OLLAMA_BASE_URL", "http://custom:8080")
	os.Setenv("OLLAMA_MODEL", "custom-model")
	defer func() {
		os.Unsetenv("CHAT_PROVIDER")
		os.Unsetenv("OLLAMA_BASE_URL")
		os.Unsetenv("OLLAMA_MODEL")
	}()

	ctx := BuildAppContext()
	
	if ctx == nil {
		t.Fatal("expected context to be created")
	}

	if _, ok := ctx.ChatProvider.(*ollama.OllamaChatProvider); !ok {
		t.Error("expected ollama chat provider when CHAT_PROVIDER=ollama")
	}
}

func TestBuildAppContext_AzureQA(t *testing.T) {
	os.Setenv("CHAT_PROVIDER", "azure-qa")
	os.Setenv("AZURE_QNA_ENDPOINT", "https://test.cognitiveservices.azure.com")
	os.Setenv("AZURE_QNA_API_KEY", "test-key")
	os.Setenv("AZURE_QNA_PROJECT_NAME", "test-project")
	os.Setenv("AZURE_QNA_DEPLOYMENT_NAME", "test-deployment")
	defer func() {
		os.Unsetenv("CHAT_PROVIDER")
		os.Unsetenv("AZURE_QNA_ENDPOINT")
		os.Unsetenv("AZURE_QNA_API_KEY")
		os.Unsetenv("AZURE_QNA_PROJECT_NAME")
		os.Unsetenv("AZURE_QNA_DEPLOYMENT_NAME")
	}()

	ctx := BuildAppContext()
	
	if ctx == nil {
		t.Fatal("expected context to be created")
	}

	if _, ok := ctx.ChatProvider.(*azure.AzureChatProvider); !ok {
		t.Error("expected azure chat provider when CHAT_PROVIDER=azure-qa")
	}
}

// Note: TestBuildAppContext_UnknownProvider and TestBuildAppContext_AzureMissingEnvs
// are commented out because they call log.Fatal/log.Fatalf which terminates the process.
// In a real test environment, you would need to refactor BuildAppContext to return
// errors instead of calling log.Fatal, or use a testing framework that can handle
// process termination.

/*
func TestBuildAppContext_UnknownProvider(t *testing.T) {
	os.Setenv("CHAT_PROVIDER", "unknown")
	defer os.Unsetenv("CHAT_PROVIDER")
	
	// This test would terminate the process due to log.Fatalf
	// In production code, consider returning an error instead
	// BuildAppContext() // This would call log.Fatalf and exit
}

func TestBuildAppContext_AzureMissingEnvs(t *testing.T) {
	os.Setenv("CHAT_PROVIDER", "azure-qa")
	defer os.Unsetenv("CHAT_PROVIDER")
	
	// This test would terminate the process due to log.Fatal
	// In production code, consider returning an error instead
	// BuildAppContext() // This would call log.Fatal and exit
}
*/