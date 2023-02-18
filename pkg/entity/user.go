package entity

type PublicUser struct {
	Id       string `json:"Id" db:"Id"`
	Username string `json:"Username" db:"Username"`
}

type User struct {
	PublicUser
	PasswordHash string   `json:"PasswordHash" db:"PasswordHash"`
	Roles        []string `json:"Roles" db:"Roles"`
}
