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
	// Serve the SPA
	e.Static("/", "./static/dist")

	// Serve the api endpoints
	e.GET("/status", handlers.StatusHandler(ctx))
	e.POST("/api/chat", handlers.ChatHandler(ctx))

	slog.Info("Starting server on localhost:8090")
	log.Fatal(e.Start(":8090"))
}
