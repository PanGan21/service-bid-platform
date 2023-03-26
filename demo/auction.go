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

	fmt.Println("Auction created!", auction)

	return auction, nil
}

func updateAuctionStatusToOpen(session string, auctionId string) error {
	updateAuctionPath := auctionApi + "/update/status?auctionId=" + auctionId

	body, err := json.Marshal(map[string]string{"Status": string(entity.Open)})
	if err != nil {
		return err
	}

	jar, err := cookiejar.New(nil)
	if err != nil {
		return err
	}

	client := http.Client{
		Jar: jar,
	}

	cookie := &http.Cookie{
		Name:  "s.id",
		Value: session,
	}

	req, err := http.NewRequest("POST", updateAuctionPath, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	req.AddCookie(cookie)
	res, err := client.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		return fmt.Errorf("auction update failed")
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	var auction entity.Auction

	err = json.Unmarshal(resBody, &auction)
	if err != nil {
		return err
	}

	fmt.Println("Auction updated!", auction)

	return nil
}

func updateAuctionWinner(session string, auctionId string) error {
	updateAuctionWinnerPath := auctionApi + "/update/winner?auctionId=" + auctionId

	body, err := json.Marshal(map[string]string{})
	if err != nil {
		return err
	}

	jar, err := cookiejar.New(nil)
	if err != nil {
		return err
	}

	client := http.Client{
		Jar: jar,
	}

	cookie := &http.Cookie{
		Name:  "s.id",
		Value: session,
	}

	req, err := http.NewRequest("POST", updateAuctionWinnerPath, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	req.AddCookie(cookie)
	res, err := client.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		return fmt.Errorf("auction winner update failed")
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	var auction entity.Auction

	err = json.Unmarshal(resBody, &auction)
	if err != nil {
		return err
	}

	fmt.Println("Auction winner updated!", auction)

	return nil
}
