package middleware

import (
	"dog-service/auth"
	"net/http"

	"github.com/gorilla/mux"
)

func RequireRole(role string) mux.MiddlewareFunc {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if auth.GetRole(r) != role {
				http.Error(w, "Not Allowed", http.StatusForbidden)
				return
			}

			h.ServeHTTP(w, r)
		})
	}
}
