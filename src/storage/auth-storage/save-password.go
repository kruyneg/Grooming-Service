package authstorage

import (
	"dog-service/auth"
	"fmt"
)

func (s *AuthStorage) SaveLoginPassword(id int64, login, password, role string) error {
	hashpass, err := auth.HashPassword(password)
	if err != nil {
		return err
	}
	res, e := s.db.Exec(`
		INSERT INTO auth (user_id, login, password_hash, role)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (login) DO NOTHING
	`, id, login, hashpass, role)

	if e != nil {
		return e
	}

	if cnt, _ := res.RowsAffected(); cnt == 0 {
		return fmt.Errorf("this login already exists")
	}
	return e
}