package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/labstack/echo/v4"

	"chat-backend/internal/app"
	"chat-backend/internal/chat"
	"chat-backend/internal/chat/azure"
	"chat-backend/internal/chat/mock"
	"chat-backend/internal/handlers"
	"chat-backend/internal/middleware"
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
	ctx := buildAppContext()

	e := echo.New()

	// Add middleware
	e.Use(middleware.Logger())
	e.Use(middleware.RateLimit())

	// Routes
	e.GET("/status", handlers.StatusHandler(ctx))
	e.POST("/api/faq", handlers.FAQHandler(ctx))

	slog.Info("Starting server on localhost:8090")
	log.Fatal(e.Start(":8090"))
}
