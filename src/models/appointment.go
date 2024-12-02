package models

import (
	"database/sql"
	"time"
)

type Appointment struct {
	// Id, PetId, SalonMasterId, ServiceId int64
	Id int64
	ReviewId sql.NullInt64
	Status string
	Time time.Time
	// Add info
	Pet Pet
	SalonMaster SalonMaster
	Service Service
}