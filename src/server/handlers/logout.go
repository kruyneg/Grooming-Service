package handlers

import (
	"dog-service/auth"
	"log/slog"
	"net/http"
)

func NewLogout(logger *slog.Logger) http.HandlerFunc {
	where := "handlers.login.NewLogin"
	return func(w http.ResponseWriter, r *http.Request) {
		if (auth.Check(r)) {
			if err := auth.Logout(w, r); err != nil {
				logger.Error("Error while logout", slog.String("error", err.Error()), slog.String("where", where))
				http.Error(w, "Logout Error", http.StatusInternalServerError)
				return
			}
		}
		
		http.Redirect(w, r, "/", http.StatusFound)
	}
}