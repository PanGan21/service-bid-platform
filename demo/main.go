package main

import (
	"fmt"
	"log"

	"github.com/PanGan21/demo/data"
)

var attempts = 10

func main() {
	resident1Session, err := login(data.Resident1User)
	if err != nil {
		log.Fatal(err)
	}

	resident2Session, err := login(data.Resident2User)
	if err != nil {
		log.Fatal(err)
	}

	bidder1Session, err := login(data.Bidder1User)
	if err != nil {
		log.Fatal(err)
	}

	bidder2Session, err := login(data.Bidder2User)
	if err != nil {
		log.Fatal(err)
	}

	bidder3Session, err := login(data.Bidder3User)
	if err != nil {
		log.Fatal(err)
	}

	superAdminSession, err := login(data.SuperAdmin)
	if err != nil {
		log.Fatal(err)
	}

	yesterdayAuctionNew, err := createAuction(resident1Session, data.YesterdatAuctionNew)
	if err != nil {
		log.Fatal(err)
	}

	err = waitUntilAuctionIsAvailableInBidding(attempts, yesterdayAuctionNew.Id)
	if err != nil {
		log.Fatal(err)
	}

	twoDaysAgoAuctionNew, err := createAuction(resident1Session, data.TwoDaysAgoAuctionNew)
	if err != nil {
		log.Fatal(err)
	}

	err = waitUntilAuctionIsAvailableInBidding(attempts, twoDaysAgoAuctionNew.Id)
	if err != nil {
		log.Fatal(err)
	}

	twoDaysAgoAuctionOpen, err := createAuction(resident2Session, data.TwoDaysAgoAuctionOpen)
	if err != nil {
		log.Fatal(err)
	}

	err = waitUntilAuctionIsAvailableInBidding(attempts, twoDaysAgoAuctionOpen.Id)
	if err != nil {
		log.Fatal(err)
	}

	twoDaysAgoAuctionAssigned, err := createAuction(resident2Session, data.TwoDaysAgoAuctionAssigned)
	if err != nil {
		log.Fatal(err)
	}

	err = waitUntilAuctionIsAvailableInBidding(attempts, twoDaysAgoAuctionAssigned.Id)
	if err != nil {
		log.Fatal(err)
	}

	err = updateAuctionStatusToOpen(superAdminSession, fmt.Sprint(twoDaysAgoAuctionOpen.Id))
	if err != nil {
		log.Fatal(err)
	}

	err = updateAuctionStatusToOpen(superAdminSession, fmt.Sprint(twoDaysAgoAuctionAssigned.Id))
	if err != nil {
		log.Fatal(err)
	}

	err = waitUntilAuctionIsOpenToBids(attempts, twoDaysAgoAuctionAssigned.Id)
	if err != nil {
		log.Fatal(err)
	}

	bid100 := data.BidData{
		Amount:    100.0,
		AuctionId: twoDaysAgoAuctionAssigned.Id,
	}
	bid_100, err := createBid(bidder1Session, bid100)
	if err != nil {
		log.Fatal(err)
	}

	err = waitUntilBidIsAvailableInAuction(attempts, bid_100.Id)
	if err != nil {
		log.Fatal(err)
	}

	bid200 := data.BidData{
		Amount:    200.0,
		AuctionId: twoDaysAgoAuctionAssigned.Id,
	}
	bid_200, err := createBid(bidder2Session, bid200)
	if err != nil {
		log.Fatal(err)
	}

	err = waitUntilBidIsAvailableInAuction(attempts, bid_200.Id)
	if err != nil {
		log.Fatal(err)
	}

	bid300 := data.BidData{
		Amount:    300.0,
		AuctionId: twoDaysAgoAuctionAssigned.Id,
	}
	bid_300, err := createBid(bidder3Session, bid300)
	if err != nil {
		log.Fatal(err)
	}

	err = waitUntilBidIsAvailableInAuction(attempts, bid_300.Id)
	if err != nil {
		log.Fatal(err)
	}

	err = updateAuctionWinner(superAdminSession, fmt.Sprint(twoDaysAgoAuctionAssigned.Id))
	if err != nil {
		log.Fatal(err)
	}

	err = logout(resident1Session)
	if err != nil {
		log.Fatal(err)
	}

	err = logout(resident2Session)
	if err != nil {
		log.Fatal(err)
	}

	err = logout(bidder1Session)
	if err != nil {
		log.Fatal(err)
	}

	err = logout(superAdminSession)
	if err != nil {
		log.Fatal(err)
	}
}
