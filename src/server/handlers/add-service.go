package handlers

import (
	"dog-service/models"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"strconv"
)

type ServiceSaver interface {
	SaveService(models.Service) (int64, error)
}

func NewAddService(tmplPath string, logger *slog.Logger, storage ServiceSaver) http.HandlerFunc {
	where := "handlers.add-service.NewAddService"

	tmpl, err := template.ParseFiles(tmplPath + "add-service.html")
	if err != nil {
		logger.Error("Cannot parse add-service.html", slog.String("where", where), slog.String("What", err.Error()))
	}
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			var (
				service models.Service
				err     error
				price64 float64
			)

			service.Type = r.FormValue("type")
			price64, _ = strconv.ParseFloat(r.FormValue("price"), 32)
			service.Price = float32(price64)
			service.Duration, _ = strconv.Atoi(r.FormValue("duration"))

			_, err = storage.SaveService(service)
			if err != nil {
				logger.Error(fmt.Sprintf("%s", err), slog.String("where", where))
				http.Error(w, "Database write error", http.StatusInternalServerError)
				return
			}

			http.Redirect(w, r, "/", http.StatusSeeOther)
		}

		// Выполняем шаблон и передаем данные в ResponseWriter
		err = tmpl.Execute(w, nil)
		if err != nil {
			logger.Error("Ошибка при рендеринге шаблона", slog.String("where", where), slog.String("error", err.Error()))
			http.Error(w, "Ошибка при рендеринге страницы", http.StatusInternalServerError)
		}
	}
}
