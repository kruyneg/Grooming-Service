package userStorage

import "database/sql"

type UserStorage struct {
	db *sql.DB
}

func New(db * sql.DB) UserStorage {
	return UserStorage{db}
}

func (s *UserStorage) Close() error {
	return s.db.Close()
}
