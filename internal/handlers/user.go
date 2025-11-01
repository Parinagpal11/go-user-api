package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/yourusername/go-user-api/internal/models"
)

// GetUsers returns all users
func GetUsers(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := `
			SELECT id, email, username, first_name, last_name, created_at, updated_at
			FROM users
			ORDER BY created_at DESC
		`
		rows, err := db.Query(query)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to fetch users"})
			return
		}
		defer rows.Close()

		users := []models.User{}
		for rows.Next() {
			var user models.User
			err := rows.Scan(
				&user.ID,
				&user.Email,
				&user.Username,
				&user.FirstName,
				&user.LastName,
				&user.CreatedAt,
				&user.UpdatedAt,
			)
			if err != nil {
				continue
			}
			users = append(users, user)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users)
	}
}

// GetUser returns a single user by ID
func GetUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid user ID"})
			return
		}

		var user models.User
		query := `
			SELECT id, email, username, first_name, last_name, created_at, updated_at
			FROM users
			WHERE id = $1
		`
		err = db.QueryRow(query, id).Scan(
			&user.ID,
			&user.Email,
			&user.Username,
			&user.FirstName,
			&user.LastName,
			&user.CreatedAt,
			&user.UpdatedAt,
		)

		if err == sql.ErrNoRows {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "User not found"})
			return
		} else if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Database error"})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
	}
}

// GetCurrentUser returns the authenticated user's information
func GetCurrentUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get user ID from context (set by auth middleware)
		userID, ok := r.Context().Value("userID").(int)
		if !ok {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Unauthorized"})
			return
		}

		var user models.User
		query := `
			SELECT id, email, username, first_name, last_name, created_at, updated_at
			FROM users
			WHERE id = $1
		`
		err := db.QueryRow(query, userID).Scan(
			&user.ID,
			&user.Email,
			&user.Username,
			&user.FirstName,
			&user.LastName,
			&user.CreatedAt,
			&user.UpdatedAt,
		)

		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to fetch user"})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
	}
}

// UpdateUser updates a user's information
func UpdateUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid user ID"})
			return
		}

		// Get authenticated user ID
		authUserID, ok := r.Context().Value("userID").(int)
		if !ok || authUserID != id {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "You can only update your own profile"})
			return
		}

		var req models.UpdateUserRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid request body"})
			return
		}

		// Build update query dynamically based on provided fields
		query := `
			UPDATE users
			SET first_name = COALESCE(NULLIF($1, ''), first_name),
			    last_name = COALESCE(NULLIF($2, ''), last_name),
			    updated_at = NOW()
			WHERE id = $3
			RETURNING id, email, username, first_name, last_name, created_at, updated_at
		`

		var user models.User
		err = db.QueryRow(query, req.FirstName, req.LastName, id).Scan(
			&user.ID,
			&user.Email,
			&user.Username,
			&user.FirstName,
			&user.LastName,
			&user.CreatedAt,
			&user.UpdatedAt,
		)

		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to update user"})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
	}
}

// DeleteUser deletes a user
func DeleteUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid user ID"})
			return
		}

		// Get authenticated user ID
		authUserID, ok := r.Context().Value("userID").(int)
		if !ok || authUserID != id {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "You can only delete your own account"})
			return
		}

		query := `DELETE FROM users WHERE id = $1`
		result, err := db.Exec(query, id)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to delete user"})
			return
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "User not found"})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "User deleted successfully"})
	}
}
