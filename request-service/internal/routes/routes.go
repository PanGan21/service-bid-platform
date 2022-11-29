package http

import (
	"net/http"

	"github.com/PanGan21/packages/auth"
	"github.com/PanGan21/packages/logger"
	"github.com/gin-gonic/gin"
)

func NewRouter(handler *gin.Engine, l logger.Interface, authService auth.AuthService) {
	// Options
	handler.Use(gin.Recovery())

	// K8s probe
	handler.GET("/healthz", func(c *gin.Context) { c.Status(http.StatusOK) })

	// JWT Middleware
	handler.Use(auth.VerifyJWT(authService))

	var requiredRoles []string
	handler.GET("/hello", auth.AuthorizeEndpoint(requiredRoles...), func(c *gin.Context) { c.Status(http.StatusOK) })
}
