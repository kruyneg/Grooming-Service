package storage

import (
	"database/sql"
	appointmentStorage "dog-service/storage/appointment-storage"
	salonGroomerStorage "dog-service/storage/salon-groomer-storage"
	serviceStorage "dog-service/storage/service-storage"
	userStorage "dog-service/storage/user-storage"
	"fmt"

	_ "github.com/lib/pq"
)

type Storage struct {
	appointmentStorage.AppointmentStorage
	salonGroomerStorage.SalonGroomerStorage
	userStorage.UserStorage
	serviceStorage.ServiceStorage
}

func New(connstr string) (Storage, error) {
	db, err := sql.Open("postgres", connstr)
	if err != nil {
		return Storage{}, fmt.Errorf("error while init storage: %s", err)
	}
	// if err := db.Ping(); err != nil {
	// 	return Storage{db}, fmt.Errorf("Error while ping database: %s", err)
	// }
	return Storage{
		AppointmentStorage: appointmentStorage.New(db),
		SalonGroomerStorage: salonGroomerStorage.New(db),
		UserStorage: userStorage.New(db),
		ServiceStorage: serviceStorage.New(db),
	}, nil
}

func (s *Storage) Close() error {
	if err := s.AppointmentStorage.Close(); err != nil {
		return err
	} else if err := s.SalonGroomerStorage.Close(); err != nil {
		return err
	} else if err := s.UserStorage.Close(); err != nil {
		return err
	} else if err := s.ServiceStorage.Close(); err != nil {
		return err
	}
	return nil
}
