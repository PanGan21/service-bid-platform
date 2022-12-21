package entity

type User struct {
	Id           int      `json:"Id" db:"Id"`
	Username     string   `json:"Username" db:"Username"`
	PasswordHash string   `json:"PasswordHash" db:"PasswordHash"`
	Roles        []string `json:"Roles" db:"Roles"`
}
