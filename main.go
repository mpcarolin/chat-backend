package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"

	"chat-backend/internal/app"
	"chat-backend/internal/chat"
	"chat-backend/internal/chat/azure"
	"chat-backend/internal/chat/mock"
	"chat-backend/internal/handlers"
)

func buildAppContext() *app.AppContext {
	var chatProvider chat.ChatProvider

	if os.Getenv("MOCK_CHAT") == "true" {
		slog.Info("Using mock chat provider")
		chatProvider = mock.NewMockChatProvider()
	} else {
		endpoint := os.Getenv("AZURE_QNA_ENDPOINT")
		apiKey := os.Getenv("AZURE_QNA_API_KEY")
		projectName := os.Getenv("AZURE_QNA_PROJECT_NAME")
		deploymentName := os.Getenv("AZURE_QNA_DEPLOYMENT_NAME")

		if endpoint == "" || apiKey == "" || projectName == "" || deploymentName == "" {
			log.Fatal("Required Azure environment variables not set: AZURE_QNA_ENDPOINT, AZURE_QNA_API_KEY, AZURE_QNA_PROJECT_NAME, AZURE_QNA_DEPLOYMENT_NAME")
		}

		slog.Info("Using Azure chat provider")
		chatProvider = azure.NewAzureChatProvider(endpoint, apiKey, projectName, deploymentName)
	}

	return app.NewAppContext(chatProvider)
}

func main() {
	appCtx := buildAppContext()

	http.HandleFunc("/status", handlers.StatusHandler(appCtx))
	http.HandleFunc("/api/faq", handlers.FAQHandler(appCtx))

	slog.Info("Starting server on localhost:8090")
	err := http.ListenAndServe(":8090", nil)
	log.Fatal(err)
}
