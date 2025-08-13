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

type FAQRequest struct {
	Question string `json:"question"`
}

type FAQResponse struct {
	Answer string `json:"answer"`
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

func FAQHandler(appCtx *app.AppContext) echo.HandlerFunc {
	return func(c echo.Context) error {
		var faqReq FAQRequest
		if err := c.Bind(&faqReq); err != nil {
			slog.Error("Failed to decode FAQ request", "error", err)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request format"})
		}

		if faqReq.Question == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Question is required"})
		}

		chatReq := &chat.ChatRequest{
			Messages: []chat.Message{
				{
					Role:    "user",
					Content: faqReq.Question,
				},
			},
		}

		chatResp, err := appCtx.ChatProvider.Chat(c.Request().Context(), chatReq)
		if err != nil {
			slog.Error("Failed to get answer from chat provider", "error", err, "question", faqReq.Question)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to process question"})
		}

		faqResp := FAQResponse{
			Answer: chatResp.Content,
		}

		return c.JSON(http.StatusOK, faqResp)
	}
}

