package routes

import (
	"dog-service/server/handlers"
	"dog-service/storage"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
)

func SetupRoutes(r *mux.Router, tmplPath, staticPath string, logger *slog.Logger, db *storage.Storage) {
	// staticPath = "/home/kruyneg/Programming/DB_Labs/static"
	r.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "pong")
	})

	r.HandleFunc("/appointments", handlers.NewAppointments(tmplPath, logger, db))
	r.HandleFunc("/profile", handlers.NewProfile(tmplPath, logger, db))
	r.HandleFunc("/profile/save", handlers.NewProfileSaver(logger, db))
	r.HandleFunc("/profile/pet", handlers.NewPetHandler(logger, db))
	r.HandleFunc("/create-appointment", handlers.NewCreateAppointment(tmplPath, logger, db))
	r.HandleFunc("/create-appointment/available-times", handlers.NewAvailableTime(logger, db))
	r.HandleFunc("/reviews", handlers.NewReviews(tmplPath, logger, db))
	r.HandleFunc("/create-review", handlers.NewCreateReview(tmplPath, logger, db))
	r.HandleFunc("/", handlers.NewHome(tmplPath, logger, db))

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(staticPath))))

	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Info("404 Not Found",
			slog.String("url", r.URL.Path),
			slog.String("method", r.Method))
		http.Error(w, "404 - Not Found", http.StatusNotFound)
	})
}
