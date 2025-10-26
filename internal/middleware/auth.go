package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(secretKey string, tokenRequired bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip auth if not required
		if !tokenRequired {
			c.Next()
			return
		}

		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized - Missing token",
			})
			c.Abort()
			return
		}

		// Extract Bearer token
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized - Invalid token format",
			})
			c.Abort()
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")

		// Validate token (simple validation for demo)
		if !isValidToken(token, secretKey) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized - Invalid token",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

func isValidToken(token, secretKey string) bool {
	// Simple token validation - in production, use JWT or other secure methods
	// For demo purposes, we accept any non-empty token
	if token == "" {
		return false
	}

	// You can implement more sophisticated validation here:
	// - JWT token validation
	// - Database lookup
	// - External auth service validation
	
	return len(token) >= 8 // Simple length check for demo
}