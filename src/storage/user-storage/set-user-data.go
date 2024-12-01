package userStorage

import (
	"dog-service/models"
	"fmt"
)

func (s *UserStorage) SetUserData(data models.UserData) error {
	_, err := s.db.Exec(`
		UPDATE host
		SET
			name=$1,
			surname=$2,
			midname=$3,
			phone_number=$4,
			email=$5
		WHERE id = 1
	`, data.Name, data.Surname, data.Midname, data.Phone, data.Email)
	if err != nil {
		return fmt.Errorf("error while update user data: %s", err)
	}
	return nil
}
