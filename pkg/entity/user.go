package entity

type User struct {
	Id           int      `json:"id" db:"id"`
	Username     string   `json:"username" db:"username"`
	PasswordHash string   `json:"passwordHash" db:"passwordHash"`
	Roles        []string `json:"roles" db:"roles"`
}
