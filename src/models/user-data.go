package models

import "database/sql"

type UserData struct {
	Name, Surname string
	Midname sql.NullString
	Phone, Email string
	Pets []Pet
}

type Pet struct {
	Id int64
	Name, Breed, Animal string
	Weight sql.NullInt32
}