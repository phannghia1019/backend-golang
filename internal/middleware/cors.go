package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// CORSMiddleware handles CORS headers
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if origin == "" {
			origin = "*"
		}

		c.Header("Access-Control-Allow-Origin", origin)
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, x-requested-with, Content-Length, Accept, Origin, X-Requested-With")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Content-Type, Content-Disposition, Access-Control-Allow-Origin")
		c.Header("Access-Control-Allow-Credentials", "true")

		// Handle preflight requests
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// LanguageMiddleware handles Accept-Language header
func LanguageMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		language := c.GetHeader("Accept-Language")
		if language != "vi-VN" && language != "en-US" {
			language = "vi-VN"
		}
		c.Set("language", language)
		c.Next()
	}
}
