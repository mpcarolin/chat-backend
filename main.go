package main

import (
	"log"
	"log/slog"

	"github.com/labstack/echo/v4"

	"chat-backend/internal/app"
	"chat-backend/internal/handlers"
	"chat-backend/internal/middleware"
)

func main() {
	ctx := app.BuildAppContext()

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
