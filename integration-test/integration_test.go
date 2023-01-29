package integration_test

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
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
var bidId = 0

var adminSessionId = ""

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
			}

			if loginSessionId == "" {
				return errors.New("Session is missing")
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

// HTTP GET: /request/count
func TestHTTPCountAllRequests(t *testing.T) {
	countOwnRoutePath := requestApiPath + "/count"
	sessionCookie := fmt.Sprintf(`s.id=%s`, sessionId)

	var allRequests = 10

	Test(t,
		Description("count all requests; success"),
		Get(countOwnRoutePath),
		Send().Headers("Cookie").Add(sessionCookie),
		Expect().Status().Equal(http.StatusOK),
		Expect().Custom(func(hit Hit) error {
			var count int
			err := hit.Response().Body().JSON().Decode(&count)
			if err != nil {
				return err
			}

			if count != allRequests {
				return fmt.Errorf("requests should be %d", allRequests)
			}

			return nil
		}))
}

// // HTTP GET: /request/own
func TestHTTPGetPaginatedOwnRequests(t *testing.T) {
	getOwnRoutePath := requestApiPath + "/own"
	sessionCookie := fmt.Sprintf(`s.id=%s`, sessionId)

	limit10 := 10
	page1 := 1

	routePathAscendingOrder := fmt.Sprintf("%s?limit=%d&page=%d&asc=true", getOwnRoutePath, limit10, page1)

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

	routePathDescendingOrder := fmt.Sprintf("%s?limit=%d&page=%d&asc=false", getOwnRoutePath, limit10, page1)

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

// HTTP GET: /request/count/own
func TestHTTPCountOwnRequests(t *testing.T) {
	countOwnRoutePath := requestApiPath + "/count/own"
	sessionCookie := fmt.Sprintf(`s.id=%s`, sessionId)

	var ownedRequests = 10

	Test(t,
		Description("count owned requests; success"),
		Get(countOwnRoutePath),
		Send().Headers("Cookie").Add(sessionCookie),
		Expect().Status().Equal(http.StatusOK),
		Expect().Custom(func(hit Hit) error {
			var count int
			err := hit.Response().Body().JSON().Decode(&count)
			if err != nil {
				return err
			}

			if count != ownedRequests {
				return fmt.Errorf("requests should be %d", ownedRequests)
			}

			return nil
		}))
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

			bidId = bid.Id

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

// HTTP GET: /bidding/count/own
func TestHTTPCountOwnnBids(t *testing.T) {
	countOwnRoutePath := biddingApiPath + "/count/own"
	createRoutePath := biddingApiPath + "/"
	sessionCookie := fmt.Sprintf(`s.id=%s`, sessionId)

	var newBid = map[string]interface{}{"RequestId": requestId, "Amount": 100.0}
	for i := 2; i <= 10; i++ {
		description := fmt.Sprintf("bid; create; success; no %d", i)
		newBid["Amount"] = i

		Test(t,
			Description(description),
			Post(createRoutePath),
			Send().Headers("Cookie").Add(sessionCookie),
			Send().Body().JSON(newBid),
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

	}

	var ownedBids = 10

	Test(t,
		Description("count owned bids; success"),
		Get(countOwnRoutePath),
		Send().Headers("Cookie").Add(sessionCookie),
		Expect().Status().Equal(http.StatusOK),
		Expect().Custom(func(hit Hit) error {
			var count int
			err := hit.Response().Body().JSON().Decode(&count)
			if err != nil {
				return err
			}

			if count != ownedBids {
				return fmt.Errorf("bids should be %d", ownedBids)
			}

			return nil
		}))
}

// // HTTP GET: /bidding/own
func TestHTTPGetPaginatedOwnBids(t *testing.T) {
	getOwnRoutePath := biddingApiPath + "/own"
	sessionCookie := fmt.Sprintf(`s.id=%s`, sessionId)

	limit10 := 10
	page1 := 1

	routePathAscendingOrder := fmt.Sprintf("%s?limit=%d&page=%d&asc=true", getOwnRoutePath, limit10, page1)

	Test(t,
		Description("get owned bids; success; ascending order"),
		Get(routePathAscendingOrder),
		Send().Headers("Cookie").Add(sessionCookie),
		Expect().Status().Equal(http.StatusOK),
		Expect().Custom(func(hit Hit) error {
			var bids []entity.Bid
			err := hit.Response().Body().JSON().Decode(&bids)
			if err != nil {
				return err
			}

			for _, bid := range bids {
				if bid.CreatorId != userId {
					return fmt.Errorf("bids creatorId: %s is not equal to the userId: %s", bid.CreatorId, userId)
				}
			}

			if len(bids) != limit10 {
				return fmt.Errorf("bids should be %d", limit10)
			}

			isAscendingOrder := sort.SliceIsSorted(bids, func(p, q int) bool {
				return bids[p].Amount < bids[q].Amount
			})

			if !isAscendingOrder {
				return errors.New("bids are not in ascending order")
			}

			return nil
		}),
	)

	routePathDescendingOrder := fmt.Sprintf("%s?limit=%d&page=%d&asc=false", getOwnRoutePath, limit10, page1)

	Test(t,
		Description("get owned bids; success; descending order"),
		Get(routePathDescendingOrder),
		Send().Headers("Cookie").Add(sessionCookie),
		Expect().Status().Equal(http.StatusOK),
		Expect().Custom(func(hit Hit) error {
			var bids []entity.Bid
			err := hit.Response().Body().JSON().Decode(&bids)
			if err != nil {
				return err
			}

			for _, request := range bids {
				if request.CreatorId != userId {
					return fmt.Errorf("bids creatorId: %s is not equal to the userId: %s", request.CreatorId, userId)
				}
			}

			if len(bids) != limit10 {
				return fmt.Errorf("bids should be %d", limit10)
			}

			isDescendingOrder := sort.SliceIsSorted(bids, func(p, q int) bool {
				return bids[p].Amount <= bids[q].Amount
			})

			if isDescendingOrder {
				return errors.New("bids are not in descending order")
			}

			return nil
		}),
	)
}

// HTTP GET: /bidding/?id
func TestHTTPGetBidById(t *testing.T) {
	createRoutePath := biddingApiPath + "/"
	sessionCookie := fmt.Sprintf(`s.id=%s`, sessionId)

	routePath := fmt.Sprintf("%s?id=%s", createRoutePath, strconv.Itoa(bidId))

	Test(t,
		Description("bid; get bid by id; success"),
		Get(routePath),
		Send().Headers("Cookie").Add(sessionCookie),
		Send().Body().JSON(testdata.MockBid),
		Expect().Status().Equal(http.StatusOK),
		Expect().Custom(func(hit Hit) error {
			var bid entity.Bid

			err := hit.Response().Body().JSON().Decode(&bid)
			if err != nil {
				return err
			}

			if bid.Amount != testdata.MockBid["Amount"] || bid.RequestId != requestId || bid.Id != bidId {
				log.Fatal("bid data do not match")
			}

			return nil
		}),
	)

	var incorrectBidId = 0
	incorrectRoutePath := fmt.Sprintf("%s?id=%s", createRoutePath, strconv.Itoa(incorrectBidId))

	Test(t,
		Description("bid; get bid by id; failure"),
		Get(incorrectRoutePath),
		Send().Headers("Cookie").Add(sessionCookie),
		Expect().Status().Equal(http.StatusInternalServerError),
	)
}

// HTTP GET: /bidding/requestId/?requestId
func TestHTTPGetPaginatedBidsByRequestId(t *testing.T) {
	createRoutePath := biddingApiPath + "/"
	sessionCookie := fmt.Sprintf(`s.id=%s`, sessionId)

	testdata.MockBid["RequestId"] = requestId

	for i := 2; i <= 10; i++ {
		description := fmt.Sprintf("bid; create; success; no %d", i)

		testdata.MockBid["Amount"] = i

		Test(t,
			Description(description),
			Post(createRoutePath),
			Send().Headers("Cookie").Add(sessionCookie),
			Send().Body().JSON(testdata.MockBid),
			Expect().Status().Equal(http.StatusOK),
		)

	}

	limit10 := 10
	page1 := 1

	routePathAscendingOrder := fmt.Sprintf("%srequestId/?limit=%d&page=%d&asc=true&requestId=%s", createRoutePath, limit10, page1, strconv.Itoa(requestId))

	Test(t,
		Description("get bids; success; ascending order"),
		Get(routePathAscendingOrder),
		Send().Headers("Cookie").Add(sessionCookie),
		Expect().Status().Equal(http.StatusOK),
		Expect().Custom(func(hit Hit) error {
			var bids []entity.Bid
			err := hit.Response().Body().JSON().Decode(&bids)
			if err != nil {
				return err
			}

			if len(bids) != limit10 {
				return fmt.Errorf("bids should be %d", limit10)
			}

			isAscendingOrder := sort.SliceIsSorted(bids, func(p, q int) bool {
				return bids[p].Id < bids[q].Id
			})

			if !isAscendingOrder {
				return errors.New("bids are not in ascending order")
			}

			return nil
		}),
	)

	routePathIncorrectRequestId := fmt.Sprintf("%srequestId/?limit=%d&page=%d&asc=true&requestId=%s", createRoutePath, limit10, page1, strconv.Itoa(0))

	Test(t,
		Description("get requests; failure; requestId does not exist"),
		Get(routePathIncorrectRequestId),
		Send().Headers("Cookie").Add(sessionCookie),
		Expect().Status().Equal(http.StatusOK),
		Expect().Custom(func(hit Hit) error {
			var resp interface{}
			err := hit.Response().Body().JSON().Decode(&resp)
			if err != nil {
				return err
			}

			if resp != nil {
				return errors.New("bids shouldn't be returned")
			}

			return nil
		}),
	)
}

// HTTP POST: /request/update/winner
func TestHTTPUpdateWinner(t *testing.T) {
	loginRoutePath := userApiPath + "/login"
	nonAdminSessionCookie := fmt.Sprintf(`s.id=%s`, sessionId)
	var yesterdayRequestId = 0
	var tomorrowRequestId = 0
	var mockBidId = 0

	Test(t,
		Description("login admin user; succes"),
		Post(loginRoutePath),
		Send().Headers("Content-Type").Add("application/json"),
		Send().Body().JSON(testdata.AdminUser),
		Expect().Status().Equal(http.StatusOK),
		Expect().Body().String().Contains("Successfully authenticated user"),
		Expect().Custom(func(hit Hit) error {
			var cookies = hit.Response().Cookies()

			for _, c := range cookies {
				if c.Name == "s.id" {
					adminSessionId = c.Value
				}
			}

			if adminSessionId == "" {
				return errors.New("Session is missing")
			}
			return nil
		}),
	)

	createReqeustPath := requestApiPath + "/"

	Test(t,
		Description("create request with deadline yesterday; success"),
		Post(createReqeustPath),
		Send().Headers("Cookie").Add(nonAdminSessionCookie),
		Send().Body().JSON(testdata.MockRequestYesterday),
		Expect().Status().Equal(http.StatusOK),
		Expect().Custom(func(hit Hit) error {
			var request entity.Request

			err := hit.Response().Body().JSON().Decode(&request)
			if err != nil {
				return err
			}

			yesterdayRequestId = request.Id

			return nil
		}),
	)

	Test(t,
		Description("create request with deadline tomorrow; success"),
		Post(createReqeustPath),
		Send().Headers("Cookie").Add(nonAdminSessionCookie),
		Send().Body().JSON(testdata.MockRequestTomorrow),
		Expect().Status().Equal(http.StatusOK),
		Expect().Custom(func(hit Hit) error {
			var request entity.Request

			err := hit.Response().Body().JSON().Decode(&request)
			if err != nil {
				return err
			}

			tomorrowRequestId = request.Id

			return nil
		}),
	)

	createBidRoutePath := biddingApiPath + "/"
	var mockBid = map[string]interface{}{"RequestId": yesterdayRequestId, "Amount": 100.0}
	Test(t,
		Description("bid; create; success"),
		Post(createBidRoutePath),
		Send().Headers("Cookie").Add(nonAdminSessionCookie),
		Send().Body().JSON(mockBid),
		Expect().Status().Equal(http.StatusOK),
		Expect().Custom(func(hit Hit) error {
			var bid entity.Bid

			err := hit.Response().Body().JSON().Decode(&bid)
			if err != nil {
				return err
			}

			mockBidId = bid.Id

			return nil
		}),
	)

	routePath := requestApiPath + "/update/winner?requestId=" + strconv.Itoa(yesterdayRequestId)
	adminSessionCookie := fmt.Sprintf(`s.id=%s`, adminSessionId)

	Test(t,
		Description("non admin user; failure"),
		Post(routePath),
		Send().Headers("Cookie").Add(nonAdminSessionCookie),
		Expect().Status().Equal(http.StatusUnauthorized),
		Expect().Body().String().Contains("incorrect permissions"),
	)

	incorrectRoutePath := requestApiPath + "/update/winner?requestId=random"
	Test(t,
		Description("admin user; validation error; failure"),
		Post(incorrectRoutePath),
		Send().Headers("Cookie").Add(adminSessionCookie),
		Expect().Status().Equal(http.StatusBadRequest),
		Expect().Body().String().Contains("Validation error"),
	)

	incorrectRoutePath = requestApiPath + "/update/winner?requestId=100000000"
	Test(t,
		Description("admin user; request not found; failure"),
		Post(incorrectRoutePath),
		Send().Headers("Cookie").Add(adminSessionCookie),
		Expect().Status().Equal(http.StatusNotFound),
		Expect().Body().String().Contains("Request not found"),
	)

	notUpdateableRoutePath := requestApiPath + "/update/winner?requestId=" + strconv.Itoa(tomorrowRequestId)
	Test(t,
		Description("admin user; request cannot be update; deadline hasn't passed; failure"),
		Post(notUpdateableRoutePath),
		Send().Headers("Cookie").Add(adminSessionCookie),
		Expect().Status().Equal(http.StatusUnauthorized),
		Expect().Body().String().Contains("Request not allowed to be resolved"),
	)

	Test(
		t,
		Description("admin user; request updated; success"),
		Post(routePath),
		Send().Headers("Cookie").Add(adminSessionCookie),
		Expect().Status().Equal(http.StatusOK),
		Expect().Custom(func(hit Hit) error {
			var bid entity.Bid

			err := hit.Response().Body().JSON().Decode(&bid)
			if err != nil {
				return err
			}

			if bid.RequestId != yesterdayRequestId || bid.Id != mockBidId {
				return errors.New("returned request is incorrect")
			}

			return nil
		}),
	)

	Test(
		t,
		Description("admin user; resolved request cannot be updated; failure"),
		Post(routePath),
		Send().Headers("Cookie").Add(adminSessionCookie),
		Expect().Status().Equal(http.StatusUnauthorized),
		Expect().Body().String().Contains("Request not allowed to be resolved"),
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
