package userStorage

import (
	"dog-service/models"
	"fmt"
	"github.com/lib/pq"
)

func (s *UserStorage) GetUserData() (models.UserData, error) {
	row := s.db.QueryRow(
		`SELECT
			host.name, host.surname, host.midname,
			host.phone_number, host.email, 
			array_agg(coalesce(pets.id, 0)) as pet_ids,
			array_agg(coalesce(pets.name, '')) as pet_names,
			array_agg(coalesce(pets.breed, '')) as pet_breeds,
			array_agg(coalesce(pets.animal_type, '')) as animal_types
		FROM host LEFT JOIN pets
			ON pets.host_id = host.id
		WHERE host.id = 1
		GROUP BY host.id
	`)

	var (
		res          models.UserData
		pet_ids      []int64
		pet_names    []string
		pet_breeds   []string
		animal_types []string
	)

	err := row.Scan(&res.Name, &res.Surname, &res.Midname, &res.Phone, &res.Email,
		pq.Array(&pet_ids), pq.Array(&pet_names), pq.Array(&pet_breeds), pq.Array(&animal_types))
	if err != nil {
		return models.UserData{},
			fmt.Errorf("error while scanning user data: %s", err)
	}

	for i := range pet_names {
		if pet_names[i] != "" &&
			pet_breeds[i] != "" &&
			animal_types[i] != "" {
			res.Pets = append(res.Pets, models.Pet{
				Id:     pet_ids[i],
				Name:   pet_names[i],
				Breed:  pet_breeds[i],
				Animal: animal_types[i],
			})
		}
	}

	return res, nil
}
