package handlers

import (
	"fmt"
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

		if chatReq.Streaming {
			// Set up Server-Sent Events headers
			c.Response().Header().Set("Content-Type", "text/event-stream")
			c.Response().Header().Set("Cache-Control", "no-cache")
			c.Response().Header().Set("Connection", "keep-alive")
			c.Response().Header().Set("Access-Control-Allow-Origin", "*")

			// Flush headers immediately
			c.Response().WriteHeader(http.StatusOK)
			c.Response().Flush()

			// Stream callback function
			streamCallback := func(chunk *chat.ChatResponse) error {
				data := fmt.Sprintf("data: {\"response\": \"%s\", \"done\": false}\n\n", chunk.Content)
				if _, err := c.Response().Write([]byte(data)); err != nil {
					return err
				}
				c.Response().Flush()
				return nil
			}

			err := appCtx.ChatProvider.ChatStream(ctx, chatRequest, streamCallback)
			if err != nil {
				slog.Error("Failed to stream chat response", "error", err, "messages_count", len(chatReq.Messages))
				errorData := fmt.Sprintf("data: {\"error\": \"Failed to process request\"}\n\n")
				c.Response().Write([]byte(errorData))
				c.Response().Flush()
				return nil
			}

			// Send final done message
			doneData := fmt.Sprintf("data: {\"done\": true}\n\n")
			c.Response().Write([]byte(doneData))
			c.Response().Flush()
			return nil
		}

		// Non-streaming response (existing behavior)
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
