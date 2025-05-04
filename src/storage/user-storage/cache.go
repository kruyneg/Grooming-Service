package userStorage

import (
	"context"
	"dog-service/models"
	"encoding/json"
	"fmt"
	"log"
)

func (s *UserStorage) updateCache(user models.UserData) {
	data, err := json.Marshal(user)
	if err != nil {
		log.Println(err)
	}

	s.cache.Set(context.Background(), fmt.Sprint(user.Id), string(data))
}

func (s *UserStorage) getFromCache(id int64, user *models.UserData) bool {
	data, err := s.cache.Get(context.Background(), fmt.Sprint(id))
	if err != nil {
		log.Println(err)
		return false
	}
	if err := json.Unmarshal([]byte(data), &user); err != nil {
		log.Println(err)
		return false
	}
	return true
}