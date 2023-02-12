package entity

import (
	"errors"
)

type Auction struct {
	Id           int           `json:"Id" db:"Id"`
	Title        string        `json:"Title" db:"Title"`
	Postcode     string        `json:"Postcode" db:"Postcode"`
	Info         string        `json:"Info" db:"Info"`
	CreatorId    string        `json:"CreatorId" db:"CreatorId"`
	Deadline     int64         `json:"Deadline" db:"Deadline"`
	Status       AuctionStatus `json:"Status" db:"Status"`
	WinningBidId string        `json:"WinningBidId" db:"WinningBidId"`
}

type ExtendedAuction struct {
	Id           int           `json:"Id" db:"Id"`
	Title        string        `json:"Title" db:"Title"`
	Postcode     string        `json:"Postcode" db:"Postcode"`
	Info         string        `json:"Info" db:"Info"`
	CreatorId    string        `json:"CreatorId" db:"CreatorId"`
	Deadline     int64         `json:"Deadline" db:"Deadline"`
	Status       AuctionStatus `json:"Status" db:"Status"`
	WinningBidId string        `json:"WinningBidId" db:"WinningBidId"`
	BidsCount    int           `json:"BidsCount" db:"BidsCount"`
}

type BidPopulatedAuction struct {
	Id        int           `json:"Id" db:"Id"`
	Title     string        `json:"Title" db:"Title"`
	Postcode  string        `json:"Postcode" db:"Postcode"`
	Info      string        `json:"Info" db:"Info"`
	CreatorId string        `json:"CreatorId" db:"CreatorId"`
	Deadline  int64         `json:"Deadline" db:"Deadline"`
	Status    AuctionStatus `json:"Status" db:"Status"`
	BidId     int           `json:"BidId" db:"BidId"`
	BidAmount float64       `json:"BidAmount" db:"BidAmount"`
}

type AuctionStatus string

const (
	New        AuctionStatus = "new"
	Rejected   AuctionStatus = "rejected"
	Open       AuctionStatus = "open"
	Assigned   AuctionStatus = "assigned"
	InProgress AuctionStatus = "in progress"
	Closed     AuctionStatus = "closed"
)

var ErrIncorrectAuctionType = errors.New("incorrect auction type")

func IsAuctionType(unknown interface{}) (Auction, error) {
	var auction Auction

	unknownMap, ok := unknown.(map[string]interface{})
	if !ok {
		return auction, ErrIncorrectAuctionType
	}

	auction.CreatorId, ok = unknownMap["CreatorId"].(string)
	if !ok {
		return auction, ErrIncorrectAuctionType
	}

	auction.Info, ok = unknownMap["Info"].(string)
	if !ok {
		return auction, ErrIncorrectAuctionType
	}

	auction.Postcode, ok = unknownMap["Postcode"].(string)
	if !ok {
		return auction, ErrIncorrectAuctionType
	}

	auction.Title, ok = unknownMap["Title"].(string)
	if !ok {
		return auction, ErrIncorrectAuctionType
	}

	floatId, ok := unknownMap["Id"].(float64)
	if !ok {
		return auction, ErrIncorrectAuctionType
	}
	auction.Id = int(floatId)

	floatDeadline, ok := unknownMap["Deadline"].(float64)
	if !ok {
		return auction, ErrIncorrectAuctionType
	}
	auction.Deadline = int64(floatDeadline)

	s, ok := unknownMap["Status"].(string)
	if !ok {
		return auction, ErrIncorrectAuctionType
	}
	status := AuctionStatus(s)

	switch status {
	case Open, New, Rejected, Assigned, InProgress, Closed:
		auction.Status = status
	default:
		return auction, ErrIncorrectAuctionType
	}

	auction.WinningBidId = unknownMap["WinningBidId"].(string)
	if !ok {
		return auction, ErrIncorrectAuctionType
	}

	return auction, nil
}
