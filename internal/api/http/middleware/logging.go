package middleware

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/Dercraker/SearchEngine/internal/shared"
	"github.com/Dercraker/SearchEngine/internal/shared/requestId"
	"github.com/google/uuid"
)

func Logging(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			startTime := time.Now()

			requestID := r.Header.Get("X-Request-ID")
			if requestID == "" {
				requestID = uuid.NewString()
			}

			r = r.WithContext(requestId.WithRunId(r.Context(), requestID))

			rw := newResponseWriter(w)
			rw.Header().Set("X-Request-ID", requestID)
			next.ServeHTTP(rw, r)

			logger.Info(
				"http_request",
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.Int("status", rw.statusCode),
				slog.String("remote", r.RemoteAddr),
				slog.Float64("duration_ms", shared.DurationMs(startTime, time.Now())),
				slog.String("request_id", requestID),
			)
		})
	}
}
