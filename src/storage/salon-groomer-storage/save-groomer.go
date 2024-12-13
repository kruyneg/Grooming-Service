package salonGroomerStorage

import (
	"dog-service/models"
	"fmt"
)

func (s *SalonGroomerStorage) SaveGroomer(groomer models.Groomer) (int64, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return 0, fmt.Errorf("could not begin transaction: %s", err)
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	err = tx.QueryRow("INSERT INTO groomers (name, surname, description) VALUES ($1, $2, $3) RETURNING id",
		groomer.Name, groomer.Surname, groomer.Description).Scan(&groomer.Id)
	if err != nil {
		return 0, fmt.Errorf("could not insert into groomers: %w", err)
	}

	for _, salon := range groomer.Salons {
		_, err = tx.Exec("INSERT INTO salon_masters (groomer_id, salon_id) VALUES ($1, $2)", groomer.Id, salon.Id)
		if err != nil {
			return 0, fmt.Errorf("could not insert into salon_masters: %w", err)
		}
	}

	return groomer.Id, nil
}
