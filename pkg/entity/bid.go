package entity

import "errors"

type Bid struct {
	Id        int     `json:"Id" db:"Id"`
	Amount    float64 `json:"Amount" db:"Amount"`
	CreatorId string  `json:"CreatorId" db:"CreatorId"`
	AuctionId int     `json:"AuctionId" db:"AuctionId"`
}

var ErrIncorrectBidType = errors.New("incorrect bid type")

func IsBidType(unknown interface{}) (Bid, error) {
	var bid Bid

	unknownMap, ok := unknown.(map[string]interface{})
	if !ok {
		return bid, ErrIncorrectBidType
	}

	floatId, ok := unknownMap["Id"].(float64)
	if !ok {
		return bid, ErrIncorrectBidType
	}

	bid.Id = int(floatId)

	bid.Amount, ok = unknownMap["Amount"].(float64)
	if !ok {
		return bid, ErrIncorrectBidType
	}

	bid.CreatorId, ok = unknownMap["CreatorId"].(string)
	if !ok {
		return bid, ErrIncorrectBidType
	}

	floatAuctionId, ok := unknownMap["AuctionId"].(float64)
	if !ok {
		return bid, ErrIncorrectBidType
	}

	bid.AuctionId = int(floatAuctionId)

	return bid, nil
}
