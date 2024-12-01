package servicestorage

import "database/sql"

type ServiceStorage struct {
	db* sql.DB
}

func New(db *sql.DB) ServiceStorage {
	return ServiceStorage{db}
}

func (s *ServiceStorage) Close() error {
	return s.db.Close()
}
