// Package middlewares provides reusable HTTP middleware for the application.
package middlewares

import (
	"log/slog"
	"net/http"
	"time"
)

// responseRecorder captures the status code written by the handler.
type responseRecorder struct {
	http.ResponseWriter
	statusCode int
	written    bool
}

func newResponseRecorder(writer http.ResponseWriter) *responseRecorder {
	return &responseRecorder{
		ResponseWriter: writer,
		statusCode:     http.StatusOK,
	}
}

func (rr *responseRecorder) WriteHeader(code int) {
	if !rr.written {
		rr.statusCode = code
		rr.written = true
		rr.ResponseWriter.WriteHeader(code)
	}
}

func (rr *responseRecorder) Write(b []byte) (int, error) {
	if !rr.written {
		rr.WriteHeader(http.StatusOK)
	}
	return rr.ResponseWriter.Write(b)
}

// Logger logs details about each completed HTTP request.
func Logger(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			start := time.Now()
			rec := newResponseRecorder(writer)

			// Call the next handler.
			next.ServeHTTP(rec, request)

			logger.InfoContext(
				request.Context(),
				"request completed",
				slog.String("method", request.Method),
				slog.String("path", request.URL.Path),
				slog.Int("status", rec.statusCode),
				slog.Duration("duration", time.Since(start)),
				slog.String("remote_addr", request.RemoteAddr),
			)
		})
	}
}
