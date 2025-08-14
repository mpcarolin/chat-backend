package handlers

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"chat-backend/internal/app"
	"chat-backend/internal/chat"
)

type Status struct {
	Date   string `json:"date"`
	Status string `json:"status"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Messages  []Message `json:"messages"`
	Streaming bool      `json:"streaming,omitempty"`
}

type ChatResponse struct {
	Response string `json:"response"`
}

func StatusHandler(appCtx *app.AppContext) echo.HandlerFunc {
	return func(c echo.Context) error {
		status := Status{
			Date:   time.Now().UTC().String(),
			Status: "Running",
		}
		return c.JSON(http.StatusOK, status)
	}
}

func ChatHandler(appCtx *app.AppContext) echo.HandlerFunc {
	return func(c echo.Context) error {
		var chatReq ChatRequest
		if err := c.Bind(&chatReq); err != nil {
			slog.Error("Failed to decode chat request", "error", err)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request format"})
		}

		// Validate messages array
		if len(chatReq.Messages) == 0 {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Messages array is required and cannot be empty"})
		}

		// Convert to internal message format
		var messages []chat.Message
		for _, msg := range chatReq.Messages {
			messages = append(messages, chat.Message{
				Role:    msg.Role,
				Content: msg.Content,
			})
		}

		chatRequest := &chat.ChatRequest{
			Messages:  messages,
			Streaming: chatReq.Streaming,
		}

		ctx := c.Request().Context()

		chatResp, err := appCtx.ChatProvider.Chat(ctx, chatRequest)
		if err != nil {
			slog.Error("Failed to get answer from chat provider", "error", err, "messages_count", len(chatReq.Messages))
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to process request"})
		}

		chatResponse := ChatResponse{
			Response: chatResp.Content,
		}

		return c.JSON(http.StatusOK, chatResponse)
	}
}
