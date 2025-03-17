package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

// User represents a user in the system
type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"`
	Email        string    `json:"email,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// UserRegistration represents the data needed to register a new user
type UserRegistration struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=8"`
	Email    string `json:"email"`
}

// UserLogin represents the data needed to log in a user
type UserLogin struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// UserResponse represents the user data returned to the client
type UserResponse struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// HashPassword hashes a password using bcrypt
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPassword checks if the provided password matches the stored hash
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// ToResponse converts a User to a UserResponse
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
	}
}
