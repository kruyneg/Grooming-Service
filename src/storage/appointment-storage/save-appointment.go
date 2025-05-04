package appointmentstorage

import (
	"context"
	"dog-service/models"
	"dog-service/pubsub"
	"fmt"
)

func (s *AppointmentStorage) SaveAppointment(
	appointment models.Appointment) (int64, error) {

	res, err := s.db.Exec(`
		INSERT INTO appointments
			(pet_id, service_id, time, salon_master_id, status)
		VALUES
			($1, $2, $3, $4, $5)
	`, appointment.Pet.Id, appointment.Service.Id,
	appointment.Time, appointment.SalonMaster.Id, appointment.Status)
	if err != nil {
		return 0,
			fmt.Errorf("error while save appointment: %s", err)
	}
	id, _ := res.LastInsertId()
	pubsub.Publish(context.Background(), "new_appointment",
		map[string]any{
			"id": appointment.Id,
			"pet_id": appointment.Pet.Id,
			"service_id": appointment.Service.Id,
			"time": appointment.Time,
			"status": appointment.Status,
		})
	return id, nil
}
