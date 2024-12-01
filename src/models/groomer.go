package models

import "database/sql"

type Groomer struct {
	Id int64
	Name string
	Salons []Salon
	Description sql.NullString
}