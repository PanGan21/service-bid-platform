package integration_test

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	. "github.com/Eun/go-hit"
	"github.com/google/uuid"
)

var (
	// Attempts connection
	host       = getHost()
	healthPath = "http://" + host + "/healthz"
	attempts   = 20

	// HTTP REST
	basePath = "http://" + host
)

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

func TestMain(m *testing.M) {
	apiHost, found := os.LookupEnv("API_HOST")
	if found {
		host = apiHost
		log.Println("Starting integration tests in Docker")
	} else {
		log.Println("Starting integration tests locally")
	}

	err := healethCheck(attempts)
	if err != nil {
		log.Fatalf("Integration tests: host %s is not available: %s", host, err)
	}

	log.Printf("Integration tests: host %s is available", host)

	code := m.Run()
	os.Exit(code)
}

func healethCheck(attempts int) error {
	var err error

	for attempts > 0 {
		err = Do(Get(healthPath), Expect().Status().Equal(http.StatusOK))
		if err == nil {
			return nil
		}

		log.Printf("Integration tests: url %s is not available, attempts left: %d", healthPath, attempts)
		time.Sleep(time.Second)

		attempts--

	}

	return err
}

// HTTP POST: /register
func TestHTTPDoRegister(t *testing.T) {
	var mockUser = map[string]interface{}{"username": uuid.New(), "password": uuid.New()}

	var sessionId string = ""

	Test(t,
		Description("register; success"),
		Post(basePath+"/register"),
		Send().Headers("Content-Type").Add("application/json"),
		Send().Body().JSON(mockUser),
		Expect().Status().Equal(http.StatusOK),
		Expect().Body().String().Contains("Successfully registered user"),
		Expect().Custom(func(hit Hit) error {
			var cookies = hit.Response().Cookies()
			for _, c := range cookies {
				if c.Name == "s.id" {
					sessionId = c.Value
				}
				if sessionId == "" {
					return errors.New("Session is missing")
				}
			}
			return nil
		}),
	)

	Test(t,
		Description("register; success; already logged in"),
		Post(basePath+"/register"),
		Send().Headers("Content-Type").Add("application/json"),
		Send().Headers("Cookie").Add(fmt.Sprintf(`s.id=%s`, sessionId)),
		Send().Body().JSON(mockUser),
		Expect().Status().Equal(http.StatusOK),
		Expect().Body().String().Contains("Already logged in"),
	)

	Test(t,
		Description("register; failure; registration failed"),
		Post(basePath+"/register"),
		Send().Headers("Content-Type").Add("application/json"),
		Send().Body().JSON(mockUser),
		Expect().Status().Equal(http.StatusInternalServerError),
		Expect().Body().String().Contains("Registration failed"),
	)
}
