package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"

	"github.com/PanGan21/demo/data"
	"github.com/PanGan21/pkg/entity"
)

var auctionApi = getBasePath("auction")

func createAuction(session string, auctionData data.AuctionData) (entity.Auction, error) {
	createAuctionPath := auctionApi + "/"

	var auction entity.Auction

	body, err := json.Marshal(auctionData)
	if err != nil {
		return auction, err
	}

	jar, err := cookiejar.New(nil)
	if err != nil {
		return auction, err
	}

	client := http.Client{
		Jar: jar,
	}

	cookie := &http.Cookie{
		Name:  "s.id",
		Value: session,
	}

	req, err := http.NewRequest("POST", createAuctionPath, bytes.NewBuffer(body))
	if err != nil {
		return auction, err
	}

	req.AddCookie(cookie)
	res, err := client.Do(req)
	if err != nil {
		return auction, err
	}

	if res.StatusCode != 200 {
		return auction, fmt.Errorf("auction creation failed")
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return auction, err
	}

	err = json.Unmarshal(resBody, &auction)
	if err != nil {
		return auction, err
	}

	return auction, nil
}
