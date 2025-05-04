package models

import "database/sql"

type UserData struct {
	Id int64 `json:"id"`
	Name string `json:"name"`
	Surname string `json:"surname"`
	Midname sql.NullString `json:"midname"`
	Phone string `json:"phone"`
	Email sql.NullString `json:"email"`
	Pets []Pet
}

type Pet struct {
	Id int64
	Name, Breed, Animal string
}