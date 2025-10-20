package middleware

import (
	"os"

	"github.com/gin-gonic/gin"
)

// backend/internal/adapters/http/middleware/cors.go - УЛУЧШИТЬ:
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// В production использовать конкретные домены
		allowedOrigins := []string{
			"http://localhost:3000",
			"http://localhost:8080",
			// Добавить production домены
		}

		isAllowed := false
		for _, allowed := range allowedOrigins {
			if origin == allowed {
				isAllowed = true
				break
			}
		}

		if isAllowed || os.Getenv("ENVIRONMENT") == "development" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		}

		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
