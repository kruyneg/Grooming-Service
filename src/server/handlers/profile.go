package handlers

import (
	"dog-service/auth"
	"dog-service/models"
	"encoding/json"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"strconv"
)

type UserDataGetter interface {
	GetUserData(int64) (models.UserData, error)
}
type UserDataSetter interface {
	SetUserData(models.UserData) error
}
type PetStorage interface {
	DeletePet(int64) error
	SavePet(int64, models.Pet) (int64, error)
}

func NewProfile(tmplPath string, logger *slog.Logger, getter UserDataGetter) http.HandlerFunc {
	where := "handlers.profile.NewProfile"

	tmpl, err := template.ParseFiles(tmplPath + "profile.html")
	if err != nil {
		logger.Error("Cannot parse profile.html", slog.String("where", where), slog.String("What", err.Error()))
	}
	return func(w http.ResponseWriter, r *http.Request) {
		userId := auth.GetId(r)

		// Данные, которые можно передать в шаблон
		data, err := getter.GetUserData(userId)
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

func NewProfileSaver(logger *slog.Logger, saver UserDataSetter) http.HandlerFunc {
	where := "handlers.profile.NewProfileSaver"

	return func(w http.ResponseWriter, r *http.Request) {
		userId := auth.GetId(r)

		if r.Method == http.MethodPut {
			r.ParseForm()

			user := models.UserData{
				Name:    r.FormValue("name"),
				Surname: r.FormValue("surname"),
				Phone:   r.FormValue("phone"),
				Id:      userId,
			}
			if r.FormValue("midname") != "" {
				user.Midname.String = r.FormValue("midname")
				user.Midname.Valid = true
			}
			if r.FormValue("email") != "" {
				user.Email.String = r.FormValue("email")
				user.Email.Valid = true
			}

			if err := saver.SetUserData(user); err != nil {
				logger.Error("Error while save user data", slog.String("error", err.Error()), slog.String("where", where))
				http.Error(w, "Save Error", http.StatusInternalServerError)
				return
			}
			logger.Info("Save user", slog.String("where", where), slog.Int("host_id", 1))
		}

		// Переадресуем пользователя на страницу профиля с успешным сообщением
		http.Redirect(w, r, "/u/profile", http.StatusSeeOther)
	}
}

func NewPetHandler(logger *slog.Logger, storage PetStorage) http.HandlerFunc {
	where := "handlers.profile.NewPetHandler"

	return func(w http.ResponseWriter, r *http.Request) {
		userId := auth.GetId(r)

		logger.Debug("Step into /profile/pet")
		if r.Method == http.MethodPost {
			r.ParseForm()

			pet := models.Pet{
				Name:   r.FormValue("petName"),
				Breed:  r.FormValue("petBreed"),
				Animal: r.FormValue("petAnimal"),
			}
			id, err := storage.SavePet(userId, pet)
			if err != nil {
				logger.Error("Error while save pet", slog.String("error", err.Error()), slog.String("where", where))
				http.Error(w, "Save Error", http.StatusInternalServerError)
				return
			}
			logger.Info("Save pet", slog.String("where", where), slog.Int64("pet_id", id))
		} else if r.Method == http.MethodDelete {
			data := struct {
				Id string `json:"petID"`
			}{}
			if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
				logger.Warn("Error while parse json", slog.String("where", where), slog.String("error", err.Error()))
				http.Error(w, "Invalid request body", http.StatusBadRequest)
				return
			}
			logger.Debug("", slog.Any("data", data))
			id, _ := strconv.ParseInt(data.Id, 10, 64)

			if err := storage.DeletePet(id); err != nil {
				logger.Error("Error while delete pet", slog.String("error", err.Error()), slog.String("where", where))
				http.Error(w, "Delete Error", http.StatusInternalServerError)
				return
			}
			logger.Info("Delete pet", slog.String("where", where), slog.Int64("pet_id", id))
		}

		// Переадресуем пользователя на страницу профиля с успешным сообщением
		http.Redirect(w, r, "/u/profile", http.StatusSeeOther)
	}
}
