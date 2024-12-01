package salonGroomerStorage

import "database/sql"

type SalonGroomerStorage struct {
	db *sql.DB
}

func New(db *sql.DB) SalonGroomerStorage {
	return SalonGroomerStorage{db}
}

func (s *SalonGroomerStorage) Close() error {
	return s.db.Close()
}
