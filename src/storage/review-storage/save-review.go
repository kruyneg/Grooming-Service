package reviewstorage

import (
	"dog-service/models"
	"fmt"
)

func (s *ReviewStorage) SaveReview(appointmentId int64, review models.Review) (id int64, err error) {
    tx, err := s.db.Begin()
    if err != nil {
        return 0, fmt.Errorf("cannot start transaction: %w", err)
    }
    defer func() {
        if err != nil {
            tx.Rollback()
        } else {
            err = tx.Commit()
        }
    }()

    // 1. Вставка отзыва в таблицу reviews
    query := `INSERT INTO review (score, content) VALUES ($1, $2) RETURNING id`
    err = tx.QueryRow(query, review.Score, review.Content).Scan(&id)
    if err != nil {
        return 0, fmt.Errorf("cannot insert review: %w", err)
    }

    // 2. Обновление review_id в таблице appointments
    updateQuery := `UPDATE appointments SET review_id = $1 WHERE id = $2`
    _, err = tx.Exec(updateQuery, id, appointmentId)
    if err != nil {
        return 0, fmt.Errorf("cannot update appointment: %w", err)
    }

    return id, nil
}
