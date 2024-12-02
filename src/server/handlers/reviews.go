package handlers

import (
	"dog-service/models"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"strconv"
)

type ReviewGetter interface {
	GetReviews() ([]models.Review, error)
}

type ReviewSaver interface {
	SaveReview(int64, models.Review) (int64, error)
}

func makeRange(n int) []int {
	result := make([]int, n)
	for i := range result {
		result[i] = i
	}
	return result
}

func sub(a, b int) int {
	return a - b
}

func NewReviews(tmplPath string, logger *slog.Logger, getter ReviewGetter) http.HandlerFunc {
	where := "handlers.reviews.NewReviews"

	tmpl := template.New("reviews.html").Funcs(template.FuncMap{
		"makeRange": makeRange,
		"sub":       sub,
	})
	tmpl, err := tmpl.ParseFiles(tmplPath + "reviews.html")
	if err != nil {
		logger.Error("Cannot parse reviews.html", slog.String("where", where), slog.String("What", err.Error()))
	}
	return func(w http.ResponseWriter, r *http.Request) {

		// Данные, которые можно передать в шаблон
		data, err := getter.GetReviews()
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

func NewCreateReview(tmplPath string, logger *slog.Logger, saver ReviewSaver) http.HandlerFunc {
	where := "handlers.reviews.NewCreateReview"

	tmpl, err := template.ParseFiles(tmplPath + "create-review.html")
	if err != nil {
		logger.Error("Cannot parse create-review.html", slog.String("where", where), slog.String("What", err.Error()))
	}
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			appointmentIdStr := r.URL.Query().Get("aid")
			appointmentId, err := strconv.ParseInt(appointmentIdStr, 10, 64)
			if err != nil {
				logger.Warn("Cannot parse appointment id", slog.String("where", where), slog.String("aid", appointmentIdStr))
				http.Error(w, "Bad appointment id", http.StatusBadRequest)
				return
			}

			var res models.Review

			res.Score, err = strconv.Atoi(r.FormValue("rating"))
			if err != nil {
				logger.Warn("Cannot parse rating", slog.String("where", where), slog.String("rating", r.FormValue("rating")))
				http.Error(w, "Bad rating", http.StatusBadRequest)
				return
			}
			res.Content = r.FormValue("review_text")

			saver.SaveReview(appointmentId, res)

			http.Redirect(w, r, "/reviews", http.StatusSeeOther)
		} else {
			tmpl.Execute(w, nil)
		}
	}
}
