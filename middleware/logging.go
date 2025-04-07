package middleware

import (
	"log/slog"
	"net/http"
	"time"

	"ddg-search/config"
)

// LoggingMiddleware logs incoming requests and their responses.
func LoggingMiddleware(cfg *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if cfg.DebugMode {
				start := time.Now()

				// Create a response recorder to capture the status code and size.
				recorder := &responseRecorder{
					ResponseWriter: w,
					statusCode:     http.StatusOK,
				}

				// Call the next handler in the chain.
				next.ServeHTTP(recorder, r)

				duration := time.Since(start)

				slog.Debug("Request processed",
					"method", r.Method,
					"path", r.URL.Path,
					"remote_addr", r.RemoteAddr,
					"duration", duration,
					"status_code", recorder.statusCode,
					"size", recorder.size,
				)
			} else {
				// If debug mode is not enabled, just call the next handler.
				next.ServeHTTP(w, r)
			}
		})
	}
}

// responseRecorder is a simple wrapper to capture the status code and size of the response.
type responseRecorder struct {
	http.ResponseWriter
	statusCode int
	size       int
}

func (r *responseRecorder) WriteHeader(statusCode int) {
	r.statusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

func (r *responseRecorder) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.size += size
	return size, err
}
