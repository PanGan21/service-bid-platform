package integration_test

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"testing"

	. "github.com/Eun/go-hit"
	"github.com/PanGan21/integration-test/testdata"
	"github.com/PanGan21/pkg/auth"
	"github.com/PanGan21/pkg/entity"
)

var userService = "user"
var requestService = "request"
var biddingService = "bidding"

var sessionId = ""
var userId = ""
var requestId = 0

var userApiPath = getBasePath(userService)
var requestApiPath = getBasePath(requestService)
var biddingApiPath = getBasePath(biddingService)

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

// HTTP GET: /user/
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

			if userDetails["Username"].(string) != testdata.MockUser["Username"].(string) {
				return errors.New("username does not match")
			}

			if len(userDetails["Roles"].([]interface{})) != len(testdata.DefaultRoles) {
				return errors.New("roles do not match")
			}

			id, ok := userDetails["Id"].(string)
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
			var request entity.Request

			err := hit.Response().Body().JSON().Decode(&request)
			if err != nil {
				return err
			}

			requestId = request.Id

			return nil
		}),
	)
}

// HTTP GET: /request/
func TestHTTPGetPaginatedRequests(t *testing.T) {
	createRoutePath := requestApiPath + "/"
	sessionCookie := fmt.Sprintf(`s.id=%s`, sessionId)

	for i := 2; i <= 10; i++ {
		description := fmt.Sprintf("request; create; success; no %d", i)
		testRequest := testdata.MockRequest
		testRequest["deadline"] = i

		Test(t,
			Description(description),
			Post(createRoutePath),
			Send().Headers("Cookie").Add(sessionCookie),
			Send().Body().JSON(testRequest),
			Expect().Status().Equal(http.StatusOK),
		)

	}

	limit10 := 10
	page1 := 1

	routePathAscendingOrder := fmt.Sprintf("%s?limit=%d&page=%d&asc=true", createRoutePath, limit10, page1)

	Test(t,
		Description("get requests; success; ascending order"),
		Get(routePathAscendingOrder),
		Send().Headers("Cookie").Add(sessionCookie),
		Expect().Status().Equal(http.StatusOK),
		Expect().Custom(func(hit Hit) error {
			var requests []entity.Request
			err := hit.Response().Body().JSON().Decode(&requests)
			if err != nil {
				return err
			}

			if len(requests) != limit10 {
				return fmt.Errorf("requests should be %d", limit10)
			}

			isAscendingOrder := sort.SliceIsSorted(requests, func(p, q int) bool {
				return requests[p].Deadline < requests[q].Deadline
			})

			if !isAscendingOrder {
				return errors.New("requests are not in ascending order")
			}

			return nil
		}),
	)

	routePathDescendingOrder := fmt.Sprintf("%s?limit=%d&page=%d&asc=false", createRoutePath, limit10, page1)

	Test(t,
		Description("get requests; success; descending order"),
		Get(routePathDescendingOrder),
		Send().Headers("Cookie").Add(sessionCookie),
		Expect().Status().Equal(http.StatusOK),
		Expect().Custom(func(hit Hit) error {
			var requests []entity.Request
			err := hit.Response().Body().JSON().Decode(&requests)
			if err != nil {
				return err
			}

			if len(requests) != limit10 {
				return fmt.Errorf("requests should be %d", limit10)
			}

			isAscendingOrder := sort.SliceIsSorted(requests, func(p, q int) bool {
				return requests[p].Deadline < requests[q].Deadline
			})

			if isAscendingOrder {
				return errors.New("requests are not in descending order")
			}

			return nil
		}),
	)
}

// // HTTP GET: /request/own
func TestHTTPGetPaginatedOwnRequests(t *testing.T) {
	createRoutePath := requestApiPath + "/own"
	sessionCookie := fmt.Sprintf(`s.id=%s`, sessionId)

	limit10 := 10
	page1 := 1

	routePathAscendingOrder := fmt.Sprintf("%s?limit=%d&page=%d&asc=true", createRoutePath, limit10, page1)

	Test(t,
		Description("get owned requests; success; ascending order"),
		Get(routePathAscendingOrder),
		Send().Headers("Cookie").Add(sessionCookie),
		Expect().Status().Equal(http.StatusOK),
		Expect().Custom(func(hit Hit) error {
			var requests []entity.Request
			err := hit.Response().Body().JSON().Decode(&requests)
			if err != nil {
				return err
			}

			for _, request := range requests {
				if request.CreatorId != userId {
					return fmt.Errorf("requests creatorId: %s is not equal to the userId: %s", request.CreatorId, userId)
				}
			}

			if len(requests) != limit10 {
				return fmt.Errorf("requests should be %d", limit10)
			}

			isAscendingOrder := sort.SliceIsSorted(requests, func(p, q int) bool {
				return requests[p].Deadline < requests[q].Deadline
			})

			if !isAscendingOrder {
				return errors.New("requests are not in ascending order")
			}

			return nil
		}),
	)

	routePathDescendingOrder := fmt.Sprintf("%s?limit=%d&page=%d&asc=false", createRoutePath, limit10, page1)

	Test(t,
		Description("get owned requests; success; ascending order"),
		Get(routePathDescendingOrder),
		Send().Headers("Cookie").Add(sessionCookie),
		Expect().Status().Equal(http.StatusOK),
		Expect().Custom(func(hit Hit) error {
			var requests []entity.Request
			err := hit.Response().Body().JSON().Decode(&requests)
			if err != nil {
				return err
			}

			for _, request := range requests {
				if request.CreatorId != userId {
					return fmt.Errorf("requests creatorId: %s is not equal to the userId: %s", request.CreatorId, userId)
				}
			}

			if len(requests) != limit10 {
				return fmt.Errorf("requests should be %d", limit10)
			}

			isAscendingOrder := sort.SliceIsSorted(requests, func(p, q int) bool {
				return requests[p].Deadline < requests[q].Deadline
			})

			if isAscendingOrder {
				return errors.New("requests are not in descending order")
			}

			return nil
		}),
	)
}

// HTTP POST: /bidding/
func TestHTTPCreateBid(t *testing.T) {
	routePath := biddingApiPath + "/"
	sessionCookie := fmt.Sprintf(`s.id=%s`, sessionId)

	testdata.MockBid["RequestId"] = requestId

	Test(t,
		Description("bid; create; success"),
		Post(routePath),
		Send().Headers("Cookie").Add(sessionCookie),
		Send().Body().JSON(testdata.MockBid),
		Expect().Status().Equal(http.StatusOK),
		Expect().Custom(func(hit Hit) error {
			var bid entity.Bid

			err := hit.Response().Body().JSON().Decode(&bid)
			if err != nil {
				return err
			}

			if bid.Id < 1 {
				return errors.New("bid id is not correct")
			}

			return nil
		}),
	)

	testdata.MockBid["RequestId"] = 0
	Test(t,
		Description("bid; create; failure; RequestId not valid"),
		Post(routePath),
		Send().Headers("Cookie").Add(sessionCookie),
		Send().Body().JSON(testdata.MockBid),
		Expect().Status().Equal(http.StatusInternalServerError),
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
