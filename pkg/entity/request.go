package entity

import "errors"

type Request struct {
	Id              int           `json:"Id" db:"Id"`
	Title           string        `json:"Title" db:"Title"`
	Postcode        string        `json:"Postcode" db:"Postcode"`
	Info            string        `json:"Info" db:"Info"`
	CreatorId       string        `json:"CreatorId" db:"CreatorId"`
	Status          RequestStatus `json:"Status" db:"Status"`
	RejectionReason string        `json:"RejectionReason" db:"RejectionReason"`
}

type RequestStatus string

const (
	NewRequest      RequestStatus = "new"
	ApprovedRequest RequestStatus = "approved"
	RejectedRequest RequestStatus = "rejected"
)

var ErrIncorrectRequestType = errors.New("incorrect request type")

func IsRequestType(unknown interface{}) (Request, error) {
	var request Request

	unknownMap, ok := unknown.(map[string]interface{})
	if !ok {
		return request, ErrIncorrectRequestType
	}

	request.CreatorId, ok = unknownMap["CreatorId"].(string)
	if !ok {
		return request, ErrIncorrectRequestType
	}

	request.Info, ok = unknownMap["Info"].(string)
	if !ok {
		return request, ErrIncorrectRequestType
	}

	request.Postcode, ok = unknownMap["Postcode"].(string)
	if !ok {
		return request, ErrIncorrectRequestType
	}

	request.Title, ok = unknownMap["Title"].(string)
	if !ok {
		return request, ErrIncorrectRequestType
	}

	floatId, ok := unknownMap["Id"].(float64)
	if !ok {
		return request, ErrIncorrectRequestType
	}
	request.Id = int(floatId)

	s, ok := unknownMap["Status"].(string)
	if !ok {
		return request, ErrIncorrectRequestType
	}
	status := RequestStatus(s)

	switch status {
	case NewRequest, ApprovedRequest, RejectedRequest:
		request.Status = status
	default:
		return request, ErrIncorrectRequestType
	}

	request.RejectionReason = unknownMap["RejectionReason"].(string)
	if !ok {
		return request, ErrIncorrectRequestType
	}

	return request, nil
}
