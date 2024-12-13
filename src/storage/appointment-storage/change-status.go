package appointmentstorage

import (
	"fmt"
)

func (s *AppointmentStorage) ChangeStatus(id int64, status string) error {
	_, err := s.db.Exec(`
		UPDATE appointments
		SET status = $2
		WHERE id = $1
	`, id, status)
	if err != nil {
		return fmt.Errorf("error while save appointment: %s", err)
	}
	return nil
}
