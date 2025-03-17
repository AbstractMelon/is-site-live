package auth

import (
	"net/http"
	"strings"

	"github.com/abstractmelon/is-site-live/internal/config"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware creates a middleware for JWT authentication
func AuthMiddleware(cfg config.JWTConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header is required"})
			c.Abort()
			return
		}

		// Check if authorization header is in the correct format
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header format must be Bearer {token}"})
			c.Abort()
			return
		}

		// Get token
		tokenString := parts[1]

		// Validate token
		claims, err := ValidateToken(tokenString, cfg)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			c.Abort()
			return
		}

		// Set user ID and username in context
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)

		c.Next()
	}
}
