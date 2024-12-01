package appointmentstorage

import (
	"dog-service/models"
	"fmt"
)

func (s *AppointmentStorage) GetAppointments() ([]models.Appointment, error) {
	rows, err := s.db.Query("SELECT status, time FROM appointments")
	if err != nil {
		return []models.Appointment{},
			fmt.Errorf("error while getting appointments: %s", err)
	}
	defer rows.Close()
	var res []models.Appointment
	for rows.Next() {
		var a models.Appointment
		err = rows.Scan(&a.Status, &a.Time)
		if err != nil {
			return []models.Appointment{},
			fmt.Errorf("error while scaning appointments: %s", err)
		}
		res = append(res, a)
	}
	return res, nil
}
