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

func NewRouter(handler *gin.Engine, l logger.Interface, corsOrigins []string, authService auth.AuthService, requestService service.RequestService) {
	requestController := requestController.NewRequestController(l, requestService)
	// Options
	handler.Use(gin.Recovery())

	// Cors
	handler.Use(cors.New(cors.Config{
		AllowOrigins: corsOrigins,
		AllowMethods: []string{"POST", "GET", "OPTIONS"},
		AllowHeaders: []string{"DNT", "X-CustomHeader", "Keep-Alive", "User-Agent", "X-Requested-With", "If-Modified-Since", "Cache-Control", "Content-Type"},
		MaxAge:       12 * time.Hour,
	}))

	// K8s probe
	handler.GET("/healthz", func(c *gin.Context) { c.Status(http.StatusOK) })

	// JWT Middleware
	handler.Use(auth.VerifyJWT(authService))

	// Routers
	handler.POST("/", requestController.Create)
	handler.GET("/status", requestController.GetByStatus)
	handler.GET("/status/count", requestController.CountByStatus)
	handler.GET("/status/own", requestController.GetOwnByStatus)
	handler.GET("/status/own/count", requestController.CountOwnByStatus)

	requireAdminRole := []string{"ADMIN"}
	handler.POST("/reject", auth.AuthorizeEndpoint(requireAdminRole...), requestController.RejectRequest)
	handler.POST("/approve", auth.AuthorizeEndpoint(requireAdminRole...), requestController.Approve)
}
