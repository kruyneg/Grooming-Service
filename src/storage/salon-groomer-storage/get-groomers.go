package salonGroomerStorage

import (
	"dog-service/models"
	"fmt"

	"github.com/lib/pq"
)

func (s *SalonGroomerStorage) GetGroomers() ([]models.Groomer, error) {
	rows, err := s.db.Query(`
		SELECT 
			g.id, g.name, g.description,
			array_agg(s.address) as addresses,
			array_agg(s.phone_number) as phones
		FROM
			salon_masters
		JOIN
			salons as s ON salon_masters.salon_id = s.id
		JOIN
			groomers as g ON salon_masters.groomer_id = g.id
		GROUP BY g.id, g.name, g.description
	`)
	if err != nil {
		return []models.Groomer{},
			fmt.Errorf("error while getting groomers: %s", err)
	}

	var res []models.Groomer
	for rows.Next() {
		var (
			g         models.Groomer
			addresses []string
			phones    []string
		)

		rows.Scan(&g.Id, &g.Name, &g.Description, pq.Array(&addresses), pq.Array(&phones))

		for i := range addresses {
			g.Salons = append(g.Salons, models.Salon{Address: addresses[i], Phone: phones[i]})
		}
		res = append(res, g)
	}
	return res, nil
}
