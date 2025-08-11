package main

import (
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
	"time"

	"chat-backend/internal/app"
	"chat-backend/internal/chat/azure"
)

type Status struct {
	Date   string `json:"date"`
	Status string `json:"status"`
}

func statusHandler(appCtx *app.AppContext) http.HandlerFunc {
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

func main() {
	azureProvider := azure.NewAzureChatProvider("", "")
	appCtx := app.NewAppContext(azureProvider)

	http.HandleFunc("/status", statusHandler(appCtx))

	slog.Info("Starting server on localhost:8090")
	err := http.ListenAndServe(":8090", nil)
	log.Fatal(err)
}
