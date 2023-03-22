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

var biddingApi = getBasePath("bidding")

func createBid(session string, bidData data.BidData) (entity.Bid, error) {
	createBidPath := biddingApi + "/"

	var bid entity.Bid

	body, err := json.Marshal(bidData)
	if err != nil {
		return bid, err
	}

	jar, err := cookiejar.New(nil)
	if err != nil {
		return bid, err
	}

	client := http.Client{
		Jar: jar,
	}

	cookie := &http.Cookie{
		Name:  "s.id",
		Value: session,
	}

	req, err := http.NewRequest("POST", createBidPath, bytes.NewBuffer(body))
	if err != nil {
		return bid, err
	}

	req.AddCookie(cookie)
	res, err := client.Do(req)
	if err != nil {
		return bid, err
	}

	if res.StatusCode != 200 {
		return bid, fmt.Errorf("bid creation failed")
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return bid, err
	}

	err = json.Unmarshal(resBody, &bid)
	if err != nil {
		return bid, err
	}

	fmt.Println("Bid created!", bid)

	return bid, nil
}
