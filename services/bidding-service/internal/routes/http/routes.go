package http

import (
	"net/http"

	bidController "github.com/PanGan21/bidding-service/internal/routes/http/bid"
	"github.com/PanGan21/bidding-service/internal/service"
	"github.com/PanGan21/pkg/auth"
	"github.com/PanGan21/pkg/logger"
	"github.com/gin-gonic/gin"
)

func NewRouter(handler *gin.Engine, l logger.Interface, authService auth.AuthService, bidService service.BidService) {
	bidController := bidController.NewBidController(l, bidService)
	// Options
	handler.Use(gin.Recovery())

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
