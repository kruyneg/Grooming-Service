package handlers

import (
	"dog-service/models"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
)

type UserSaver interface {
	SaveUser(models.UserData) (int64, error)
}
type PasswordSaver interface {
	SaveLoginPassword(id int64, login, password, role string) error
}

type RegisterStorage interface {
	UserSaver
	PasswordSaver
}

func NewRegister(role, tmplPath string, logger *slog.Logger, storage RegisterStorage) http.HandlerFunc {
	where := "handlers.registration.NewRegister"

	tmpl, err := template.ParseFiles(tmplPath + "registration.html")
	if err != nil {
		logger.Error("Cannot parse registration.html", slog.String("where", where), slog.String("What", err.Error()))
	}
	return func(w http.ResponseWriter, r *http.Request) {
		if (r.Method == http.MethodPost) {
			username := r.FormValue("username")
			password := r.FormValue("password")

			var user models.UserData
			user.Name = r.FormValue("name")
			user.Surname = r.FormValue("surname")
			user.Midname.String = r.FormValue("midname")
			user.Phone = r.FormValue("phone")

			if user.Midname.String != "" {
				user.Midname.Valid = true
			}
			
			id, err := storage.SaveUser(user)
			if err != nil {
				logger.Error(fmt.Sprintf("%s", err), slog.String("where", where))
				http.Error(w, "Database write error", http.StatusInternalServerError)
				return
			}

			err = storage.SaveLoginPassword(id, username, password, role)
			if err != nil {
				logger.Error(fmt.Sprintf("%s", err), slog.String("where", where))
				http.Error(w, "Database write error", http.StatusInternalServerError)
				return
			}
		}

		// Выполняем шаблон и передаем данные в ResponseWriter
		err = tmpl.Execute(w, nil)
		if err != nil {
			logger.Error("Ошибка при рендеринге шаблона", slog.String("where", where), slog.String("error", err.Error()))
			http.Error(w, "Ошибка при рендеринге страницы", http.StatusInternalServerError)
		}
	}
}