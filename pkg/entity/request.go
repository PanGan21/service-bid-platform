package entity

type Request struct {
	Id              int           `json:"Id" db:"Id"`
	Title           string        `json:"Title" db:"Title"`
	Postcode        string        `json:"Postcode" db:"Postcode"`
	Info            string        `json:"Info" db:"Info"`
	CreatorId       string        `json:"CreatorId" db:"CreatorId"`
	Deadline        int64         `json:"Deadline" db:"Deadline"`
	Status          RequestStatus `json:"Status" db:"Status"`
	RejectionReason string        `json:"RejectionReason" db:"RejectionReason"`
}

type RequestStatus string

const (
	NewRequest      RequestStatus = "new"
	ApprovedRequest RequestStatus = "approved"
	RejectedRequest RequestStatus = "rejected"
)
