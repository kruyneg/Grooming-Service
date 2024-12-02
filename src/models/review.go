package models

type Review struct {
	Host UserData
	SalonMaster SalonMaster
	Score int
	Content string
}