package models

import (
	"database/sql"
	"time"
)

type Appointment struct {
	Id, PetId, SalonMasterId, ServiceId int64
	ReviewId sql.NullInt64
	Status string
	Time time.Time
}