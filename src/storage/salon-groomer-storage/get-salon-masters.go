package salonGroomerStorage

import (
	"dog-service/models"
	"fmt"
)

func (s SalonGroomerStorage) GetSalonMasters() ([]models.SalonMaster, error) {
	rows, err := s.db.Query(`
		SELECT 
			salon_masters.id,
			groomers.name,
			salons.address
		FROM
			salon_masters
		JOIN
			salons ON salons.id = salon_masters.salon_id
		JOIN
			groomers ON groomers.id = salon_masters.groomer_id
	`)
	if err != nil {
		return []models.SalonMaster{},
			fmt.Errorf("error while reading salon_masters: %s", err)
	}

	var res []models.SalonMaster
	for rows.Next() {
		var sm models.SalonMaster
		rows.Scan(&sm.Id, &sm.Name, &sm.Address)
		res = append(res, sm)
	}
	return res, nil
}