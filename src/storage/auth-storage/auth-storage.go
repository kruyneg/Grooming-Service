package authstorage

import "database/sql"

type AuthStorage struct {
	db *sql.DB
}

func New(db *sql.DB) AuthStorage {
	return AuthStorage{db}
}