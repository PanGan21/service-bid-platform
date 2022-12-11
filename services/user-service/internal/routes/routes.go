package http

import (
	"net/http"

	"github.com/PanGan21/pkg/auth"
	"github.com/PanGan21/pkg/logger"
	userController "github.com/PanGan21/user-service/internal/routes/user"
	"github.com/PanGan21/user-service/internal/service"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

func NewRouter(handler *gin.Engine, l logger.Interface, store sessions.RedisStore, userService service.UserService, authService auth.AuthService) {
	userController := userController.NewUserController(l, userService, authService)
	// Options
	handler.Use(gin.Recovery())

	// K8s probe
	handler.GET("/healthz", func(c *gin.Context) { c.Status(http.StatusOK) })

	// Session
	handler.Use(sessions.Sessions("s.id", store))

	// Routers
	handler.GET("/", userController.GetUserDetails)
	handler.POST("/login", userController.Login)
	handler.POST("/logout", userController.Logout)
	handler.POST("/register", userController.Register)
	handler.GET("/authenticate", userController.Authenticate)
}