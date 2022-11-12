package user

import (
	"net/http"

	"github.com/PanGan21/packages/logger"
	"github.com/PanGan21/user-service/internal/service"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

type UserController interface {
	Login(c *gin.Context)
	Logout(c *gin.Context)
	Register(c *gin.Context)
}

type userController struct {
	logger      logger.Interface
	userService service.UserService
}

const userKey = "userId"

func NewUserController(logger logger.Interface, userServ service.UserService) UserController {
	return &userController{
		logger:      logger,
		userService: userServ,
	}
}

type UserData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (controller *userController) Login(c *gin.Context) {
	var userData UserData
	if err := c.BindJSON(&userData); err != nil {
		controller.logger.Error(err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Validation error"})
		return
	}

	userId, err := controller.userService.Login(c.Request.Context(), userData.Username, userData.Password)
	if err != nil {
		controller.logger.Error(err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
		return
	}

	session := sessions.Default(c)
	// Save the id in the session
	session.Set(userKey, userId)
	if err := session.Save(); err != nil {
		controller.logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully authenticated user"})
}

func (controller *userController) Logout(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(userKey)
	if user == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session token"})
		return
	}
	session.Delete(userKey)
	if err := session.Save(); err != nil {
		controller.logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}

func (controller *userController) Register(c *gin.Context) {
	var userData UserData
	if err := c.BindJSON(&userData); err != nil {
		controller.logger.Error(err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Validation error"})
		return
	}

	session := sessions.Default(c)
	user := session.Get(userKey)
	if user != nil {
		c.JSON(http.StatusOK, gin.H{"message": "Already logged in"})
		return
	}

	userId, err := controller.userService.Register(c.Request.Context(), userData.Username, userData.Password)
	if err != nil {
		controller.logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Registration failed"})
		return
	}

	session.Set(userKey, userId)
	if err := session.Save(); err != nil {
		controller.logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully registered user"})
}
