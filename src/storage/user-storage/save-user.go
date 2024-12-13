package userStorage

import (
	"dog-service/models"
	"fmt"
)

func (s *UserStorage) SaveUser(user models.UserData) (int64, error) {
	res := s.db.QueryRow(`
		INSERT INTO host (name, surname, midname, phone_number)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`, user.Name, user.Surname, user.Midname, user.Phone)
	if res.Err() != nil {
		return 0,
			fmt.Errorf("error while save user: %s", res.Err())
	}

	var id int64
	res.Scan(&id)
	return id, nil
}
