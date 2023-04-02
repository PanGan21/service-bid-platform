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
	"github.com/PanGan21/pkg/auth"
)

var userService = "user"
var auctionService = "auction"
var biddingService = "bidding"

var sessionId = ""
var userId = ""
var auctionId = 0
var bidId = 0

var adminSessionId = ""

var userApiPath = getBasePath(userService)
var auctionApiPath = getBasePath(auctionService)
var biddingApiPath = getBasePath(biddingService)

func TestMain(m *testing.M) {
	fmt.Println("Sleep for 50 seconds to allow services and kafka stabilize")
	time.Sleep(50 * time.Second)
	fmt.Println("Start integration tests")

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
func TestHTTPDoGetLoggedInDetails(t *testing.T) {
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

// HTTP GET: /user/details
func TestHTTPDoGetDetails(t *testing.T) {
	sessionCookie := fmt.Sprintf(`s.id=%s`, sessionId)
	routePath := userApiPath + "/details?userId=" + userId

	Test(t,
		Description("get user details; success; user exists"),
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

	incorrectRoutePath := userApiPath + "/details?userId=" + "100000000"

	Test(t,
		Description("get user details; failure; user does not exists"),
		Get(incorrectRoutePath),
		Send().Headers("Cookie").Add(sessionCookie),
		Expect().Status().Equal(http.StatusUnauthorized),
		Expect().Body().String().Contains("unauthorized"),
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
		Description("login admin user; succes"),
		Post(routePath),
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

// HTTP GET: /auction/hello
func TestHTTPDoHello(t *testing.T) {
	routePath := auctionApiPath + "/hello"
	sessionCookie := fmt.Sprintf(`s.id=%s`, sessionId)

	Test(t,
		Description("auction; hello; success"),
		Get(routePath),
		Send().Headers("Cookie").Add(sessionCookie),
		Expect().Status().Equal(http.StatusOK),
	)
}

// HTTP POST: /auction/
// func TestHTTPCreateAuction(t *testing.T) {
// 	routePath := auctionApiPath + "/"
// 	sessionCookie := fmt.Sprintf(`s.id=%s`, sessionId)

// 	Test(t,
// 		Description("auction; create; success"),
// 		Post(routePath),
// 		Send().Headers("Cookie").Add(sessionCookie),
// 		Send().Body().JSON(testdata.MockAuction),
// 		Expect().Status().Equal(http.StatusOK),
// 		Expect().Custom(func(hit Hit) error {
// 			var auction entity.Auction

// 			err := hit.Response().Body().JSON().Decode(&auction)
// 			if err != nil {
// 				return err
// 			}

// 			if auction.Status != entity.New {
// 				return fmt.Errorf("created auction status should be %s", entity.New)
// 			}

// 			auctionId = auction.Id

// 			return nil
// 		}),
// 	)

// 	err := waitUntilAuctionIsAvailableInBidding(50, auctionId)
// 	if err != nil {
// 		log.Fatal(err)
// 		t.Fail()
// 	}
// }

// HTTP POST: /auction/update/status
// func TestHTTPUpdateAuctionStatusToOpen(t *testing.T) {
// 	routePath := auctionApiPath + "/update/status?auctionId=" + strconv.Itoa(auctionId)
// 	adminSessionCookie := fmt.Sprintf(`s.id=%s`, adminSessionId)

// 	var updateStatusBody = map[string]interface{}{"Status": entity.Open}

// 	Test(
// 		t,
// 		Description("update auction status; open; success"),
// 		Post(routePath),
// 		Send().Body().JSON(updateStatusBody),
// 		Send().Headers("Cookie").Add(adminSessionCookie),
// 		Expect().Status().Equal(http.StatusOK),
// 		Expect().Custom(func(hit Hit) error {
// 			var auction entity.Auction
// 			err := hit.Response().Body().JSON().Decode(&auction)
// 			if err != nil {
// 				return err
// 			}

// 			if auction.Status != entity.Open {
// 				return fmt.Errorf("auction should have been update to status %s", entity.Open)
// 			}

// 			return nil
// 		}),
// 	)
// }

// // HTTP GET: /auction/
// func TestHTTPGetPaginatedAuctions(t *testing.T) {
// 	createRoutePath := auctionApiPath + "/"
// 	sessionCookie := fmt.Sprintf(`s.id=%s`, sessionId)

// 	for i := 2; i <= 10; i++ {
// 		description := fmt.Sprintf("auction; create; success; no %d", i)
// 		testAuction := testdata.MockAuction
// 		testAuction["deadline"] = i

// 		Test(t,
// 			Description(description),
// 			Post(createRoutePath),
// 			Send().Headers("Cookie").Add(sessionCookie),
// 			Send().Body().JSON(testAuction),
// 			Expect().Status().Equal(http.StatusOK),
// 		)

// 	}

// 	limit10 := 10
// 	page1 := 1

// 	routePathAscendingOrder := fmt.Sprintf("%s?limit=%d&page=%d&asc=true", createRoutePath, limit10, page1)

// 	Test(t,
// 		Description("get auctions; success; ascending order"),
// 		Get(routePathAscendingOrder),
// 		Send().Headers("Cookie").Add(sessionCookie),
// 		Expect().Status().Equal(http.StatusOK),
// 		Expect().Custom(func(hit Hit) error {
// 			var auctions []entity.Auction
// 			err := hit.Response().Body().JSON().Decode(&auctions)
// 			if err != nil {
// 				return err
// 			}

// 			if len(auctions) != limit10 {
// 				return fmt.Errorf("auctions should be %d", limit10)
// 			}

// 			isAscendingOrder := sort.SliceIsSorted(auctions, func(p, q int) bool {
// 				return auctions[p].Deadline < auctions[q].Deadline
// 			})

// 			if !isAscendingOrder {
// 				return errors.New("auctions are not in ascending order")
// 			}

// 			return nil
// 		}),
// 	)

// 	routePathDescendingOrder := fmt.Sprintf("%s?limit=%d&page=%d&asc=false", createRoutePath, limit10, page1)

// 	Test(t,
// 		Description("get auctions; success; descending order"),
// 		Get(routePathDescendingOrder),
// 		Send().Headers("Cookie").Add(sessionCookie),
// 		Expect().Status().Equal(http.StatusOK),
// 		Expect().Custom(func(hit Hit) error {
// 			var auctions []entity.Auction
// 			err := hit.Response().Body().JSON().Decode(&auctions)
// 			if err != nil {
// 				return err
// 			}

// 			if len(auctions) != limit10 {
// 				return fmt.Errorf("auctions should be %d", limit10)
// 			}

// 			isAscendingOrder := sort.SliceIsSorted(auctions, func(p, q int) bool {
// 				return auctions[p].Deadline < auctions[q].Deadline
// 			})

// 			if isAscendingOrder {
// 				return errors.New("auctions are not in descending order")
// 			}

// 			return nil
// 		}),
// 	)
// }

// // HTTP GET: /auction/count
// func TestHTTPCountAllAuctions(t *testing.T) {
// 	countOwnRoutePath := auctionApiPath + "/count"
// 	sessionCookie := fmt.Sprintf(`s.id=%s`, sessionId)

// 	var allAuctions = 10

// 	Test(t,
// 		Description("count all auctions; success"),
// 		Get(countOwnRoutePath),
// 		Send().Headers("Cookie").Add(sessionCookie),
// 		Expect().Status().Equal(http.StatusOK),
// 		Expect().Custom(func(hit Hit) error {
// 			var count int
// 			err := hit.Response().Body().JSON().Decode(&count)
// 			if err != nil {
// 				return err
// 			}

// 			if count != allAuctions {
// 				return fmt.Errorf("auctions should be %d", allAuctions)
// 			}

// 			return nil
// 		}))
// }

// // // HTTP GET: /auction/own
// func TestHTTPGetPaginatedOwnAuctions(t *testing.T) {
// 	getOwnRoutePath := auctionApiPath + "/own"
// 	sessionCookie := fmt.Sprintf(`s.id=%s`, sessionId)

// 	limit10 := 10
// 	page1 := 1

// 	routePathAscendingOrder := fmt.Sprintf("%s?limit=%d&page=%d&asc=true", getOwnRoutePath, limit10, page1)

// 	Test(t,
// 		Description("get owned auctions; success; ascending order"),
// 		Get(routePathAscendingOrder),
// 		Send().Headers("Cookie").Add(sessionCookie),
// 		Expect().Status().Equal(http.StatusOK),
// 		Expect().Custom(func(hit Hit) error {
// 			var auctions []entity.Auction
// 			err := hit.Response().Body().JSON().Decode(&auctions)
// 			if err != nil {
// 				return err
// 			}

// 			for _, auction := range auctions {
// 				if auction.CreatorId != userId {
// 					return fmt.Errorf("auctions creatorId: %s is not equal to the userId: %s", auction.CreatorId, userId)
// 				}
// 			}

// 			if len(auctions) != limit10 {
// 				return fmt.Errorf("auctions should be %d", limit10)
// 			}

// 			isAscendingOrder := sort.SliceIsSorted(auctions, func(p, q int) bool {
// 				return auctions[p].Deadline < auctions[q].Deadline
// 			})

// 			if !isAscendingOrder {
// 				return errors.New("auctions are not in ascending order")
// 			}

// 			return nil
// 		}),
// 	)

// 	routePathDescendingOrder := fmt.Sprintf("%s?limit=%d&page=%d&asc=false", getOwnRoutePath, limit10, page1)

// 	Test(t,
// 		Description("get owned auctions; success; ascending order"),
// 		Get(routePathDescendingOrder),
// 		Send().Headers("Cookie").Add(sessionCookie),
// 		Expect().Status().Equal(http.StatusOK),
// 		Expect().Custom(func(hit Hit) error {
// 			var auctions []entity.Auction
// 			err := hit.Response().Body().JSON().Decode(&auctions)
// 			if err != nil {
// 				return err
// 			}

// 			for _, auction := range auctions {
// 				if auction.CreatorId != userId {
// 					return fmt.Errorf("auctions creatorId: %s is not equal to the userId: %s", auction.CreatorId, userId)
// 				}
// 			}

// 			if len(auctions) != limit10 {
// 				return fmt.Errorf("auctions should be %d", limit10)
// 			}

// 			isAscendingOrder := sort.SliceIsSorted(auctions, func(p, q int) bool {
// 				return auctions[p].Deadline < auctions[q].Deadline
// 			})

// 			if isAscendingOrder {
// 				return errors.New("auctions are not in descending order")
// 			}

// 			return nil
// 		}),
// 	)
// }

// // HTTP GET: /auction/count/own
// func TestHTTPCountOwnAuctions(t *testing.T) {
// 	countOwnRoutePath := auctionApiPath + "/count/own"
// 	sessionCookie := fmt.Sprintf(`s.id=%s`, sessionId)

// 	var ownedAuctions = 10

// 	Test(t,
// 		Description("count owned auctions; success"),
// 		Get(countOwnRoutePath),
// 		Send().Headers("Cookie").Add(sessionCookie),
// 		Expect().Status().Equal(http.StatusOK),
// 		Expect().Custom(func(hit Hit) error {
// 			var count int
// 			err := hit.Response().Body().JSON().Decode(&count)
// 			if err != nil {
// 				return err
// 			}

// 			if count != ownedAuctions {
// 				return fmt.Errorf("auctions should be %d", ownedAuctions)
// 			}

// 			return nil
// 		}))
// }

// // HTTP POST: /bidding/
// func TestHTTPCreateBid(t *testing.T) {
// 	routePath := biddingApiPath + "/"
// 	sessionCookie := fmt.Sprintf(`s.id=%s`, sessionId)

// 	testdata.MockBid["AuctionId"] = auctionId

// 	waitUntilAuctionIsOpenToBids(50, auctionId)

// 	Test(t,
// 		Description("bid; create; success"),
// 		Post(routePath),
// 		Send().Headers("Cookie").Add(sessionCookie),
// 		Send().Body().JSON(testdata.MockBid),
// 		Expect().Status().Equal(http.StatusOK),
// 		Expect().Custom(func(hit Hit) error {
// 			var bid entity.Bid

// 			err := hit.Response().Body().JSON().Decode(&bid)
// 			if err != nil {
// 				return err
// 			}

// 			if bid.Id < 1 {
// 				return errors.New("bid id is not correct")
// 			}

// 			bidId = bid.Id

// 			return nil
// 		}),
// 	)

// 	testdata.MockBid["AuctionId"] = 0
// 	Test(t,
// 		Description("bid; create; failure; AuctionId not valid"),
// 		Post(routePath),
// 		Send().Headers("Cookie").Add(sessionCookie),
// 		Send().Body().JSON(testdata.MockBid),
// 		Expect().Status().Equal(http.StatusUnauthorized),
// 		Expect().Body().String().Contains("Auction doesn't receive bids"),
// 	)
// }

// // HTTP GET: /bidding/count/own
// func TestHTTPCountOwnnBids(t *testing.T) {
// 	countOwnRoutePath := biddingApiPath + "/count/own"
// 	createRoutePath := biddingApiPath + "/"
// 	sessionCookie := fmt.Sprintf(`s.id=%s`, sessionId)

// 	var newBid = map[string]interface{}{"AuctionId": auctionId, "Amount": 100.0}
// 	for i := 2; i <= 10; i++ {
// 		description := fmt.Sprintf("bid; create; success; no %d", i)
// 		newBid["Amount"] = i

// 		Test(t,
// 			Description(description),
// 			Post(createRoutePath),
// 			Send().Headers("Cookie").Add(sessionCookie),
// 			Send().Body().JSON(newBid),
// 			Expect().Status().Equal(http.StatusOK),
// 			Expect().Custom(func(hit Hit) error {
// 				var bid entity.Bid

// 				err := hit.Response().Body().JSON().Decode(&bid)
// 				if err != nil {
// 					return err
// 				}

// 				if bid.Id < 1 {
// 					return errors.New("bid id is not correct")
// 				}

// 				return nil
// 			}),
// 		)

// 	}

// 	var ownedBids = 10

// 	Test(t,
// 		Description("count owned bids; success"),
// 		Get(countOwnRoutePath),
// 		Send().Headers("Cookie").Add(sessionCookie),
// 		Expect().Status().Equal(http.StatusOK),
// 		Expect().Custom(func(hit Hit) error {
// 			var count int
// 			err := hit.Response().Body().JSON().Decode(&count)
// 			if err != nil {
// 				return err
// 			}

// 			if count != ownedBids {
// 				return fmt.Errorf("bids should be %d", ownedBids)
// 			}

// 			return nil
// 		}))
// }

// // HTTP GET: /bidding/own
// func TestHTTPGetPaginatedOwnBids(t *testing.T) {
// 	getOwnRoutePath := biddingApiPath + "/own"
// 	sessionCookie := fmt.Sprintf(`s.id=%s`, sessionId)

// 	limit10 := 10
// 	page1 := 1

// 	routePathAscendingOrder := fmt.Sprintf("%s?limit=%d&page=%d&asc=true", getOwnRoutePath, limit10, page1)

// 	Test(t,
// 		Description("get owned bids; success; ascending order"),
// 		Get(routePathAscendingOrder),
// 		Send().Headers("Cookie").Add(sessionCookie),
// 		Expect().Status().Equal(http.StatusOK),
// 		Expect().Custom(func(hit Hit) error {
// 			var bids []entity.Bid
// 			err := hit.Response().Body().JSON().Decode(&bids)
// 			if err != nil {
// 				return err
// 			}

// 			for _, bid := range bids {
// 				if bid.CreatorId != userId {
// 					return fmt.Errorf("bids creatorId: %s is not equal to the userId: %s", bid.CreatorId, userId)
// 				}
// 			}

// 			if len(bids) != limit10 {
// 				return fmt.Errorf("bids should be %d", limit10)
// 			}

// 			isAscendingOrder := sort.SliceIsSorted(bids, func(p, q int) bool {
// 				return bids[p].Amount < bids[q].Amount
// 			})

// 			if !isAscendingOrder {
// 				return errors.New("bids are not in ascending order")
// 			}

// 			return nil
// 		}),
// 	)

// 	routePathDescendingOrder := fmt.Sprintf("%s?limit=%d&page=%d&asc=false", getOwnRoutePath, limit10, page1)

// 	Test(t,
// 		Description("get owned bids; success; descending order"),
// 		Get(routePathDescendingOrder),
// 		Send().Headers("Cookie").Add(sessionCookie),
// 		Expect().Status().Equal(http.StatusOK),
// 		Expect().Custom(func(hit Hit) error {
// 			var bids []entity.Bid
// 			err := hit.Response().Body().JSON().Decode(&bids)
// 			if err != nil {
// 				return err
// 			}

// 			for _, auction := range bids {
// 				if auction.CreatorId != userId {
// 					return fmt.Errorf("bids creatorId: %s is not equal to the userId: %s", auction.CreatorId, userId)
// 				}
// 			}

// 			if len(bids) != limit10 {
// 				return fmt.Errorf("bids should be %d", limit10)
// 			}

// 			isDescendingOrder := sort.SliceIsSorted(bids, func(p, q int) bool {
// 				return bids[p].Amount <= bids[q].Amount
// 			})

// 			if isDescendingOrder {
// 				return errors.New("bids are not in descending order")
// 			}

// 			return nil
// 		}),
// 	)
// }

// // HTTP GET: /bidding/?id
// func TestHTTPGetBidById(t *testing.T) {
// 	createRoutePath := biddingApiPath + "/"
// 	sessionCookie := fmt.Sprintf(`s.id=%s`, sessionId)

// 	routePath := fmt.Sprintf("%s?id=%s", createRoutePath, strconv.Itoa(bidId))

// 	Test(t,
// 		Description("bid; get bid by id; success"),
// 		Get(routePath),
// 		Send().Headers("Cookie").Add(sessionCookie),
// 		Send().Body().JSON(testdata.MockBid),
// 		Expect().Status().Equal(http.StatusOK),
// 		Expect().Custom(func(hit Hit) error {
// 			var bid entity.Bid

// 			err := hit.Response().Body().JSON().Decode(&bid)
// 			if err != nil {
// 				return err
// 			}

// 			if bid.Amount != testdata.MockBid["Amount"] || bid.AuctionId != auctionId || bid.Id != bidId {
// 				log.Fatal("bid data do not match")
// 			}

// 			return nil
// 		}),
// 	)

// 	var incorrectBidId = 0
// 	incorrectRoutePath := fmt.Sprintf("%s?id=%s", createRoutePath, strconv.Itoa(incorrectBidId))

// 	Test(t,
// 		Description("bid; get bid by id; failure"),
// 		Get(incorrectRoutePath),
// 		Send().Headers("Cookie").Add(sessionCookie),
// 		Expect().Status().Equal(http.StatusInternalServerError),
// 	)
// }

// // HTTP GET: /bidding/auctionId/?auctionId
// func TestHTTPGetPaginatedBidsByAuctionId(t *testing.T) {
// 	createRoutePath := biddingApiPath + "/"
// 	sessionCookie := fmt.Sprintf(`s.id=%s`, sessionId)

// 	testdata.MockBid["AuctionId"] = auctionId

// 	for i := 2; i <= 10; i++ {
// 		description := fmt.Sprintf("bid; create; success; no %d", i)

// 		testdata.MockBid["Amount"] = i

// 		Test(t,
// 			Description(description),
// 			Post(createRoutePath),
// 			Send().Headers("Cookie").Add(sessionCookie),
// 			Send().Body().JSON(testdata.MockBid),
// 			Expect().Status().Equal(http.StatusOK),
// 		)

// 	}

// 	limit10 := 10
// 	page1 := 1

// 	routePathAscendingOrder := fmt.Sprintf("%sauctionId/?limit=%d&page=%d&asc=true&auctionId=%s", createRoutePath, limit10, page1, strconv.Itoa(auctionId))

// 	Test(t,
// 		Description("get bids; success; ascending order"),
// 		Get(routePathAscendingOrder),
// 		Send().Headers("Cookie").Add(sessionCookie),
// 		Expect().Status().Equal(http.StatusOK),
// 		Expect().Custom(func(hit Hit) error {
// 			var bids []entity.Bid
// 			err := hit.Response().Body().JSON().Decode(&bids)
// 			if err != nil {
// 				return err
// 			}

// 			if len(bids) != limit10 {
// 				return fmt.Errorf("bids should be %d", limit10)
// 			}

// 			isAscendingOrder := sort.SliceIsSorted(bids, func(p, q int) bool {
// 				return bids[p].Id < bids[q].Id
// 			})

// 			if !isAscendingOrder {
// 				return errors.New("bids are not in ascending order")
// 			}

// 			return nil
// 		}),
// 	)

// 	routePathIncorrectAuctionId := fmt.Sprintf("%sauctionId/?limit=%d&page=%d&asc=true&auctionId=%s", createRoutePath, limit10, page1, strconv.Itoa(0))

// 	Test(t,
// 		Description("get auctions; failure; auctionId does not exist"),
// 		Get(routePathIncorrectAuctionId),
// 		Send().Headers("Cookie").Add(sessionCookie),
// 		Expect().Status().Equal(http.StatusOK),
// 		Expect().Custom(func(hit Hit) error {
// 			var resp interface{}
// 			err := hit.Response().Body().JSON().Decode(&resp)
// 			if err != nil {
// 				return err
// 			}

// 			if resp != nil {
// 				return errors.New("bids shouldn't be returned")
// 			}

// 			return nil
// 		}),
// 	)
// }

// // HTTP POST: /auction/update/winner
// func TestHTTPUpdateWinner(t *testing.T) {
// 	nonAdminSessionCookie := fmt.Sprintf(`s.id=%s`, sessionId)
// 	var yesterdayAuctionId = 0
// 	var tomorrowAuctionId = 0

// 	var mockFirstLowestBidId = 0
// 	var mockSecondLowestBidId = 0
// 	var firstLowestAmount = 50.0
// 	var secondLowestAmount = 100.0

// 	createReqeustPath := auctionApiPath + "/"

// 	Test(t,
// 		Description("create auction with deadline yesterday; success"),
// 		Post(createReqeustPath),
// 		Send().Headers("Cookie").Add(nonAdminSessionCookie),
// 		Send().Body().JSON(testdata.MockAuctionYesterday),
// 		Expect().Status().Equal(http.StatusOK),
// 		Expect().Custom(func(hit Hit) error {
// 			var auction entity.Auction

// 			err := hit.Response().Body().JSON().Decode(&auction)
// 			if err != nil {
// 				return err
// 			}

// 			yesterdayAuctionId = auction.Id

// 			return nil
// 		}),
// 	)

// 	updateStatusRoutePath := auctionApiPath + "/update/status?auctionId=" + strconv.Itoa(yesterdayAuctionId)
// 	var updateStatusBody = map[string]interface{}{"Status": entity.Open}
// 	adminSessionCookie := fmt.Sprintf(`s.id=%s`, adminSessionId)
// 	Test(
// 		t,
// 		Description("update yesterday auction status; open; success"),
// 		Post(updateStatusRoutePath),
// 		Send().Body().JSON(updateStatusBody),
// 		Send().Headers("Cookie").Add(adminSessionCookie),
// 		Expect().Status().Equal(http.StatusOK),
// 		Expect().Custom(func(hit Hit) error {
// 			var auction entity.Auction
// 			err := hit.Response().Body().JSON().Decode(&auction)
// 			if err != nil {
// 				return err
// 			}

// 			if auction.Status != entity.Open {
// 				return fmt.Errorf("auction should have been update to status %s", entity.Open)
// 			}

// 			return nil
// 		}),
// 	)

// 	Test(t,
// 		Description("create auction with deadline tomorrow; success"),
// 		Post(createReqeustPath),
// 		Send().Headers("Cookie").Add(nonAdminSessionCookie),
// 		Send().Body().JSON(testdata.MockAuctionTomorrow),
// 		Expect().Status().Equal(http.StatusOK),
// 		Expect().Custom(func(hit Hit) error {
// 			var auction entity.Auction

// 			err := hit.Response().Body().JSON().Decode(&auction)
// 			if err != nil {
// 				return err
// 			}

// 			tomorrowAuctionId = auction.Id

// 			return nil
// 		}),
// 	)

// 	updateStatusRoutePath = auctionApiPath + "/update/status?auctionId=" + strconv.Itoa(tomorrowAuctionId)
// 	Test(
// 		t,
// 		Description("update yesterday auction status; open; success"),
// 		Post(updateStatusRoutePath),
// 		Send().Body().JSON(updateStatusBody),
// 		Send().Headers("Cookie").Add(adminSessionCookie),
// 		Expect().Status().Equal(http.StatusOK),
// 		Expect().Custom(func(hit Hit) error {
// 			var auction entity.Auction
// 			err := hit.Response().Body().JSON().Decode(&auction)
// 			if err != nil {
// 				return err
// 			}

// 			if auction.Status != entity.Open {
// 				return fmt.Errorf("auction should have been update to status %s", entity.Open)
// 			}

// 			return nil
// 		}),
// 	)

// 	createBidRoutePath := biddingApiPath + "/"
// 	var mockBid = map[string]interface{}{"AuctionId": yesterdayAuctionId, "Amount": firstLowestAmount}
// 	Test(t,
// 		Description("bid; create; success"),
// 		Post(createBidRoutePath),
// 		Send().Headers("Cookie").Add(nonAdminSessionCookie),
// 		Send().Body().JSON(mockBid),
// 		Expect().Status().Equal(http.StatusOK),
// 		Expect().Custom(func(hit Hit) error {
// 			var bid entity.Bid

// 			err := hit.Response().Body().JSON().Decode(&bid)
// 			if err != nil {
// 				return err
// 			}

// 			mockFirstLowestBidId = bid.Id

// 			return nil
// 		}),
// 	)

// 	var mockSecondBid = map[string]interface{}{"AuctionId": yesterdayAuctionId, "Amount": secondLowestAmount}
// 	Test(t,
// 		Description("bid; create; success"),
// 		Post(createBidRoutePath),
// 		Send().Headers("Cookie").Add(nonAdminSessionCookie),
// 		Send().Body().JSON(mockSecondBid),
// 		Expect().Status().Equal(http.StatusOK),
// 		Expect().Custom(func(hit Hit) error {
// 			var bid entity.Bid

// 			err := hit.Response().Body().JSON().Decode(&bid)
// 			if err != nil {
// 				return err
// 			}

// 			mockSecondLowestBidId = bid.Id

// 			return nil
// 		}),
// 	)

// 	err := waitUntilBidIsAvailableInAuction(50, mockFirstLowestBidId)
// 	if err != nil {
// 		log.Fatal(err)
// 		t.Fail()
// 	}

// 	err = waitUntilBidIsAvailableInAuction(50, mockSecondLowestBidId)
// 	if err != nil {
// 		log.Fatal(err)
// 		t.Fail()
// 	}

// 	routePath := auctionApiPath + "/update/winner?auctionId=" + strconv.Itoa(yesterdayAuctionId)

// 	Test(t,
// 		Description("non admin user; failure"),
// 		Post(routePath),
// 		Send().Headers("Cookie").Add(nonAdminSessionCookie),
// 		Expect().Status().Equal(http.StatusUnauthorized),
// 		Expect().Body().String().Contains("incorrect permissions"),
// 	)

// 	incorrectRoutePath := auctionApiPath + "/update/winner?auctionId=random"
// 	Test(t,
// 		Description("admin user; validation error; failure"),
// 		Post(incorrectRoutePath),
// 		Send().Headers("Cookie").Add(adminSessionCookie),
// 		Expect().Status().Equal(http.StatusBadRequest),
// 		Expect().Body().String().Contains("Validation error"),
// 	)

// 	incorrectRoutePath = auctionApiPath + "/update/winner?auctionId=100000000"
// 	Test(t,
// 		Description("admin user; auction not found; failure"),
// 		Post(incorrectRoutePath),
// 		Send().Headers("Cookie").Add(adminSessionCookie),
// 		Expect().Status().Equal(http.StatusNotFound),
// 		Expect().Body().String().Contains("Auction not found"),
// 	)

// 	notUpdateableRoutePath := auctionApiPath + "/update/winner?auctionId=" + strconv.Itoa(tomorrowAuctionId)
// 	Test(t,
// 		Description("admin user; auction cannot be update; deadline hasn't passed; failure"),
// 		Post(notUpdateableRoutePath),
// 		Send().Headers("Cookie").Add(adminSessionCookie),
// 		Expect().Status().Equal(http.StatusUnauthorized),
// 		Expect().Body().String().Contains("Auction not allowed to be resolved"),
// 	)

// 	Test(
// 		t,
// 		Description("admin user; auction updated; success"),
// 		Post(routePath),
// 		Send().Headers("Cookie").Add(adminSessionCookie),
// 		Expect().Status().Equal(http.StatusOK),
// 		Expect().Custom(func(hit Hit) error {
// 			var auction entity.Auction

// 			err := hit.Response().Body().JSON().Decode(&auction)
// 			if err != nil {
// 				return err
// 			}

// 			if auction.Id != yesterdayAuctionId || auction.WinningBidId != strconv.Itoa(mockFirstLowestBidId) || auction.WinnerId != userId || auction.WinningAmount != secondLowestAmount {
// 				return errors.New("returned auction is incorrect")
// 			}

// 			return nil
// 		}),
// 	)

// 	Test(
// 		t,
// 		Description("admin user; resolved auction cannot be updated again; failure"),
// 		Post(routePath),
// 		Send().Headers("Cookie").Add(adminSessionCookie),
// 		Expect().Status().Equal(http.StatusUnauthorized),
// 		Expect().Body().String().Contains("Auction not allowed to be resolved"),
// 	)
// }

// // HTTP GET: /auction/open/past-deadline
// func TestHTTPGetOpenPastDeadlineAuctions(t *testing.T) {
// 	baseRoutePath := auctionApiPath + "/open/past-deadline"
// 	adminSessionCookie := fmt.Sprintf(`s.id=%s`, adminSessionId)

// 	limit10 := 10
// 	page1 := 1

// 	routePathAscendingOrder := fmt.Sprintf("%s?limit=%d&page=%d&asc=true", baseRoutePath, limit10, page1)

// 	now := time.Now().UTC().UnixMilli()

// 	Test(
// 		t,
// 		Description("get open auctions past dealine; asc order; success"),
// 		Get(routePathAscendingOrder),
// 		Send().Headers("Cookie").Add(adminSessionCookie),
// 		Expect().Status().Equal(http.StatusOK),
// 		Expect().Custom(func(hit Hit) error {
// 			var auctions []entity.Auction
// 			err := hit.Response().Body().JSON().Decode(&auctions)
// 			if err != nil {
// 				return err
// 			}

// 			for _, req := range auctions {
// 				if req.Status != entity.Open || req.Deadline >= now {
// 					return fmt.Errorf("auction with id %d is not open or not past the deadline", req.Id)
// 				}
// 			}

// 			if len(auctions) != 1 {
// 				return fmt.Errorf("auctions should be %d", 1)
// 			}

// 			isAscendingOrder := sort.SliceIsSorted(auctions, func(p, q int) bool {
// 				return auctions[p].Deadline < auctions[q].Deadline
// 			})

// 			if !isAscendingOrder {
// 				return errors.New("auctions are not in ascending order")
// 			}

// 			return nil
// 		}),
// 	)

// 	routePathDescendingOrder := fmt.Sprintf("%s?limit=%d&page=%d&asc=false", baseRoutePath, limit10, page1)

// 	Test(
// 		t,
// 		Description("get open auctions past dealine; desc order; success"),
// 		Get(routePathDescendingOrder),
// 		Send().Headers("Cookie").Add(adminSessionCookie),
// 		Expect().Status().Equal(http.StatusOK),
// 		Expect().Custom(func(hit Hit) error {
// 			var auctions []entity.Auction
// 			err := hit.Response().Body().JSON().Decode(&auctions)
// 			if err != nil {
// 				return err
// 			}

// 			for _, req := range auctions {
// 				if req.Status != entity.Open || req.Deadline >= now {
// 					return fmt.Errorf("auction with id %d is not open or not past the deadline", req.Id)
// 				}
// 			}

// 			if len(auctions) != 1 {
// 				return fmt.Errorf("auctions should be %d", 1)
// 			}

// 			return nil
// 		}),
// 	)
// }

// // HTTP GET: /auction/open/past-deadline
// func TestHTTPCountOpenPastDeadlineAuctions(t *testing.T) {
// 	routePath := auctionApiPath + "/open/past-deadline/count"
// 	adminSessionCookie := fmt.Sprintf(`s.id=%s`, adminSessionId)

// 	var pastDeadlineOpenAuctions = 1

// 	Test(
// 		t,
// 		Description("count open auctions past dealine; asc order; success"),
// 		Get(routePath),
// 		Send().Headers("Cookie").Add(adminSessionCookie),
// 		Expect().Status().Equal(http.StatusOK),
// 		Expect().Custom(func(hit Hit) error {
// 			var count int
// 			err := hit.Response().Body().JSON().Decode(&count)
// 			if err != nil {
// 				return err
// 			}

// 			if count != pastDeadlineOpenAuctions {
// 				return fmt.Errorf("past deadline open auctions should be %d", pastDeadlineOpenAuctions)
// 			}

// 			return nil
// 		}),
// 	)
// }

// // HTTP POST: /auction/update/status
// func TestHTTPUpdateAuctionStatus(t *testing.T) {
// 	routePath := auctionApiPath + "/update/status?auctionId=" + strconv.Itoa(auctionId)
// 	adminSessionCookie := fmt.Sprintf(`s.id=%s`, adminSessionId)

// 	var updateStatusBody = map[string]interface{}{"Status": entity.InProgress}

// 	Test(
// 		t,
// 		Description("update auction status; success"),
// 		Post(routePath),
// 		Send().Body().JSON(updateStatusBody),
// 		Send().Headers("Cookie").Add(adminSessionCookie),
// 		Expect().Status().Equal(http.StatusOK),
// 		Expect().Custom(func(hit Hit) error {
// 			var auction entity.Auction
// 			err := hit.Response().Body().JSON().Decode(&auction)
// 			if err != nil {
// 				return err
// 			}

// 			if auction.Status != entity.InProgress {
// 				return fmt.Errorf("auction should have been update to status %s", entity.InProgress)
// 			}

// 			return nil
// 		}),
// 	)
// }

// // HTTP GET: /auction/status?status
// func TestHTTPGetPaginatedAuctionsByStatus(t *testing.T) {
// 	baseRoutePath := auctionApiPath + "/status"
// 	adminSessionCookie := fmt.Sprintf(`s.id=%s`, adminSessionId)

// 	limit10 := 10
// 	page1 := 1

// 	routePathIncorrectStatus := fmt.Sprintf("%s?status=%s&limit=%d&page=%d&asc=true", baseRoutePath, "random", limit10, page1)

// 	Test(t,
// 		Description("get auctions by status; incorrect status; failure"),
// 		Get(routePathIncorrectStatus),
// 		Send().Headers("Cookie").Add(adminSessionCookie),
// 		Expect().Body().String().Contains("Validation error"),
// 	)

// 	routePathAscendingOrderOpenAuctions := fmt.Sprintf("%s?status=%s&limit=%d&page=%d&asc=true", baseRoutePath, "open", limit10, page1)

// 	Test(t,
// 		Description("get auctions by status open; success; ascending order"),
// 		Get(routePathAscendingOrderOpenAuctions),
// 		Send().Headers("Cookie").Add(adminSessionCookie),
// 		Expect().Status().Equal(http.StatusOK),
// 		Expect().Custom(func(hit Hit) error {
// 			var auctions []entity.Auction
// 			err := hit.Response().Body().JSON().Decode(&auctions)
// 			if err != nil {
// 				return err
// 			}

// 			if len(auctions) > limit10 {
// 				return fmt.Errorf("auctions should be less than %d", limit10)
// 			}

// 			for _, r := range auctions {
// 				if r.Status != entity.Open {
// 					return fmt.Errorf("auction with id: %d is not open", r.Id)
// 				}
// 			}

// 			isAscendingOrder := sort.SliceIsSorted(auctions, func(p, q int) bool {
// 				return auctions[p].Deadline < auctions[q].Deadline
// 			})

// 			if !isAscendingOrder {
// 				return errors.New("auctions are not in ascending order")
// 			}

// 			return nil
// 		}),
// 	)

// 	routePathDescendingOrderOpenAuctions := fmt.Sprintf("%s?status=%s&limit=%d&page=%d&asc=false", baseRoutePath, entity.Open, limit10, page1)

// 	Test(t,
// 		Description("get auctions by status assigned; success; descending order"),
// 		Get(routePathDescendingOrderOpenAuctions),
// 		Send().Headers("Cookie").Add(adminSessionCookie),
// 		Expect().Status().Equal(http.StatusOK),
// 		Expect().Custom(func(hit Hit) error {
// 			var auctions []entity.Auction
// 			err := hit.Response().Body().JSON().Decode(&auctions)
// 			if err != nil {
// 				return err
// 			}

// 			if len(auctions) > limit10 {
// 				return fmt.Errorf("auctions should be less than %d", limit10)
// 			}

// 			for _, r := range auctions {
// 				if r.Status != entity.Open {
// 					return fmt.Errorf("auction with id: %d is not open", r.Id)
// 				}
// 			}

// 			return nil
// 		}),
// 	)
// }

// // HTTP GET: /auction/status/count
// func TestHTTPCountAllAuctionsByStatus(t *testing.T) {
// 	countBaseRoutePath := auctionApiPath + "/status/count"
// 	adminSessionCookie := fmt.Sprintf(`s.id=%s`, adminSessionId)

// 	var assignedAuctions = 1

// 	countAssignedRoutePath := fmt.Sprintf("%s?status=%s", countBaseRoutePath, entity.Assigned)

// 	Test(t,
// 		Description("count owned auctions; success"),
// 		Get(countAssignedRoutePath),
// 		Send().Headers("Cookie").Add(adminSessionCookie),
// 		Expect().Status().Equal(http.StatusOK),
// 		Expect().Custom(func(hit Hit) error {
// 			var count int
// 			err := hit.Response().Body().JSON().Decode(&count)
// 			if err != nil {
// 				return err
// 			}

// 			if count != assignedAuctions {
// 				return fmt.Errorf("auctions should be %d", assignedAuctions)
// 			}

// 			return nil
// 		}))
// }

// // HTTP GET: /own/assigned-bids
// func TestHTTPGetOwnAssignedAuctions(t *testing.T) {
// 	basePath := auctionApiPath + "/own/assigned-bids"
// 	sessionCookie := fmt.Sprintf(`s.id=%s`, sessionId)

// 	limit10 := 10
// 	page1 := 1

// 	routePathAscendingOrder := fmt.Sprintf("%s?limit=%d&page=%d&asc=true", basePath, limit10, page1)

// 	Test(t,
// 		Description("get own assigned auctions; success"),
// 		Get(routePathAscendingOrder),
// 		Send().Headers("Content-Type").Add("application/json"),
// 		Send().Headers("Cookie").Add(sessionCookie),
// 		Expect().Status().Equal(http.StatusOK),
// 		Expect().Custom(func(hit Hit) error {
// 			var auctions []entity.Auction

// 			err := hit.Response().Body().JSON().Decode(&auctions)
// 			if err != nil {
// 				return err
// 			}

// 			if len(auctions) != 1 {
// 				return fmt.Errorf("auctions should be %d", 1)
// 			}

// 			return nil
// 		}))
// }

// // HTTP GET: /own/assigned-bids/count
// func TestHTTPCountOwnAssignedAuctions(t *testing.T) {
// 	routePath := auctionApiPath + "/own/assigned-bids/count"
// 	sessionCookie := fmt.Sprintf(`s.id=%s`, sessionId)

// 	Test(t,
// 		Description("count own assigned extended auctions; success"),
// 		Get(routePath),
// 		Send().Headers("Content-Type").Add("application/json"),
// 		Send().Headers("Cookie").Add(sessionCookie),
// 		Expect().Status().Equal(http.StatusOK),
// 		Expect().Custom(func(hit Hit) error {
// 			var count int
// 			err := hit.Response().Body().JSON().Decode(&count)
// 			if err != nil {
// 				return err
// 			}

// 			if count != 1 {
// 				return fmt.Errorf("auctions should be %d", 1)
// 			}

// 			return nil
// 		}))
// }

// // HTTP POST: /auction/update/reject
// func TestHTTPRejectAuction(t *testing.T) {
// 	routePath := auctionApiPath + "/update/reject?auctionId=" + strconv.Itoa(auctionId)
// 	adminSessionCookie := fmt.Sprintf(`s.id=%s`, adminSessionId)

// 	Test(
// 		t,
// 		Description("reject auction; success"),
// 		Post(routePath),
// 		Send().Body().JSON(testdata.MockRejectionReason),
// 		Send().Headers("Cookie").Add(adminSessionCookie),
// 		Expect().Status().Equal(http.StatusOK),
// 		Expect().Custom(func(hit Hit) error {
// 			var auction entity.Auction
// 			err := hit.Response().Body().JSON().Decode(&auction)
// 			if err != nil {
// 				return err
// 			}

// 			if auction.Id != auctionId {
// 				return fmt.Errorf("incorrect auction returned")
// 			}

// 			if auction.Status != entity.Rejected {
// 				return fmt.Errorf("auction should have been update to status %s", entity.Rejected)
// 			}

// 			if auction.RejectionReason != testdata.MockRejectionReason["RejectionReason"] {
// 				return fmt.Errorf("incorrect rejection reason returned")
// 			}

// 			return nil
// 		}),
// 	)
// }

// // HTTP POST: /auction/own/rejected
// func TestHTTPGetPaginatedOwnRejectedAuctions(t *testing.T) {
// 	getOwnRoutePath := auctionApiPath + "/own/rejected"
// 	sessionCookie := fmt.Sprintf(`s.id=%s`, sessionId)

// 	limit10 := 10
// 	page1 := 1

// 	routePathAscendingOrder := fmt.Sprintf("%s?limit=%d&page=%d&asc=true", getOwnRoutePath, limit10, page1)

// 	Test(t,
// 		Description("get owned rejected auctions; success; ascending order"),
// 		Get(routePathAscendingOrder),
// 		Send().Headers("Cookie").Add(sessionCookie),
// 		Expect().Status().Equal(http.StatusOK),
// 		Expect().Custom(func(hit Hit) error {
// 			var auctions []entity.Auction
// 			err := hit.Response().Body().JSON().Decode(&auctions)
// 			if err != nil {
// 				return err
// 			}

// 			for _, auction := range auctions {
// 				if auction.CreatorId != userId {
// 					return fmt.Errorf("auctions creatorId: %s is not equal to the userId: %s", auction.CreatorId, userId)
// 				}
// 			}

// 			if len(auctions) != 1 {
// 				return fmt.Errorf("auctions should be %d", 1)
// 			}

// 			isAscendingOrder := sort.SliceIsSorted(auctions, func(p, q int) bool {
// 				return auctions[p].Deadline < auctions[q].Deadline
// 			})

// 			if !isAscendingOrder {
// 				return errors.New("auctions are not in ascending order")
// 			}

// 			return nil
// 		}),
// 	)

// 	routePathDescendingOrder := fmt.Sprintf("%s?limit=%d&page=%d&asc=false", getOwnRoutePath, limit10, page1)

// 	Test(t,
// 		Description("get owned rejected auctions; success; ascending order"),
// 		Get(routePathDescendingOrder),
// 		Send().Headers("Cookie").Add(sessionCookie),
// 		Expect().Status().Equal(http.StatusOK),
// 		Expect().Custom(func(hit Hit) error {
// 			var auctions []entity.Auction
// 			err := hit.Response().Body().JSON().Decode(&auctions)
// 			if err != nil {
// 				return err
// 			}

// 			for _, auction := range auctions {
// 				if auction.CreatorId != userId {
// 					return fmt.Errorf("auctions creatorId: %s is not equal to the userId: %s", auction.CreatorId, userId)
// 				}
// 			}

// 			if len(auctions) != 1 {
// 				return fmt.Errorf("auctions should be %d", 1)
// 			}

// 			return nil
// 		}),
// 	)
// }

// // HTTP POST: /auction/own/rejected/count
// func TestHTTPCountOwnRejectedAuctions(t *testing.T) {
// 	countOwnRoutePath := auctionApiPath + "/own/rejected/count"
// 	sessionCookie := fmt.Sprintf(`s.id=%s`, sessionId)

// 	var ownedAuctions = 1

// 	Test(t,
// 		Description("count owned rejected auctions; success"),
// 		Get(countOwnRoutePath),
// 		Send().Headers("Cookie").Add(sessionCookie),
// 		Expect().Status().Equal(http.StatusOK),
// 		Expect().Custom(func(hit Hit) error {
// 			var count int
// 			err := hit.Response().Body().JSON().Decode(&count)
// 			if err != nil {
// 				return err
// 			}

// 			if count != ownedAuctions {
// 				return fmt.Errorf("auctions should be %d", ownedAuctions)
// 			}

// 			return nil
// 		}))
// }

// // HTTP POST: /user/logout
// func TestHTTPDoLogout(t *testing.T) {
// 	routePath := userApiPath + "/logout"
// 	sessionCookie := fmt.Sprintf(`s.id=%s`, sessionId)

// 	Test(t,
// 		Description("logout; success"),
// 		Post(routePath),
// 		Send().Headers("Content-Type").Add("application/json"),
// 		Send().Headers("Cookie").Add(sessionCookie),
// 		Send().Body().JSON(testdata.MockUser),
// 		Expect().Status().Equal(http.StatusOK),
// 		Expect().Body().String().Contains("Successfully logged out"),
// 	)

// 	Test(t,
// 		Description("logout; failure; invalid session"),
// 		Post(routePath),
// 		Send().Headers("Content-Type").Add("application/json"),
// 		Send().Headers("Cookie").Add("s.id=123"),
// 		Send().Body().JSON(testdata.MockUser),
// 		Expect().Status().Equal(http.StatusBadRequest),
// 		Expect().Body().String().Contains("Invalid session token"),
// 	)
// }
