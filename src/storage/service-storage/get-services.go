package servicestorage

import (
	"dog-service/models"
	"fmt"
)

func (s *ServiceStorage) GetServices() ([]models.Service, error) {
	rows, err := s.db.Query(`
		SELECT
			id, type, price, duration
		FROM services
	`)
	if err != nil {
		return []models.Service{},
			fmt.Errorf("error while selecting services: %s", err)
	}
	defer rows.Close()

	var res []models.Service
	for rows.Next() {
		var service models.Service
		rows.Scan(&service.Id, &service.Type, &service.Price, &service.Duration)
		res = append(res, service)
	}
	return res, nil
}
