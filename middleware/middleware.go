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

		// Set the username in the context for further processing
		c.Set("username", username)

		// Continue processing the request
		c.Next()
	}
}
