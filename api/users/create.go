package users

import (
	"encoding/json"
	"net/http"

	"github.com/holonet/core/logger"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var request UserRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if request.Username == "" || request.Email == "" || request.Password == "" {
		http.Error(w, "Username, email, and password are required", http.StatusBadRequest)
		return
	}

	passwordHash := "hashed_" + request.Password

	status := "active"
	if request.Status != "" {
		status = request.Status
	}

	query := `
		INSERT INTO users (username, email, password_hash, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		RETURNING id, username, email, status, created_at, updated_at
	`

	user := User{}
	err := dbHandler.QueryRow(
		query,
		request.Username,
		request.Email,
		passwordHash,
		status,
	).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Status,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		logger.Error("Failed to create user: %v", err)
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}
