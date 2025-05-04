package userStorage

import (
	"context"
	"dog-service/models"
	"dog-service/pubsub"
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

	user.Id = id

	pubsub.Publish(context.Background(), "new_user", user)

	s.updateCache(user)
	return id, nil
}
