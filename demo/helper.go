package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/PanGan21/pkg/entity"
	"github.com/PanGan21/pkg/postgres"
)

var localHost = "localhost"
var postgresUrl = "postgres://postgres:password@localhost:5432"

func getBasePath(service string) string {
	apiUrl, found := os.LookupEnv("API_URL")
	if !found {
		apiUrl = localHost
	}

	return fmt.Sprintf(`http://%s/%s`, apiUrl, service)
}

func getPostgresUrl() string {
	url, found := os.LookupEnv("POSTGRES_URL")
	if !found {
		url = postgresUrl
	}

	return url
}

func waitUntilAuctionIsAvailableInBidding(attempts int, auctionId int) error {
	var err error
	ctx := context.Background()

	biddingDbUrl := getPostgresUrl() + "/bidding"

	for attempts > 0 {
		var auction entity.Auction

		pg, err := postgres.New(biddingDbUrl, postgres.MaxPoolSize(2))
		if err != nil {
			fmt.Println("Error connecting with the db", err)
			return err
		}

		c, err := pg.Pool.Acquire(ctx)
		if err != nil {
			return err
		}
		defer c.Release()

		const query = `
			SELECT * FROM auctions WHERE Id=$1;
		`

		err = c.QueryRow(ctx, query, auctionId).Scan(&auction.Id, &auction.Title, &auction.Postcode, &auction.Info, &auction.CreatorId, &auction.Deadline, &auction.Status, &auction.WinningBidId, &auction.RejectionReason, &auction.WinnerId, &auction.WinningAmount)
		if err == nil && auction.Id == auctionId {
			fmt.Println("Auction available!", auction)
			return nil
		}

		log.Printf("Demo: auction with id %d is not available, attempts left: %d", auctionId, attempts)
		time.Sleep(time.Second)

		attempts--

	}

	return err
}
