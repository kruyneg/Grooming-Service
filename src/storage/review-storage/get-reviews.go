package reviewstorage

import (
	"dog-service/models"
	"fmt"
)

func (s *ReviewStorage) GetReviews() ([]models.Review, error) {
	rows, err := s.db.Query(`
		SELECT
			h.name, h.surname,
			g.name, s.address,
			r.score,
			r.content
		FROM
			review r
		JOIN
			appointments ON appointments.review_id = r.id
		JOIN
			pets ON pets.id = appointments.pet_id
		JOIN
			host h ON pets.host_id = h.id
		JOIN 
			salon_masters ON salon_masters.id = appointments.salon_master_id
		JOIN
			salons s ON salon_masters.salon_id = s.id
		JOIN
			groomers g ON g.id = salon_masters.groomer_id
	`)
	if err != nil {
		return []models.Review{},
			fmt.Errorf("error while get reviews: %s", err)
	}

	var res []models.Review
	for rows.Next() {
		var r models.Review
		rows.Scan(
			&r.Host.Name, &r.Host.Surname,
			&r.SalonMaster.Name, &r.SalonMaster.Address,
			&r.Score, &r.Content,
		)
		res = append(res, r)
	}

	return res, nil
}
