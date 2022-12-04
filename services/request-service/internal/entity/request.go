package entity

import "github.com/google/uuid"

type Request struct {
	Id        uuid.UUID `json:"id" db:"id"`
	Title     string    `json:"title" db:"title"`
	Postcode  string    `json:"postcode" db:"postcode"`
	Info      string    `json:"info" db:"info"`
	CreatorId string    `json:"creatorId" db:"creatorId"`
	Deadline  int64     `json:"deadline" db:"deadline"`
}
