package reviewstorage

import "database/sql"

type ReviewStorage struct {
	db *sql.DB
}

func New(db *sql.DB) ReviewStorage {
	return ReviewStorage{db}
}