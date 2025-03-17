package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/abstractmelon/is-site-live/internal/config"
	"github.com/golang-jwt/jwt/v5"
)

// Claims represents the JWT claims
type Claims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// GenerateToken generates a JWT token for a user
func GenerateToken(userID int, username string, cfg config.JWTConfig) (string, error) {
	// Set expiration time to 24 hours
	expirationTime := time.Now().Add(24 * time.Hour)

	// Create claims
	claims := &Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   fmt.Sprintf("%d", userID),
		},
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token with secret key
	tokenString, err := token.SignedString([]byte(cfg.Secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken validates a JWT token
func ValidateToken(tokenString string, cfg config.JWTConfig) (*Claims, error) {
	// Parse token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(cfg.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	// Validate token
	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	// Get claims
	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}
