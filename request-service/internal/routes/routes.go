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

	handler.Use(auth.VerifyJWT(authService))
	// K8s probe

	var emptyRoles []string
	handler.GET("/healthz", auth.AuthorizeEndpoint(emptyRoles...), func(c *gin.Context) { c.Status(http.StatusOK) })
	// handler.Use(auth.VerifyJWT(authService))
}
