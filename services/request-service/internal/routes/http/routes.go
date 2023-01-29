package http

import (
	"net/http"
	"time"

	"github.com/PanGan21/pkg/auth"
	"github.com/PanGan21/pkg/logger"
	requestController "github.com/PanGan21/request-service/internal/routes/http/request"
	"github.com/PanGan21/request-service/internal/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewRouter(handler *gin.Engine, l logger.Interface, authService auth.AuthService, requestService service.RequestService, bidService service.BidService) {
	requestController := requestController.NewRequestController(l, requestService, bidService)
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
	handler.GET("/", requestController.GetAll)
	handler.GET("/count", requestController.CountAll)
	handler.POST("/", requestController.Create)
	handler.GET("/count/own", requestController.CountOwn)
	handler.GET("/own", requestController.GetOwn)

	updateRequestRoles := []string{"ADMIN"}
	handler.POST("/update/winner", auth.AuthorizeEndpoint(updateRequestRoles...), requestController.UpdateWinnerByRequestId)

	var requiredRoles []string
	handler.GET("/hello", auth.AuthorizeEndpoint(requiredRoles...), func(c *gin.Context) { c.Status(http.StatusOK) })
}
