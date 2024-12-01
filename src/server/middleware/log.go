package middleware

import (
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
)

type responseWriterWrapper struct {
	http.ResponseWriter
	statusCode int
}

func (w *responseWriterWrapper) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

func LogRequest(logger *slog.Logger) mux.MiddlewareFunc {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Info("Request",
				slog.String("url", r.URL.Path),
				slog.String("Method", r.Method))

			wrappedWriter := responseWriterWrapper{ResponseWriter: w, statusCode: http.StatusOK}

			defer func() {
				logger.Info("Response",
					slog.String("Status", http.StatusText(wrappedWriter.statusCode)),
					slog.Int("Code", wrappedWriter.statusCode))
			}()

			h.ServeHTTP(&wrappedWriter, r)
		})
	}
}
