package servicestorage

import "dog-service/models"

func (s *ServiceStorage) SaveService(service models.Service) (int64, error) {
	err := s.db.QueryRow(`
		INSERT INTO services (type, price, duration)
		VALUES ($1, $2, $3)
		RETURNING id
	`, service.Type, service.Price, service.Duration).Scan(&service.Id)

	return service.Id, err
}