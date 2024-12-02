package appointmentstorage

import (
	"fmt"
	"time"
)

func (s *AppointmentStorage) GetTime(date time.Time, id int64) ([]struct {
	Time     time.Time
	Duration int
}, error) {
	rows, err := s.db.Query(`
		SELECT appointments.time, services.duration
		FROM appointments
			JOIN services ON services.id = appointments.service_id
		WHERE salon_master_id = $1 AND
			DATE_PART('year', time) = $2 AND
			DATE_PART('month', time) = $3 AND
			DATE_PART('day', time) = $4
	`, id, date.Year(), date.Month(), date.Day())
	if err != nil {
		return []struct {
				Time     time.Time
				Duration int
			}{},
			fmt.Errorf("error while select time: %s", err)
	}
	defer rows.Close()

	var res []struct {
		Time     time.Time
		Duration int
	}

	for rows.Next() {
		var t struct {
			Time     time.Time
			Duration int
		}
		rows.Scan(&t.Time, &t.Duration)
		res = append(res, t)
	}
	return res, nil
}
