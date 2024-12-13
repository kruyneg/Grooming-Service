package handlers

import (
	"dog-service/auth"
	"html/template"
	"log/slog"
	"net/http"
)

type PasswordGetter interface {
	GetPassword(login, role string) (int64, string, error)
}

func NewLogin(role, tmplPath string, logger *slog.Logger, storage PasswordGetter) http.HandlerFunc {
	where := "handlers.login.NewLogin"

	tmpl, err := template.ParseFiles(tmplPath + "login.html")
	if err != nil {
		logger.Error("Cannot parse login.html", slog.String("where", where), slog.String("What", err.Error()))
	}
	return func(w http.ResponseWriter, r *http.Request) {
		if (auth.Check(r)) {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		
		data := struct {
			Msg *string
			Role string
		}{Role: role}
		var err error
		var msg string

		if (r.Method == http.MethodPost) {
			var passHash string
			var id int64
			username := r.FormValue("username")
			password := r.FormValue("password")
			
			id, passHash, err = storage.GetPassword(username, role)

			if passHash == "" {
				msg = "Not found login or password"
				data.Msg = &msg
			}

			if err == nil && data.Msg == nil {
				err = auth.CheckPassword(passHash, password)
				if err == nil {
					auth.CreateSession(w, r, id, role)
					http.Redirect(w, r, "/", http.StatusFound)
					return
				}
			}
		}

		if err != nil {
			msg = err.Error()
			data.Msg = &msg
		}

		// Выполняем шаблон и передаем данные в ResponseWriter
		err = tmpl.Execute(w, data)
		if err != nil {
			logger.Error("Ошибка при рендеринге шаблона", slog.String("where", where), slog.String("error", err.Error()))
			http.Error(w, "Ошибка при рендеринге страницы", http.StatusInternalServerError)
		}
	}
}