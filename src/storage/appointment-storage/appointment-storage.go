package appointmentstorage

import "database/sql"

type AppointmentStorage struct {
	db *sql.DB
}

func New(db *sql.DB) AppointmentStorage {
	return AppointmentStorage{db}
}

func (s *AppointmentStorage) Close() error {
	return s.db.Close()
}
