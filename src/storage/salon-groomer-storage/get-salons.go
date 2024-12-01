package salonGroomerStorage

import (
	"dog-service/models"
	"fmt"
)

func (s *SalonGroomerStorage) GetSalons() ([]models.Salon, error) {
	rows, err := s.db.Query(`
		SELECT id, address, phone_number
		FROM salons
	`)
	if err != nil {
		return []models.Salon{},
			fmt.Errorf("error while getting salons: %s", err)
	}

	var res []models.Salon
	for rows.Next() {
		var s models.Salon
		rows.Scan(&s.Id, &s.Address, &s.Phone)
		res = append(res, s)
	}
	return res, nil
}