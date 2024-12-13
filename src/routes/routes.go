package routes

import (
	"dog-service/auth"
	"dog-service/server/handlers"
	"dog-service/server/middleware"
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

	// setup middlewares
	r.Use(middleware.LogRequest(logger))
	userRouter := r.PathPrefix("/u").Subrouter()
	employeeRouter := r.PathPrefix("/e").Subrouter()
	adminRouter := r.PathPrefix("/a").Subrouter()
	userRouter.Use(middleware.Auth(logger))
	employeeRouter.Use(middleware.Auth(logger))
	employeeRouter.Use(middleware.RequireRole(auth.RoleEmployee))
	adminRouter.Use(middleware.Auth(logger))
	adminRouter.Use(middleware.RequireRole(auth.RoleAdmin))

	// setup routes
	userRouter.HandleFunc("/appointments", handlers.NewUAppointments(tmplPath, logger, db))
	userRouter.HandleFunc("/profile", handlers.NewProfile(tmplPath, logger, db))
	userRouter.HandleFunc("/profile/save", handlers.NewProfileSaver(logger, db))
	userRouter.HandleFunc("/profile/pet", handlers.NewPetHandler(logger, db))
	userRouter.HandleFunc("/create-appointment", handlers.NewCreateAppointment(tmplPath, logger, db))
	userRouter.HandleFunc("/create-appointment/available-times", handlers.NewAvailableTime(logger, db))
	userRouter.HandleFunc("/create-review", handlers.NewCreateReview(tmplPath, logger, db))
	employeeRouter.HandleFunc("/appointments", handlers.NewEAppointments(tmplPath, logger, db))
	adminRouter.HandleFunc("/add-groomer", handlers.NewAddGroomer(tmplPath, logger, db))
	adminRouter.HandleFunc("/add-service", handlers.NewAddService(tmplPath, logger, db))
	r.HandleFunc("/reviews", handlers.NewReviews(tmplPath, logger, db))
	r.HandleFunc("/", handlers.NewHome(tmplPath, logger, db))
	r.HandleFunc("/u/login", handlers.NewLogin(auth.RoleUser, tmplPath, logger, db))
	r.HandleFunc("/e/login", handlers.NewLogin(auth.RoleEmployee, tmplPath, logger, db))
	r.HandleFunc("/a/login", handlers.NewLogin(auth.RoleAdmin, tmplPath, logger, db))
	r.HandleFunc("/u/register", handlers.NewRegister(auth.RoleUser, tmplPath, logger, db))
	r.HandleFunc("/logout", handlers.NewLogout(logger))

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(staticPath))))

	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Info("404 Not Found",
			slog.String("url", r.URL.Path),
			slog.String("method", r.Method))
		http.Error(w, "404 - Not Found", http.StatusNotFound)
	})
}
