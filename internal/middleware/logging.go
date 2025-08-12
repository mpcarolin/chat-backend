package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

func Logger(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		info := fmt.Sprintf("[%s] %s", time.Now().UTC().String(), req.Pattern)
		slog.Info(info)

		if next != nil {
			next(w, req)
		}
	}
}
