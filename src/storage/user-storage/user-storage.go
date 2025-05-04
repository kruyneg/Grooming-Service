package userStorage

import (
	"context"
	"database/sql"
)

type Cache interface {
	Set(context.Context, string, string) error
	Get(context.Context, string) (string, error)
	Delete(context.Context, string) error
}

type UserStorage struct {
	db *sql.DB
	cache Cache
}

func New(db *sql.DB, cache Cache) UserStorage {
	return UserStorage{db, cache}
}

func (s *UserStorage) Close() error {
	return s.db.Close()
}
