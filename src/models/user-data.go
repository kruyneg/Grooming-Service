package models

import "database/sql"

type UserData struct {
	Id int64
	Name, Surname string
	Midname sql.NullString
	Phone string
	Email sql.NullString
	Pets []Pet
}

type Pet struct {
	Id int64
	Name, Breed, Animal string
	Weight sql.NullInt32
}