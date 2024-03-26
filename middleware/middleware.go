package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AnonymousLoginMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.GetHeader("Username")
		if username == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing username in header"})
			c.Abort()
			return
		}

		
		c.Set("username", username)
		c.Next()
	}
}
