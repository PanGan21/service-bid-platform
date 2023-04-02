package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"

	"github.com/PanGan21/pkg/entity"
)

var auctionApi = getBasePath("auction")

func updateAuctionWinner(session string, auctionId string) error {
	updateAuctionWinnerPath := auctionApi + "/update/winner?auctionId=" + auctionId

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

	req, err := http.NewRequest("POST", updateAuctionWinnerPath, nil)
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
