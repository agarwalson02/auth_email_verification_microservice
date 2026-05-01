package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthHeaderMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		token := c.GetHeader("Authorization")

		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "missing authorization header",
			})
			return
		}

		// store token in context for later use
		c.Set("auth", token)

		c.Next()
	}
}