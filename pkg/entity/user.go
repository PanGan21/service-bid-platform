package entity

type PublicUser struct {
	Id       string `json:"Id" db:"Id"`
	Username string `json:"Username" db:"Username"`
}

type UserDetails struct {
	Email string `json:"Email" db:"Email"`
	Phone string `json:"Phone" db:"Phone"`
}

type User struct {
	Id           string   `json:"Id" db:"Id"`
	Username     string   `json:"Username" db:"Username"`
	Email        string   `json:"Email" db:"Email"`
	Phone        string   `json:"Phone" db:"Phone"`
	PasswordHash string   `json:"PasswordHash" db:"PasswordHash"`
	Roles        []string `json:"Roles" db:"Roles"`
}
