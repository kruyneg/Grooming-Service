package userStorage

import (
	"dog-service/models"
	"fmt"

	"github.com/lib/pq"
)

func (s *UserStorage) GetUserData(id int64) (models.UserData, error) {
	var user models.UserData
	if ok := s.getFromCache(id, &user); ok {
		var err error
		user.Pets, err = s.getPets(id)
		return user, err
	}

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
		WHERE host.id = $1
		GROUP BY host.id
	`, id)

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

	s.updateCache(user)

	return res, nil
}

func (s *UserStorage) getPets(userID int64) ([]models.Pet, error) {
	rows, err := s.db.Query(`
		SELECT
			id, name, breed, animal_type
		FROM pets
		WHERE host_id = $1
	`, userID)
	if err != nil {
		return []models.Pet{}, err
	}

	var result []models.Pet
	for rows.Next() {
		var pet models.Pet
		if err = rows.Scan(&pet.Id, &pet.Name, &pet.Breed, &pet.Animal);
			err != nil {
			return []models.Pet{}, err
		}
		result = append(result, pet)
	}

	return result, nil
}