package http

import (
	"net/http"
	"time"

	auctionController "github.com/PanGan21/auction-service/internal/routes/http/auction"
	"github.com/PanGan21/auction-service/internal/service"
	"github.com/PanGan21/pkg/auth"
	"github.com/PanGan21/pkg/logger"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewRouter(handler *gin.Engine, l logger.Interface, authService auth.AuthService, auctionService service.AuctionService, bidService service.BidService) {
	auctionController := auctionController.NewAuctionController(l, auctionService, bidService)
	// Options
	handler.Use(gin.Recovery())

	// Cors
	handler.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:3000"},
		AllowMethods: []string{"POST", "GET", "OPTIONS"},
		AllowHeaders: []string{"DNT", "X-CustomHeader", "Keep-Alive", "User-Agent", "X-Requested-With", "If-Modified-Since", "Cache-Control", "Content-Type"},
		MaxAge:       12 * time.Hour,
	}))

	// K8s probe
	handler.GET("/healthz", func(c *gin.Context) { c.Status(http.StatusOK) })

	// JWT Middleware
	handler.Use(auth.VerifyJWT(authService))

	// Routers
	handler.GET("/", auctionController.GetAll)
	handler.GET("/count", auctionController.CountAll)
	handler.POST("/", auctionController.Create)
	handler.GET("/count/own", auctionController.CountOwn)
	handler.GET("/own", auctionController.GetOwn)
	handler.GET("/own/assigned-bids", auctionController.GetOwnAssignedByStatuses)
	handler.GET("/own/assigned-bids/count", auctionController.CountOwnAssignedByStatuses)
	handler.GET("/status", auctionController.GetByStatus)
	handler.GET("/status/count", auctionController.CountByStatus)

	requireAdminRole := []string{"ADMIN"}
	handler.GET("/open/past-deadline", auth.AuthorizeEndpoint(requireAdminRole...), auctionController.GetOpenPastDeadline)
	handler.GET("/open/past-deadline/count", auth.AuthorizeEndpoint(requireAdminRole...), auctionController.CountOpenPastDeadline)
	handler.POST("/update/winner", auth.AuthorizeEndpoint(requireAdminRole...), auctionController.UpdateWinnerByAuctionId)
	handler.POST("/update/status", auth.AuthorizeEndpoint(requireAdminRole...), auctionController.UpdateStatus)

	var requiredRoles []string
	handler.GET("/hello", auth.AuthorizeEndpoint(requiredRoles...), func(c *gin.Context) { c.Status(http.StatusOK) })
}
