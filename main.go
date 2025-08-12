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
	mw "chat-backend/internal/middleware"
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
			log.Fatal("All required Azure envs must be set: AZURE_QNA_ENDPOINT, AZURE_QNA_API_KEY, AZURE_QNA_PROJECT_NAME, AZURE_QNA_DEPLOYMENT_NAME")
		}

		slog.Info("Using Azure chat provider")
		chatProvider = azure.NewAzureChatProvider(endpoint, apiKey, projectName, deploymentName)
	}

	return app.NewAppContext(chatProvider)
}

func main() {
	// TODO: change this to use echo/labstack library and rate limiting middleware
	ctx := buildAppContext()

	mux := http.NewServeMux()

	mux.HandleFunc("GET /status", mw.Logger(handlers.StatusHandler(ctx)))
	mux.HandleFunc("POST /api/faq", mw.Logger(handlers.FAQHandler(ctx)))

	slog.Info("Starting server on localhost:8090")
	err := http.ListenAndServe(":8090", mux)
	log.Fatal(err)
}
