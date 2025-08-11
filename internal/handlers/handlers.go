package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

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

func StatusHandler(appCtx *app.AppContext) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		encoded, err := json.Marshal(&Status{
			Date:   time.Now().UTC().String(),
			Status: "Running",
		})
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(encoded)
	}
}

func FAQHandler(appCtx *app.AppContext) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var faqReq FAQRequest
		if err := json.NewDecoder(req.Body).Decode(&faqReq); err != nil {
			slog.Error("Failed to decode FAQ request", "error", err)
			http.Error(w, "Invalid request format", http.StatusBadRequest)
			return
		}

		if faqReq.Question == "" {
			http.Error(w, "Question is required", http.StatusBadRequest)
			return
		}

		chatReq := &chat.ChatRequest{
			Messages: []chat.Message{
				{
					Role:    "user",
					Content: faqReq.Question,
				},
			},
		}

		chatResp, err := appCtx.ChatProvider.Chat(req.Context(), chatReq)
		if err != nil {
			slog.Error("Failed to get answer from chat provider", "error", err, "question", faqReq.Question)
			http.Error(w, "Failed to process question", http.StatusInternalServerError)
			return
		}

		faqResp := FAQResponse{
			Answer: chatResp.Content,
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(faqResp); err != nil {
			slog.Error("Failed to encode FAQ response", "error", err)
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	}
}