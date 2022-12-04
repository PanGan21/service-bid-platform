package integration_test

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"

	. "github.com/Eun/go-hit"
	"github.com/PanGan21/integration-test/testdata"
	"github.com/PanGan21/packages/auth"
	"github.com/google/uuid"
)

var userService = "user"
var requestService = "request"
var sessionId = ""
var userId = ""
var requestId = ""
var userApiPath = getBasePath(userService)
var requestApiPath = getBasePath(requestService)

func TestMain(m *testing.M) {
	err := healthCheck(Attempts, userService)
	if err != nil {
		log.Fatalf("Integration tests: host %s is not available: %s", Host, err)
	}

	log.Printf("Integration tests: Host %s is available", Host)

	sessionId, err = getSessionForMockUser()
	if err != nil || sessionId == "" {
		log.Fatalf("Integration tests: session not set for mockUser: %s", err)
	}

	code := m.Run()
	os.Exit(code)
}

// HTTP POST: /user/register
func TestHTTPDoRegister(t *testing.T) {
	sessionCookie := fmt.Sprintf(`s.id=%s`, sessionId)
	routePath := userApiPath + "/register"

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

// HTTP POST: /user/
func TestHTTPDoGetDetails(t *testing.T) {
	sessionCookie := fmt.Sprintf(`s.id=%s`, sessionId)
	routePath := userApiPath + "/"

	Test(t,
		Description("get user details; success; user exists; valid session"),
		Get(routePath),
		Send().Headers("Cookie").Add(sessionCookie),
		Expect().Status().Equal(http.StatusOK),
		Expect().Custom(func(hit Hit) error {
			var userDetails map[string]interface{}

			err := hit.Response().Body().JSON().Decode(&userDetails)
			if err != nil {
				return err
			}

			if userDetails["username"].(string) != testdata.MockUser["username"].(uuid.UUID).String() {
				return errors.New("username does not match")
			}

			if len(userDetails["roles"].([]interface{})) != len(testdata.DefaultRoles) {
				return errors.New("roles do not match")
			}

			id, ok := userDetails["id"].(string)
			if !ok {
				return errors.New("id does not exist")
			}

			userId = id

			return nil
		}),
	)

	Test(t,
		Description("get user details; failure; user exists; invalid session"),
		Get(routePath),
		Send().Headers("Cookie").Add("s.id=123"),
		Expect().Status().Equal(http.StatusUnauthorized),
		Expect().Body().String().Contains("Invalid session token"),
	)
}

// HTTP POST: /user/login
func TestHTTPDoLogin(t *testing.T) {
	routePath := userApiPath + "/login"

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

// HTTP POST: /user/authenticate
func TestHTTPDoAuthenticate(t *testing.T) {
	routePath := userApiPath + "/authenticate"
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

// HTTP GET: /request/hello
func TestHTTPDoHello(t *testing.T) {
	routePath := requestApiPath + "/hello"
	sessionCookie := fmt.Sprintf(`s.id=%s`, sessionId)

	Test(t,
		Description("request; hello; success"),
		Get(routePath),
		Send().Headers("Cookie").Add(sessionCookie),
		Expect().Status().Equal(http.StatusOK),
	)
}

// HTTP POST: /request/
func TestHTTPCreateRequest(t *testing.T) {
	routePath := requestApiPath + "/"
	sessionCookie := fmt.Sprintf(`s.id=%s`, sessionId)

	Test(t,
		Description("request; create; success"),
		Post(routePath),
		Send().Headers("Cookie").Add(sessionCookie),
		Send().Body().JSON(testdata.MockRequest),
		Expect().Status().Equal(http.StatusOK),
		Expect().Custom(func(hit Hit) error {
			id, err := hit.Response().Body().String()
			if err != nil {
				return err
			}

			if id == "" {
				return errors.New("request id is empty")
			}

			requestId = id

			return nil
		}),
	)
}

// HTTP POST: /user/logout
func TestHTTPDoLogout(t *testing.T) {
	routePath := userApiPath + "/logout"
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
