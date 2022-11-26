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
	"github.com/PanGan21/integration-test/testdata"
	"github.com/PanGan21/packages/auth"
)

var (
	// Attempts connection
	host       = getHost()
	healthPath = "http://" + host + "/healthz"
	attempts   = 20

	// HTTP REST
	basePath  = "http://" + host
	sessionId = ""
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

	err = setSessionForMockUser()
	if err != nil {
		log.Fatalf("Integration tests: session not set for mockUser: %s", err)
	}

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

func setSessionForMockUser() error {
	routePath := basePath + "/register"

	Do(Post(routePath),
		Send().Headers("Content-Type").Add("application/json"),
		Send().Body().JSON(testdata.MockUser),
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

	return nil
}

// HTTP POST: /register
func TestHTTPDoRegister(t *testing.T) {
	sessionCookie := fmt.Sprintf(`s.id=%s`, sessionId)
	routePath := basePath + "/register"

	Test(t,
		Description("register; success; user exists; valid session"),
		Post(routePath),
		Send().Headers("Content-Type").Add("application/json"),
		Send().Headers("Cookie").Add(sessionCookie),
		Send().Body().JSON(testdata.MockUser),
		Expect().Status().Equal(http.StatusOK),
		Expect().Body().String().Contains("Already logged in"),
	)

	Test(t,
		Description("register; failure; user exists; invalid session"),
		Post(routePath),
		Send().Headers("Content-Type").Add("application/json"),
		Send().Body().JSON(testdata.MockUser),
		Expect().Status().Equal(http.StatusInternalServerError),
		Expect().Body().String().Contains("Registration failed"),
	)

	Test(t,
		Description("register; failure; validation error"),
		Post(routePath),
		Send().Headers("Content-Type").Add("application/json"),
		Send().Body().JSON(map[string]interface{}{"username": 123}),
		Expect().Status().Equal(http.StatusBadRequest),
		Expect().Body().String().Contains("Validation error"),
	)
}

// HTTP POST: /login
func TestHTTPDoLogin(t *testing.T) {
	routePath := basePath + "/login"

	Test(t,
		Description("login; success"),
		Post(routePath),
		Send().Headers("Content-Type").Add("application/json"),
		Send().Body().JSON(testdata.MockUser),
		Expect().Status().Equal(http.StatusOK),
		Expect().Body().String().Contains("Successfully authenticated user"),
		Expect().Custom(func(hit Hit) error {
			var cookies = hit.Response().Cookies()

			var loginSessionId = ""
			for _, c := range cookies {
				if c.Name == "s.id" {
					loginSessionId = c.Value
				}
				if loginSessionId == "" {
					return errors.New("Session is missing")
				}
			}
			return nil
		}),
	)

	Test(t,
		Description("login; failure; validation error"),
		Post(routePath),
		Send().Headers("Content-Type").Add("application/json"),
		Send().Body().JSON(map[string]interface{}{"username": "RANDOM", "password": "RANDOM"}),
		Expect().Status().Equal(http.StatusUnauthorized),
		Expect().Body().String().Contains("Authentication failed"),
	)

	Test(t,
		Description("login; failure; validation error"),
		Post(routePath),
		Send().Headers("Content-Type").Add("application/json"),
		Send().Body().JSON(map[string]interface{}{"username": 123}),
		Expect().Status().Equal(http.StatusBadRequest),
		Expect().Body().String().Contains("Validation error"),
	)
}

// HTTP POST: /authenticate
func TestHTTPDoAuthenticate(t *testing.T) {
	routePath := basePath + "/authenticate"
	sessionCookie := fmt.Sprintf(`s.id=%s`, sessionId)

	Test(t,
		Description("authenticate; success"),
		Get(routePath),
		Send().Headers("Cookie").Add(sessionCookie),
		Expect().Status().Equal(http.StatusOK),
		Expect().Body().String().Contains("Successfully authenticated user"),
		Expect().Custom(func(hit Hit) error {
			jwtHeader := hit.Response().Header.Get("X-Internal-Jwt")
			if jwtHeader == "" {
				return errors.New("No jwt in X-Internal-Jwt header")
			}

			var secret = "auth_secret"
			authService := auth.NewAuthService([]byte(secret))
			_, err := authService.VerifyJWT(jwtHeader, "/authenticate")
			if err != nil {
				return err
			}

			return nil
		}),
	)
}

// HTTP POST: /logout
func TestHTTPDoLogout(t *testing.T) {
	routePath := basePath + "/logout"
	sessionCookie := fmt.Sprintf(`s.id=%s`, sessionId)

	Test(t,
		Description("logout; success"),
		Post(routePath),
		Send().Headers("Content-Type").Add("application/json"),
		Send().Headers("Cookie").Add(sessionCookie),
		Send().Body().JSON(testdata.MockUser),
		Expect().Status().Equal(http.StatusOK),
		Expect().Body().String().Contains("Successfully logged out"),
	)

	Test(t,
		Description("logout; failure; invalid session"),
		Post(routePath),
		Send().Headers("Content-Type").Add("application/json"),
		Send().Headers("Cookie").Add("s.id=123"),
		Send().Body().JSON(testdata.MockUser),
		Expect().Status().Equal(http.StatusBadRequest),
		Expect().Body().String().Contains("Invalid session token"),
	)
}
