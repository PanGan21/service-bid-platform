package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

const userKey = "userId"

func AuthRequired(c *gin.Context) {
	session := sessions.Default(c)
	sessionID := session.Get(userKey)
	fmt.Println("sessionID", sessionID)
	if sessionID == nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
}
