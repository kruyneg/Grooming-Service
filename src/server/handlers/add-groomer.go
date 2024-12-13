package handlers

import (
	"dog-service/auth"
	"dog-service/models"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"strconv"
)

type GroomerSaver interface {
	SaveGroomer(models.Groomer) (int64, error)
}

type AddGroomerStorage interface {
	GroomerSaver
	PasswordSaver
	SalonsGetter
}

func NewAddGroomer(tmplPath string, logger *slog.Logger, storage AddGroomerStorage) http.HandlerFunc {
	where := "handlers.add-groomer.NewAddGroomer"

	tmpl, err := template.ParseFiles(tmplPath + "add-groomer.html")
	if err != nil {
		logger.Error("Cannot parse add-groomer.html", slog.String("where", where), slog.String("What", err.Error()))
	}
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			username := r.FormValue("username")
			password := r.FormValue("password")

			var user models.Groomer
			user.Name = r.FormValue("name")
			user.Surname = r.FormValue("surname")
			user.Description.String = r.FormValue("description")

			if r.FormValue("description") == "" {
				user.Description.Valid = false
			} else {
				user.Description.Valid = true
			}

			salons := r.Form["salons[]"]
			for _, idStr := range salons {
				id, err := strconv.Atoi(idStr)
				if err != nil {
					logger.Warn("Cannot parse salon id", slog.String("where", where), slog.String("id", idStr))
					http.Error(w, "Bad salon id", http.StatusBadRequest)
					return
				}
				var salon models.Salon
				salon.Id = int64(id)
				user.Salons = append(user.Salons, salon)
			}

			id, err := storage.SaveGroomer(user)
			if err != nil {
				logger.Error(fmt.Sprintf("%s", err), slog.String("where", where))
				http.Error(w, "Database write error", http.StatusInternalServerError)
				return
			}

			err = storage.SaveLoginPassword(id, username, password, auth.RoleEmployee)
			if err != nil {
				logger.Error(fmt.Sprintf("%s", err), slog.String("where", where))
				http.Error(w, "Database write error", http.StatusInternalServerError)
				return
			}

			http.Redirect(w, r, "/", http.StatusSeeOther)
		}

		data, err := storage.GetSalons()
		if err != nil {
			logger.Error(fmt.Sprintf("%s", err), slog.String("where", where))
			http.Error(w, "Database read error", http.StatusInternalServerError)
			return
		}

		// Выполняем шаблон и передаем данные в ResponseWriter
		err = tmpl.Execute(w, data)
		if err != nil {
			logger.Error("Ошибка при рендеринге шаблона", slog.String("where", where), slog.String("error", err.Error()))
			http.Error(w, "Ошибка при рендеринге страницы", http.StatusInternalServerError)
		}
	}
}
