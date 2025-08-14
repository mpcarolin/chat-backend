package app

import (
	"log"
	"log/slog"
	"os"

	"chat-backend/internal/chat"
	"chat-backend/internal/chat/azure"
	"chat-backend/internal/chat/mock"
	"chat-backend/internal/chat/ollama"
)

type AppContext struct {
	ChatProvider chat.ChatProvider
}

func NewAppContext(chatProvider chat.ChatProvider) *AppContext {
	return &AppContext{
		ChatProvider: chatProvider,
	}
}

func providerName() string {
	provider := os.Getenv("CHAT_PROVIDER")
	if provider == "" {
		provider = "mock"
	}
	return provider
}

// Builds context object containing app dependencies used by handlers
// In particular, contains chat provider, depending on whichever provider
// you chose to set with CHAT_PROVIDER env
func BuildAppContext() *AppContext {
	var chatProvider chat.ChatProvider

	provider := providerName()

	switch provider {
	case "mock":
		slog.Info("Using mock chat provider")
		chatProvider = mock.NewMockChatProvider()

	case "azure-qa":
		endpoint := os.Getenv("AZURE_QNA_ENDPOINT")
		apiKey := os.Getenv("AZURE_QNA_API_KEY")
		projectName := os.Getenv("AZURE_QNA_PROJECT_NAME")
		deploymentName := os.Getenv("AZURE_QNA_DEPLOYMENT_NAME")

		if endpoint == "" || apiKey == "" || projectName == "" || deploymentName == "" {
			log.Fatal("All required Azure envs must be set: AZURE_QNA_ENDPOINT, AZURE_QNA_API_KEY, AZURE_QNA_PROJECT_NAME, AZURE_QNA_DEPLOYMENT_NAME")
		}

		slog.Info("Using Azure chat provider")
		chatProvider = azure.NewAzureChatProvider(endpoint, apiKey, projectName, deploymentName)

	case "ollama":
		baseURL := os.Getenv("OLLAMA_BASE_URL")
		model := os.Getenv("OLLAMA_MODEL")

		slog.Info("Using Ollama chat provider", "baseURL", baseURL, "model", model)
		chatProvider = ollama.NewOllamaChatProvider(baseURL, model)

	default:
		log.Fatalf("Unknown CHAT_PROVIDER: %s. Supported values: mock, azure-qa, ollama", provider)
	}

	return NewAppContext(chatProvider)
}
