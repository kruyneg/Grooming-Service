package handlers

import (
	"dog-service/models"
	"encoding/json"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

type ServiceGetter interface {
	GetServices() ([]models.Service, error)
}

type AppointmentSaver interface {
	SaveAppointment(models.Appointment) (int64, error)
}

type SalonMasterGetter interface {
	GetSalonMasters() ([]models.SalonMaster, error)
}

type CreateAppointmentStorage interface {
	ServiceGetter
	UserDataGetter
	AppointmentSaver
	SalonMasterGetter
}

type TimeGetter interface {
	GetTime(time.Time, int64) ([]struct {
		Time     time.Time
		Duration int
	}, error)
}

func NewCreateAppointment(tmplPath string, logger *slog.Logger, storage CreateAppointmentStorage) http.HandlerFunc {
	where := "handlers.create-appointment.NewProfile"

	tmpl, err := template.ParseFiles(tmplPath + "create-appointment.html")
	if err != nil {
		logger.Error("Cannot parse create-appointment.html", slog.String("where", where), slog.String("What", err.Error()))
	}
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			r.ParseForm()
			var res models.Appointment

			res.Service.Id, _ = strconv.ParseInt(r.FormValue("serviceID"), 10, 64)
			res.Pet.Id, _ = strconv.ParseInt(r.FormValue("petID"), 10, 64)
			datePart, _ := time.Parse("2006-01-02", r.FormValue("date")) // Парсим дату
			timePart, _ := time.Parse("15:04", r.FormValue("time"))      // Парсим время
			res.Time = time.Date(
				datePart.Year(),
				datePart.Month(),
				datePart.Day(),
				timePart.Hour(),
				timePart.Minute(),
				0,                   // Секунды
				0,                   // Наносекунды
				timePart.Location(), // Локация (часовой пояс)
			)
			res.SalonMaster.Id, _ = strconv.ParseInt(r.FormValue("salonMasterID"), 10, 64)
			res.Status = "created"

			logger.Debug(fmt.Sprintf("Appointment: %v", res))

			_, err := storage.SaveAppointment(res)
			if err != nil {
				logger.Error(fmt.Sprintf("%s", err), slog.String("where", where))
				http.Error(w, "Database Reading Error", http.StatusInternalServerError)
				return
			}
			http.Redirect(w, r, "/appointments", http.StatusSeeOther)
			return
		}

		// Данные, которые можно передать в шаблон
		var (
			user         models.UserData
			services     []models.Service
			salonMasters []models.SalonMaster
			err          error
		)
		user, err = storage.GetUserData()
		if err != nil {
			logger.Error(fmt.Sprintf("%s", err), slog.String("where", where))
			http.Error(w, "Database Reading Error", http.StatusInternalServerError)
			return
		}
		services, err = storage.GetServices()
		if err != nil {
			logger.Error(fmt.Sprintf("%s", err), slog.String("where", where))
			http.Error(w, "Database Reading Error", http.StatusInternalServerError)
			return
		}
		salonMasters, err = storage.GetSalonMasters()
		if err != nil {
			logger.Error(fmt.Sprintf("%s", err), slog.String("where", where))
			http.Error(w, "Database Reading Error", http.StatusInternalServerError)
			return
		}

		data := struct {
			Pets         []models.Pet
			SalonMasters []models.SalonMaster
			Services     []models.Service
		}{
			Pets:         user.Pets,
			SalonMasters: salonMasters,
			Services:     services,
		}

		// Выполняем шаблон и передаем данные в ResponseWriter
		err = tmpl.Execute(w, data)
		if err != nil {
			logger.Error("Ошибка при рендеринге шаблона", slog.String("where", where), slog.String("error", err.Error()))
			http.Error(w, "Ошибка при рендеринге страницы", http.StatusInternalServerError)
		}
	}
}

func NewAvailableTime(logger *slog.Logger, db TimeGetter) http.HandlerFunc {
	where := "handlers.create-appointment.NewAvailableTime"

	return func(w http.ResponseWriter, r *http.Request) {
		var (
			date          time.Time
			duration      int64
			salonMasterId int64
			err           error
		)
		dateStr := r.URL.Query().Get("date")
		date, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			logger.Warn("Cannot parse date", slog.String("where", where), slog.String("date", dateStr))
			http.Error(w, "Bad date", http.StatusBadRequest)
		}
		durationStr := r.URL.Query().Get("duration")
		duration, err = strconv.ParseInt(durationStr, 10, 32)
		if err != nil {
			logger.Warn("Cannot parse service id", slog.String("where", where), slog.String("duration", durationStr))
			http.Error(w, "Bad serviceId", http.StatusBadRequest)
		}
		salonMasterIdStr := r.URL.Query().Get("salonMasterId")
		salonMasterId, err = strconv.ParseInt(salonMasterIdStr, 10, 64)
		if err != nil {
			logger.Warn("Cannot parse salon-master id", slog.String("where", where), slog.String("salonMasterId", salonMasterIdStr))
			http.Error(w, "Bad salonMasterId", http.StatusBadRequest)
		}

		var notAvailableTimesWithDuration []struct {
			Time     time.Time
			Duration int
		}
		notAvailableTimesWithDuration, err = db.GetTime(date, salonMasterId)
		if err != nil {
			logger.Error("Error while get time", slog.String("error", err.Error()), slog.String("where", where))
			http.Error(w, "Get time error", http.StatusInternalServerError)
			return
		}

		occupied := make(map[int]bool, 0)
		for _, time := range notAvailableTimesWithDuration {
			for i := time.Time.Hour(); i < time.Time.Hour()+time.Duration; i++ {
				occupied[time.Time.Hour()] = true
			}
		}

		var startTime = 8
		if date.YearDay() == time.Now().YearDay() &&
			date.Year() == time.Now().Year() {
			startTime = time.Now().Hour()
		} else if date.Before(time.Now()) {
			startTime = 19
		}

		var res []string
		for i, j := startTime, startTime; j <= 18; i++ {
			j = min(i, j)
			if occupied[i] {
				continue
			}
			for j <= 18 && j-i+1 < int(duration) && !occupied[j] {
				j++
			}
			if i-j+1 == int(duration) && !occupied[j] {
				res = append(res, fmt.Sprintf("%02d:00", i))
			}
		}
		// for i := startTime; i <= 18; i++ {
		// 	if !occupied[i] {
		// 		res = append(res, fmt.Sprintf("%02d:00", i))
		// 	}
		// }
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(res)
	}
}
