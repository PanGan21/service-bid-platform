package http

import (
	"net/http"
	"time"

	bidController "github.com/PanGan21/bidding-service/internal/routes/http/bid"
	"github.com/PanGan21/bidding-service/internal/service"
	"github.com/PanGan21/pkg/auth"
	"github.com/PanGan21/pkg/logger"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewRouter(handler *gin.Engine, l logger.Interface, authService auth.AuthService, bidService service.BidService) {
	bidController := bidController.NewBidController(l, bidService)
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
	var requiredRoles []string
	handler.POST("/", auth.AuthorizeEndpoint(requiredRoles...), bidController.Create)
	handler.GET("/", bidController.GetOneById)
	handler.GET("/requestId/", bidController.GetManyByRequestId)
}
