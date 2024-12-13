package handlers

import (
	"dog-service/auth"
	"dog-service/models"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
)

type AppointmentGetter interface {
	GetUAppointments(int64) ([]models.Appointment, error)
	GetEAppointments(int64) ([]models.Appointment, error)
}

func NewUAppointments(tmplPath string, logger *slog.Logger, getter AppointmentGetter) http.HandlerFunc {
	where := "handlers.appointments.NewUAppointments"

	tmpl, err := template.ParseFiles(tmplPath + "appointments.html")
	if err != nil {
		logger.Error("Cannot parse appointments.html", slog.String("where", where), slog.String("What", err.Error()))
	}
	return func(w http.ResponseWriter, r *http.Request) {

		userId := auth.GetId(r)
		// Данные, которые можно передать в шаблон
		data, err := getter.GetUAppointments(userId)
		if err != nil {
			logger.Error(fmt.Sprintf("%s", err), slog.String("where", where))
			http.Error(w, "Database Reading Error", http.StatusInternalServerError)
			return
		}

		logger.Info(fmt.Sprint(data), slog.String("where", where))

		// Выполняем шаблон и передаем данные в ResponseWriter
		err = tmpl.Execute(w, data)
		if err != nil {
			logger.Error("Ошибка при рендеринге шаблона", slog.String("where", where), slog.String("error", err.Error()))
			http.Error(w, "Ошибка при рендеринге страницы", http.StatusInternalServerError)
		}
	}

}
