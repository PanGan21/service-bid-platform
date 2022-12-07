package entity

import "github.com/google/uuid"

type User struct {
	Id           uuid.UUID `json:"id" db:"id"`
	Username     string    `json:"username" db:"username"`
	PasswordHash string    `json:"passwordHash" db:"passwordHash"`
	Roles        []string  `json:"roles" db:"roles"`
}
