package http

import (
	"net/http"

	"github.com/PanGan21/pkg/auth"
	"github.com/PanGan21/pkg/logger"
	requestController "github.com/PanGan21/request-service/internal/routes/request"
	"github.com/PanGan21/request-service/internal/service"
	"github.com/gin-gonic/gin"
)

func NewRouter(handler *gin.Engine, l logger.Interface, authService auth.AuthService, requestService service.RequestService) {
	requestController := requestController.NewRequestController(l, requestService)
	// Options
	handler.Use(gin.Recovery())

	// K8s probe
	handler.GET("/healthz", func(c *gin.Context) { c.Status(http.StatusOK) })

	// JWT Middleware
	handler.Use(auth.VerifyJWT(authService))

	// Routers
	handler.GET("/", requestController.GetAll)
	handler.POST("/", requestController.Create)

	var requiredRoles []string
	handler.GET("/hello", auth.AuthorizeEndpoint(requiredRoles...), func(c *gin.Context) { c.Status(http.StatusOK) })
}
