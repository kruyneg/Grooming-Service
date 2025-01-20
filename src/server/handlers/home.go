package handlers

import (
	"dog-service/auth"
	"dog-service/models"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
)

type GroomersGetter interface {
	GetGroomers() ([]models.Groomer, error)
}

type SalonsGetter interface {
	GetSalons() ([]models.Salon, error)
}

type HomeStorage interface {
	GroomersGetter
	ServiceGetter
	SalonsGetter
}

func NewHome(tmplPath string, logger *slog.Logger, getter HomeStorage) http.HandlerFunc {
	where := "handlers.home.NewHome"

	tmpl, err := template.ParseFiles(tmplPath + "home.html")
	if err != nil {
		logger.Error("Cannot parse home.html", slog.String("where", where), slog.String("What", err.Error()))
	}
	return func(w http.ResponseWriter, r *http.Request) {
		// Данные, которые можно передать в шаблон
		var data = struct {
			Services []models.Service
			Groomers []models.Groomer
			Role string
		}{Role: "none"}

		if auth.Check(r) {
			data.Role = auth.GetRole(r)
		}
		var err error

		data.Services, err = getter.GetServices()
		if err != nil {
			logger.Error(fmt.Sprintf("%s", err), slog.String("where", where))
			http.Error(w, "Database Reading Error", http.StatusInternalServerError)
			return
		}
		data.Groomers, err = getter.GetGroomers()
		if err != nil {
			logger.Error(fmt.Sprintf("%s", err), slog.String("where", where))
			http.Error(w, "Database Reading Error", http.StatusInternalServerError)
		}
		// data.Salons, err = getter.GetSalons()
		// if err != nil {
		// 	logger.Error(fmt.Sprintf("%s", err), slog.String("where", where))
		// 	http.Error(w, "Database Reading Error", http.StatusInternalServerError)
		// }

		logger.Info(fmt.Sprint(data), slog.String("where", where))

		// Выполняем шаблон и передаем данные в ResponseWriter
		err = tmpl.Execute(w, data)
		if err != nil {
			logger.Error("Ошибка при рендеринге шаблона", slog.String("where", where), slog.String("error", err.Error()))
			http.Error(w, "Ошибка при рендеринге страницы", http.StatusInternalServerError)
		}
	}
}
