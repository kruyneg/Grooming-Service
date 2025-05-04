package userStorage

import (
	"context"
	"dog-service/models"
	"dog-service/pubsub"
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
		WHERE id = $6
	`, data.Name, data.Surname, data.Midname, data.Phone, data.Email, data.Id)
	if err != nil {
		return fmt.Errorf("error while update user data: %s", err)
	}

	pubsub.Publish(context.Background(), "user_changed", data)

	s.updateCache(data)
	return nil
}
