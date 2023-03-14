package http

import (
	"net/http"
	"time"

	"github.com/PanGan21/pkg/auth"
	"github.com/PanGan21/pkg/logger"
	userController "github.com/PanGan21/user-service/internal/routes/http/user"
	"github.com/PanGan21/user-service/internal/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

func NewRouter(handler *gin.Engine, l logger.Interface, store sessions.RedisStore, userService service.UserService, authService auth.AuthService) {
	userController := userController.NewUserController(l, userService, authService)
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

	// Session
	handler.Use(sessions.Sessions("s.id", store))

	// Routers
	handler.GET("/", userController.GetLoggedInUserDetails)
	handler.POST("/login", userController.Login)
	handler.POST("/logout", userController.Logout)
	handler.POST("/register", userController.Register)
	handler.GET("/authenticate", userController.Authenticate)
	handler.GET("/details", userController.GetDetailsById)
}
