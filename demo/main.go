package main

import (
	"log"

	"github.com/PanGan21/demo/data"
)

var attempts = 10

func main() {
	resident1Session, err := login(data.Resident1User)
	if err != nil {
		log.Fatal(err)
	}

	yesterdayAuction, err := createAuction(resident1Session, data.YesterdatAuction)
	if err != nil {
		log.Fatal(err)
	}

	err = waitUntilAuctionIsAvailableInBidding(attempts, yesterdayAuction.Id)
	if err != nil {
		log.Fatal(err)
	}

	twoDaysAgoAuction, err := createAuction(resident1Session, data.TwoDaysAgoAuction)
	if err != nil {
		log.Fatal(err)
	}

	err = waitUntilAuctionIsAvailableInBidding(attempts, twoDaysAgoAuction.Id)
	if err != nil {
		log.Fatal(err)
	}

	err = logout(resident1Session)
	if err != nil {
		log.Fatal(err)
	}
}
