package middleware

import (
	"dog-service/auth"
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
)

func Auth(logger *slog.Logger) mux.MiddlewareFunc {
	where := "middleware.auth.Auth"
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !auth.Check(r) {
				http.Redirect(w, r, "/u/login", http.StatusFound)
				logger.Debug("User need login", slog.String("where", where))
				return
			}
			h.ServeHTTP(w, r)
		})
	}
}
