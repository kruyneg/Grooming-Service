package userStorage

import (
	"dog-service/models"
	"fmt"
)

func (s *UserStorage) SavePet(data models.Pet) (int64, error) {
	res, err := s.db.Exec(`
		INSERT INTO pets (name, breed, host_id, animal_type)
		VALUES($1, $2, 1, $3)
	`, data.Name, data.Breed, data.Animal)
	if err != nil {
		return 0, fmt.Errorf("error while inserting pet: %s", err)
	}
	id, _ := res.LastInsertId()
	return id, nil
}

func (s *UserStorage) DeletePet(id int64) error {
	_, err := s.db.Exec(`
		DELETE FROM pets WHERE id=$1
	`, id)
	if err != nil {
		return fmt.Errorf("error while deleting pet: %s", err)
	}
	return nil
}
