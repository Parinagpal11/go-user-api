package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/yourusername/go-user-api/internal/models"
	"github.com/yourusername/go-user-api/internal/utils"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

// Register creates a new user account
func Register(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.RegisterRequest

		// Parse request body
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid request body"})
			return
		}

		// Validate input
		if err := req.Validate(); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
			return
		}

		// Hash password
		hashedPassword, err := utils.HashPassword(req.Password)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to process password"})
			return
		}

		// Insert user into database
		var user models.User
		query := `
			INSERT INTO users (email, username, password_hash, first_name, last_name)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING id, email, username, first_name, last_name, created_at, updated_at
		`
		err = db.QueryRow(
			query,
			req.Email,
			req.Username,
			hashedPassword,
			req.FirstName,
			req.LastName,
		).Scan(
			&user.ID,
			&user.Email,
			&user.Username,
			&user.FirstName,
			&user.LastName,
			&user.CreatedAt,
			&user.UpdatedAt,
		)

		if err != nil {
			// Check for duplicate email/username
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Email or username already exists"})
			return
		}

		// Generate JWT token
		token, err := utils.GenerateToken(user.ID)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to generate token"})
			return
		}

		// Return success response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(models.LoginResponse{
			Token: token,
			User:  user,
		})
	}
}

// Login authenticates a user and returns a JWT token
func Login(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.LoginRequest

		// Parse request body
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid request body"})
			return
		}

		// Validate input
		if req.Email == "" || req.Password == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Email and password are required"})
			return
		}

		// Find user by email
		var user models.User
		query := `
			SELECT id, email, username, password_hash, first_name, last_name, created_at, updated_at
			FROM users
			WHERE email = $1
		`
		err := db.QueryRow(query, req.Email).Scan(
			&user.ID,
			&user.Email,
			&user.Username,
			&user.PasswordHash,
			&user.FirstName,
			&user.LastName,
			&user.CreatedAt,
			&user.UpdatedAt,
		)

		if err == sql.ErrNoRows {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid email or password"})
			return
		} else if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Database error"})
			return
		}

		// Check password
		if !utils.CheckPassword(user.PasswordHash, req.Password) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid email or password"})
			return
		}

		// Generate JWT token
		token, err := utils.GenerateToken(user.ID)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to generate token"})
			return
		}

		// Return success response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(models.LoginResponse{
			Token: token,
			User:  user,
		})
	}
}
