package appointmentstorage

import (
	"dog-service/models"
	"fmt"
)

func (s *AppointmentStorage) SaveAppointment(
	appointment models.Appointment) (int64, error) {

	res, err := s.db.Exec(`
		INSERT INTO appointments
			(pet_id, service_id, time, salon_master_id, status)
		VALUES
			($1, $2, $3, $4, $5)
	`, appointment.PetId, appointment.ServiceId,
	appointment.Time, appointment.SalonMasterId, appointment.Status)
	if err != nil {
		return 0,
			fmt.Errorf("error while save appointment: %s", err)
	}
	id, _ := res.LastInsertId()
	return id, nil
}
