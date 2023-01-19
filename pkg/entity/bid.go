package entity

type Bid struct {
	Id        int     `json:"Id" db:"Id"`
	Amount    float64 `json:"Amount" db:"Amount"`
	CreatorId string  `json:"CreatorId" db:"CreatorId"`
	RequestId int     `json:"RequestId" db:"RequestId"`
}
