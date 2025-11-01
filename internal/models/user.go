package models

import (
	"errors"
	"regexp"
	"time"
)

type User struct {
	ID           int       `json:"id"`
	Email        string    `json:"email"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"` // Never send password hash to client
	FirstName    string    `json:"first_name,omitempty"`
	LastName     string    `json:"last_name,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type RegisterRequest struct {
	Email     string `json:"email"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type UpdateUserRequest struct {
	Email     string `json:"email,omitempty"`
	Username  string `json:"username,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
}

// Validate checks if user data is valid
func (r *RegisterRequest) Validate() error {
	if r.Email == "" {
		return errors.New("email is required")
	}
	if !isValidEmail(r.Email) {
		return errors.New("invalid email format")
	}
	if r.Username == "" {
		return errors.New("username is required")
	}
	if len(r.Username) < 3 {
		return errors.New("username must be at least 3 characters")
	}
	if r.Password == "" {
		return errors.New("password is required")
	}
	if len(r.Password) < 6 {
		return errors.New("password must be at least 6 characters")
	}
	return nil
}

func isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}
