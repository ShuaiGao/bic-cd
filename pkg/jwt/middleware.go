package jwt

import (
	"bic-cd/internal/util"
	"github.com/gin-gonic/gin"
	"net/http"
)

// JWT is jwt middleware
func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		payload, err := parseToken(token)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.Set(util.Username, payload.Username)
		c.Set(util.UserID, payload.UserId)
		c.Next()
	}
}
