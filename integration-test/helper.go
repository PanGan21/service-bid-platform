package integration_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/PanGan21/integration-test/testdata"
	"github.com/PanGan21/pkg/entity"
	"github.com/PanGan21/pkg/postgres"
)

var (
	// Attempts connection
	Host     = getHost()
	Attempts = 20

	// HTTP REST
	BasePath = "http://" + Host
)

func healthCheck(attempts int, service string) error {
	var servicePath = ""
	_, found := os.LookupEnv("API_HOST")
	if found {
		// in docker
		servicePath = "/" + service
	}

	var healthPath = "http://" + Host + servicePath + "/healthz"
	var err error

	for attempts > 0 {
		res, err := http.Get(healthPath)
		if err == nil && res.StatusCode == 200 {
			return nil
		}

		log.Printf("Integration tests: url %s is not available, attempts left: %d", healthPath, attempts)
		time.Sleep(time.Second)

		attempts--

	}

	return err
}

func getHost() string {
	var localHost = "localhost"
	var localPort = "8000"

	var host string
	var port string

	apiHost, found := os.LookupEnv("API_HOST")
	if !found {
		host = localHost
	} else {
		host = apiHost
	}

	apiPort, found := os.LookupEnv("API_PORT")
	if !found {
		port = localPort
	} else {
		port = apiPort
	}

	return fmt.Sprintf(`%s:%s`, host, port)
}

func getSessionForMockUser() (string, error) {
	var sessionId = ""

	var servicePath = ""
	_, found := os.LookupEnv("API_HOST")
	if found {
		// in docker
		servicePath = "/" + "user"
	}

	routePath := BasePath + servicePath + "/register"

	jsonBody, err := json.Marshal(testdata.MockUser)
	if err != nil {
		return "", err
	}
	res, err := http.Post(routePath, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", err
	}

	for _, c := range res.Cookies() {
		if c.Name == "s.id" {
			sessionId = c.Value
		}
		if sessionId == "" {
			return "", errors.New("session is missing")
		}
	}

	return sessionId, nil
}

func getBasePath(service string) string {
	var localHost = "localhost"
	var localPort = "8000"

	var host string
	var port string
	var pathService = ""

	apiHost, found := os.LookupEnv("API_HOST")
	if !found {
		host = localHost
	} else {
		host = apiHost
		pathService = "/" + service
	}

	apiPort, found := os.LookupEnv("API_PORT")
	if !found {
		port = localPort
	} else {
		port = apiPort
	}

	return fmt.Sprintf(`http://%s:%s%s`, host, port, pathService)
}

func waitUntilAuctionIsAvailableInBidding(attempts int, auctionId int) error {
	var err error
	ctx := context.Background()

	for attempts > 0 {
		var auction entity.Auction

		pg, err := postgres.New("postgres://postgres:password@postgres:5432/bidding", postgres.MaxPoolSize(2))
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

		err = c.QueryRow(ctx, query, auctionId).Scan(&auction.Id, &auction.Title, &auction.Postcode, &auction.Info, &auction.CreatorId, &auction.Deadline, &auction.Status, &auction.WinningBidId)
		if err == nil && auction.Id == auctionId {
			fmt.Println("Auction available!", auction)
			return nil
		}

		log.Printf("Integration tests: auction with id %d is not available, attempts left: %d", auctionId, attempts)
		time.Sleep(time.Second)

		attempts--

	}

	return err
}

func waitUntilBidIsAvailableInAuction(attempts int, bidId int) error {
	var err error
	ctx := context.Background()

	for attempts > 0 {
		var bid entity.Bid

		pg, err := postgres.New("postgres://postgres:password@postgres:5432/auction", postgres.MaxPoolSize(2))
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
			SELECT * FROM bids WHERE Id=$1;
		`

		err = c.QueryRow(ctx, query, bidId).Scan(&bid.Id, &bid.Amount, &bid.CreatorId, &bid.AuctionId)
		if err == nil && bid.Id == bidId {
			fmt.Println("Bid available!", bid)
			return nil
		}

		log.Printf("Integration tests: bid with id %d is not available, attempts left: %d", bidId, attempts)
		time.Sleep(time.Second)

		attempts--

	}

	return err
}
