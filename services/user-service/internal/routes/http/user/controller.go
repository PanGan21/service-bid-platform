package user

import (
	"net/http"
	"strconv"

	"github.com/PanGan21/pkg/auth"
	"github.com/PanGan21/pkg/logger"
	"github.com/PanGan21/user-service/internal/service"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

type UserController interface {
	Login(c *gin.Context)
	Logout(c *gin.Context)
	Register(c *gin.Context)
	GetUserDetails(c *gin.Context)
	Authenticate(c *gin.Context)
}

type userController struct {
	logger      logger.Interface
	userService service.UserService
	authService auth.AuthService
}

const userKey = "userId"

func NewUserController(logger logger.Interface, userServ service.UserService, authServ auth.AuthService) UserController {
	return &userController{
		logger:      logger,
		userService: userServ,
		authService: authServ,
	}
}

type UserData struct {
	Username string `json:"Username"`
	Password string `json:"Password"`
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
	userId := session.Get(userKey)
	if userId == nil {
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
	userId := session.Get(userKey)
	if userId != nil {
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

func (controller *userController) Authenticate(c *gin.Context) {
	session := sessions.Default(c)
	userId := session.Get(userKey)

	id, ok := userId.(string)
	if userId == nil || !ok {
		controller.logger.Error("Invalid session token")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid session token"})
		return
	}

	parsedId, err := strconv.Atoi(id)
	if err != nil {
		controller.logger.Error("Cannot convert id internally")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid session token"})
		return
	}

	// Find user
	user, err := controller.userService.GetById(c.Request.Context(), parsedId)
	if err != nil {
		controller.logger.Error(err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "aunauthorized"})
		return
	}

	methodHeader := c.Request.Header.Get("x-forwarded-method")
	uriHeader := c.Request.Header.Get("x-forwarded-uri")
	if uriHeader == "" {
		uriHeader = "/user/authenticate"
	}

	var method = methodHeader
	if method == "OPTIONS" {
		c.JSON(http.StatusOK, gin.H{})
		return
	}

	if method == "" {
		method = "GET"
	}

	token, err := controller.authService.SignJWT(userId.(string), strconv.Itoa(user.Id), uriHeader, user.Roles...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "JWT signing error"})
	}

	c.Header("x-internal-jwt", token)

	c.JSON(http.StatusOK, gin.H{"message": "Successfully authenticated user"})
}

type UserDetails struct {
	Id       string   `json:"Id"`
	Username string   `json:"Username"`
	Roles    []string `json:"Roles"`
}

func (controller *userController) GetUserDetails(c *gin.Context) {
	session := sessions.Default(c)
	userId := session.Get(userKey)
	id, ok := userId.(string)
	if userId == nil || !ok {
		controller.logger.Error("Invalid session token")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid session token"})
		return
	}

	parsedId, err := strconv.Atoi(id)
	if err != nil {
		controller.logger.Error("Cannot convert id internally")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid session token"})
		return
	}

	// Find user
	user, err := controller.userService.GetById(c.Request.Context(), parsedId)
	if err != nil {
		controller.logger.Error(err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "aunauthorized"})
		return
	}

	userDetails := &UserDetails{
		Id:       strconv.Itoa(user.Id),
		Username: user.Username,
		Roles:    user.Roles,
	}

	c.JSON(http.StatusOK, userDetails)
}
