package appointmentstorage

import (
	"dog-service/models"
	"fmt"
)

func (s *AppointmentStorage) GetAppointments(hostId int64) ([]models.Appointment, error) {
	rows, err := s.db.Query(`
		SELECT
			a.id, a.status, a.time,
			pets.name, pets.breed,
			salons.address,
			groomers.name,
			services.type
		FROM
			appointments a
		FULL JOIN
			pets ON a.pet_id = pets.id
		FULL JOIN
			salon_masters ON salon_masters.id = a.salon_master_id
		FULL JOIN
			salons ON salon_masters.salon_id = salons.id
		FULL JOIN
			groomers ON salon_masters.groomer_id = groomers.id
		FULL JOIN
			services ON a.service_id = services.id
		FULL JOIN 
			host ON host.id = pets.host_id
		WHERE host.id = $1
	`, hostId)
	if err != nil {
		return []models.Appointment{},
			fmt.Errorf("error while getting appointments: %s", err)
	}
	defer rows.Close()
	var res []models.Appointment
	for rows.Next() {
		var a models.Appointment
		err = rows.Scan(
			&a.Id, &a.Status, &a.Time,
			&a.Pet.Name, &a.Pet.Breed,
			&a.SalonMaster.Address, &a.SalonMaster.Name,
			&a.Service.Type)
		if err != nil {
			return []models.Appointment{},
				fmt.Errorf("error while scaning appointments: %s", err)
		}
		res = append(res, a)
	}
	return res, nil
}
