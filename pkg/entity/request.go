package entity

import (
	"errors"
)

type Request struct {
	Id           int           `json:"Id" db:"Id"`
	Title        string        `json:"Title" db:"Title"`
	Postcode     string        `json:"Postcode" db:"Postcode"`
	Info         string        `json:"Info" db:"Info"`
	CreatorId    string        `json:"CreatorId" db:"CreatorId"`
	Deadline     int64         `json:"Deadline" db:"Deadline"`
	Status       RequestStatus `json:"Status" db:"Status"`
	WinningBidId string        `json:"WinningBidId" db:"WinningBidId"`
}

type ExtendedRequest struct {
	Id           int           `json:"Id" db:"Id"`
	Title        string        `json:"Title" db:"Title"`
	Postcode     string        `json:"Postcode" db:"Postcode"`
	Info         string        `json:"Info" db:"Info"`
	CreatorId    string        `json:"CreatorId" db:"CreatorId"`
	Deadline     int64         `json:"Deadline" db:"Deadline"`
	Status       RequestStatus `json:"Status" db:"Status"`
	WinningBidId string        `json:"WinningBidId" db:"WinningBidId"`
	BidsCount    int           `json:"BidsCount" db:"BidsCount"`
}

type BidPopulatedRequest struct {
	Id        int           `json:"Id" db:"Id"`
	Title     string        `json:"Title" db:"Title"`
	Postcode  string        `json:"Postcode" db:"Postcode"`
	Info      string        `json:"Info" db:"Info"`
	CreatorId string        `json:"CreatorId" db:"CreatorId"`
	Deadline  int64         `json:"Deadline" db:"Deadline"`
	Status    RequestStatus `json:"Status" db:"Status"`
	BidId     int           `json:"BidId" db:"BidId"`
	BidAmount float64       `json:"BidAmount" db:"BidAmount"`
}

type RequestStatus string

const (
	Open       RequestStatus = "open"
	Assigned   RequestStatus = "assigned"
	InProgress RequestStatus = "in progress"
	Closed     RequestStatus = "closed"
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

	floatDeadline, ok := unknownMap["Deadline"].(float64)
	if !ok {
		return request, ErrIncorrectRequestType
	}
	request.Deadline = int64(floatDeadline)

	s, ok := unknownMap["Status"].(string)
	if !ok {
		return request, ErrIncorrectRequestType
	}
	status := RequestStatus(s)

	switch status {
	case Open, Assigned, InProgress, Closed:
		request.Status = status
	default:
		return request, ErrIncorrectRequestType
	}

	request.WinningBidId = unknownMap["WinningBidId"].(string)
	if !ok {
		return request, ErrIncorrectRequestType
	}

	return request, nil
}
