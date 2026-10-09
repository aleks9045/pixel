package middleware

import (
	"net/http"

	"github.com/aleks9045/pixels/internal/config"
	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", config.Load().HomepageURL)

		c.Writer.Header().
			Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PATCH, DELETE")

		c.Writer.Header().
			Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")

		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		// Handle preflight OPTIONS request.
		if c.Request.Method == "OPTIONS" {
			// Cache preflight response for 24 hours.
			c.Writer.Header().Set("Access-Control-Max-Age", "86400")

			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		// Pass control to the next middleware/handler.
		c.Next()
	}
}
