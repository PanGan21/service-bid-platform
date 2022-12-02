package entity

import "github.com/google/uuid"

type Request struct {
	Id uuid.UUID `json:"id" db:"id"`
}
