package api

import (
	"context"
	"net/http"

	"github.com/abstractmelon/is-site-live/internal/auth"
	"github.com/abstractmelon/is-site-live/internal/models"
	"github.com/gin-gonic/gin"
)

// register handles user registration
func (s *Server) register(c *gin.Context) {
	var registration models.UserRegistration

	// Bind request body
	if err := c.ShouldBindJSON(&registration); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if username already exists
	var count int
	err := s.db.Pool.QueryRow(context.Background(), `
		SELECT COUNT(*) FROM users WHERE username = $1
	`, registration.Username).Scan(&count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check username"})
		return
	}
	if count > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
		return
	}

	// Hash password
	passwordHash, err := models.HashPassword(registration.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Create user
	var user models.User
	err = s.db.Pool.QueryRow(context.Background(), `
		INSERT INTO users (username, password_hash, email)
		VALUES ($1, $2, $3)
		RETURNING id, username, email, created_at, updated_at
	`, registration.Username, passwordHash, registration.Email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Generate JWT token
	token, err := auth.GenerateToken(user.ID, user.Username, s.config.JWT)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Return user and token
	c.JSON(http.StatusCreated, gin.H{
		"user":  user.ToResponse(),
		"token": token,
	})
}

// login handles user login
func (s *Server) login(c *gin.Context) {
	var login models.UserLogin

	// Bind request body
	if err := c.ShouldBindJSON(&login); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user by username
	var user models.User
	err := s.db.Pool.QueryRow(context.Background(), `
		SELECT id, username, password_hash, email, created_at, updated_at
		FROM users
		WHERE username = $1
	`, login.Username).Scan(
		&user.ID,
		&user.Username,
		&user.PasswordHash,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// Check password
	if !models.CheckPassword(login.Password, user.PasswordHash) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// Generate JWT token
	token, err := auth.GenerateToken(user.ID, user.Username, s.config.JWT)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Return user and token
	c.JSON(http.StatusOK, gin.H{
		"user":  user.ToResponse(),
		"token": token,
	})
}

// getCurrentUser gets the current user
func (s *Server) getCurrentUser(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Get user from database
	var user models.User
	err := s.db.Pool.QueryRow(context.Background(), `
		SELECT id, username, email, created_at, updated_at
		FROM users
		WHERE id = $1
	`, userID).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
		return
	}

	// Return user
	c.JSON(http.StatusOK, user.ToResponse())
}

// updateUser updates the current user
func (s *Server) updateUser(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Bind request body
	var update struct {
		Email    *string `json:"email"`
		Password *string `json:"password"`
	}
	if err := c.ShouldBindJSON(&update); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update user
	var user models.User
	var err error

	if update.Password != nil {
		// Hash new password
		passwordHash, err := models.HashPassword(*update.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}

		// Update user with new password
		err = s.db.Pool.QueryRow(context.Background(), `
			UPDATE users
			SET password_hash = $1, email = COALESCE($2, email), updated_at = NOW()
			WHERE id = $3
			RETURNING id, username, email, created_at, updated_at
		`, passwordHash, update.Email, userID).Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
	} else {
		// Update user without changing password
		err = s.db.Pool.QueryRow(context.Background(), `
			UPDATE users
			SET email = COALESCE($1, email), updated_at = NOW()
			WHERE id = $2
			RETURNING id, username, email, created_at, updated_at
		`, update.Email, userID).Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	// Return updated user
	c.JSON(http.StatusOK, user.ToResponse())
}

// getUserProfile gets a user's public profile
func (s *Server) getUserProfile(c *gin.Context) {
	// Get username from URL
	username := c.Param("username")

	// Get user from database
	var user models.User
	err := s.db.Pool.QueryRow(context.Background(), `
		SELECT id, username, created_at
		FROM users
		WHERE username = $1
	`, username).Scan(
		&user.ID,
		&user.Username,
		&user.CreatedAt,
	)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Get user's sites
	rows, err := s.db.Pool.Query(context.Background(), `
		SELECT id, user_id, name, url, created_at, updated_at
		FROM sites
		WHERE user_id = $1
	`, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user's sites"})
		return
	}
	defer rows.Close()

	var sites []models.Site
	for rows.Next() {
		var site models.Site
		err := rows.Scan(
			&site.ID,
			&site.UserID,
			&site.Name,
			&site.URL,
			&site.CreatedAt,
			&site.UpdatedAt,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan site"})
			return
		}
		sites = append(sites, site)
	}

	// Return user profile
	c.JSON(http.StatusOK, gin.H{
		"user":  user.ToResponse(),
		"sites": sites,
	})
}
