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

var requestApi = getBasePath("request")

func createRequest(session string, requestData data.RequestData) (entity.Request, error) {
	createRequestPath := requestApi + "/"

	var request entity.Request

	body, err := json.Marshal(requestData)
	if err != nil {
		return request, err
	}

	jar, err := cookiejar.New(nil)
	if err != nil {
		return request, err
	}

	client := http.Client{
		Jar: jar,
	}

	cookie := &http.Cookie{
		Name:  "s.id",
		Value: session,
	}

	req, err := http.NewRequest("POST", createRequestPath, bytes.NewBuffer(body))
	if err != nil {
		return request, err
	}

	req.AddCookie(cookie)
	res, err := client.Do(req)
	if err != nil {
		return request, err
	}

	if res.StatusCode != 200 {
		return request, fmt.Errorf("request creation failed")
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return request, err
	}

	err = json.Unmarshal(resBody, &request)
	if err != nil {
		return request, err
	}

	fmt.Println("Request created!", request)

	return request, nil
}

func approveRequest(session string, requestId string) error {
	updateRequestPath := requestApi + "/approve?requestId=" + requestId

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

	req, err := http.NewRequest("POST", updateRequestPath, nil)
	if err != nil {
		return err
	}

	req.AddCookie(cookie)
	res, err := client.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		return fmt.Errorf("request update failed")
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	var request entity.Request

	err = json.Unmarshal(resBody, &request)
	if err != nil {
		return err
	}

	fmt.Println("Request updated!", request)

	return nil
}
