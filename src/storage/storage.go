package storage

import (
	"database/sql"
	appointmentStorage "dog-service/storage/appointment-storage"
	authstorage "dog-service/storage/auth-storage"
	reviewstorage "dog-service/storage/review-storage"
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
	reviewstorage.ReviewStorage
	authstorage.AuthStorage
}

func New(connstr string, cache userStorage.Cache) (Storage, error) {
	db, err := sql.Open("postgres", connstr)
	if err != nil {
		return Storage{}, fmt.Errorf("error while init storage: %s", err)
	}
	// if err := db.Ping(); err != nil {
	// 	return Storage{db}, fmt.Errorf("Error while ping database: %s", err)
	// }
	return Storage{
		AppointmentStorage:  appointmentStorage.New(db),
		SalonGroomerStorage: salonGroomerStorage.New(db),
		UserStorage:         userStorage.New(db, cache),
		ServiceStorage:      serviceStorage.New(db),
		ReviewStorage:       reviewstorage.New(db),
		AuthStorage:         authstorage.New(db),
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
