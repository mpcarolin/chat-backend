package main

import (
	"embed"
	"log"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	emiddleware "github.com/labstack/echo/v4/middleware"

	"chat-backend/internal/app"
	"chat-backend/internal/handlers"
	"chat-backend/internal/middleware"
)

//go:embed static/dist
var webAssets embed.FS

func main() {
	ctx := app.BuildAppContext()

	e := echo.New()

	// Add middleware
	e.Use(middleware.Logger())
	e.Use(middleware.RateLimit())
	e.Use(emiddleware.StaticWithConfig(emiddleware.StaticConfig{
		HTML5:      true,
		Root:       "static/dist",
		Filesystem: http.FS(webAssets),
	}))

	// Routes
	// Serve the SPA
	e.Static("/", "./static/dist")

	// Serve the api endpoints
	e.GET("/status", handlers.StatusHandler(ctx))
	e.POST("/api/chat", handlers.ChatHandler(ctx))

	slog.Info("Starting server on localhost:8090")
	log.Fatal(e.Start(":8090"))
}
