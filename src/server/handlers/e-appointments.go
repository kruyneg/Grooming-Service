package handlers

import (
	"dog-service/auth"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"strconv"
)

type StatusChanger interface {
	AppointmentGetter
	ChangeStatus(int64, string) (error)
}

func NewEAppointments(tmplPath string, logger *slog.Logger, storage StatusChanger) http.HandlerFunc {
	where := "handlers.appointments.NewEAppointments"

	tmpl, err := template.ParseFiles(tmplPath + "e-appointments.html")
	if err != nil {
		logger.Error("Cannot parse e-appointments.html", slog.String("where", where), slog.String("What", err.Error()))
	}
	return func(w http.ResponseWriter, r *http.Request) {

		userId := auth.GetId(r)
		// Данные, которые можно передать в шаблон

		if r.Method == http.MethodPost {
			idStr := r.URL.Query().Get("aid")
			id, err := strconv.ParseInt(idStr, 10, 64)
			if err != nil {
				logger.Warn(fmt.Sprintf("%s", err), slog.String("where", where))
				http.Error(w, "Bad appointment id", http.StatusBadRequest)
				return
			}

			err = storage.ChangeStatus(id, r.FormValue("status"))
			if err != nil {
				logger.Error(fmt.Sprintf("%s", err), slog.String("where", where))
				http.Error(w, "Database Saving Error", http.StatusInternalServerError)
				return	
			}
		}

		data, err := storage.GetEAppointments(userId)
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
