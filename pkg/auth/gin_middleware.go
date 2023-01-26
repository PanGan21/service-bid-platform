package auth

import (
	"net/http"

	"github.com/PanGan21/pkg/utils"
	"github.com/gin-gonic/gin"
)

func VerifyJWT(authService AuthService) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		jwt := c.Request.Header.Get("x-internal-jwt")
		if jwt == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "no token found"})
			return
		}

		authTokenData, err := authService.VerifyJWT(jwt, c.Request.URL.Path)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "jwt not verified"})
			return
		}

		c.Set("url", authTokenData.Route)
		c.Set("jwt", jwt)
		c.Set("userId", authTokenData.UserId)
		c.Set("roles", authTokenData.Roles)

		c.Next()
	}

	return gin.HandlerFunc(fn)
}

func AuthorizeEndpoint(allowedRoles ...string) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		requestRoles := c.Copy().Value("roles").([]string)

		if !utils.Subslice(allowedRoles, requestRoles) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "incorrect permissions"})
			return
		}
	}

	return gin.HandlerFunc(fn)
}
