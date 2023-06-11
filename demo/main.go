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

	_, err = createRequest(resident1Session, data.YesterdayRequestNew)
	if err != nil {
		log.Fatal(err)
	}

	_, err = createRequest(resident1Session, data.TwoDaysAgoRequestNew)
	if err != nil {
		log.Fatal(err)
	}

	twoDaysAgoRequestOpen, err := createRequest(resident2Session, data.TwoDaysAgoRequestOpen)
	if err != nil {
		log.Fatal(err)
	}

	twoDaysAgoRequestAssigned, err := createRequest(resident2Session, data.TwoDaysAgoRequestAssigned)
	if err != nil {
		log.Fatal(err)
	}

	err = approveRequest(superAdminSession, fmt.Sprint(twoDaysAgoRequestOpen.Id), 0)
	if err != nil {
		log.Fatal(err)
	}

	err = approveRequest(superAdminSession, fmt.Sprint(twoDaysAgoRequestAssigned.Id), 0)
	if err != nil {
		log.Fatal(err)
	}

	err = waitUntilAuctionIsAvailableInAuction(attempts, twoDaysAgoRequestOpen.Id)
	if err != nil {
		log.Fatal(err)
	}

	err = waitUntilAuctionIsAvailableInBidding(attempts, twoDaysAgoRequestOpen.Id)
	if err != nil {
		log.Fatal(err)
	}

	err = waitUntilAuctionIsAvailableInAuction(attempts, twoDaysAgoRequestAssigned.Id)
	if err != nil {
		log.Fatal(err)
	}

	err = waitUntilAuctionIsAvailableInBidding(attempts, twoDaysAgoRequestAssigned.Id)
	if err != nil {
		log.Fatal(err)
	}

	bid100 := data.BidData{
		Amount:    100.0,
		AuctionId: twoDaysAgoRequestAssigned.Id,
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
		AuctionId: twoDaysAgoRequestAssigned.Id,
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
		AuctionId: twoDaysAgoRequestAssigned.Id,
	}
	bid_300, err := createBid(bidder3Session, bid300)
	if err != nil {
		log.Fatal(err)
	}

	err = waitUntilBidIsAvailableInAuction(attempts, bid_300.Id)
	if err != nil {
		log.Fatal(err)
	}

	err = updateAuctionWinner(superAdminSession, fmt.Sprint(twoDaysAgoRequestAssigned.Id))
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
