package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/cookiejar"

	"github.com/PanGan21/demo/data"
)

var userApi = getBasePath("user")

func login(userData data.UserData) (string, error) {
	loginPath := userApi + "/login"
	var session string

	body, err := json.Marshal(userData)
	if err != nil {
		return session, err
	}

	res, err := http.Post(loginPath, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return session, err
	}

	if res.StatusCode != 200 {
		return session, fmt.Errorf("login failed")
	}

	for _, c := range res.Cookies() {
		if c.Name == "s.id" {
			session = c.Value
		}
	}

	return session, nil
}

func logout(session string) error {
	logoutPath := userApi + "/logout"

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

	req, err := http.NewRequest("POST", logoutPath, nil)
	if err != nil {
		return err
	}

	req.AddCookie(cookie)
	res, err := client.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		return fmt.Errorf("logout failed")
	}

	return nil
}
