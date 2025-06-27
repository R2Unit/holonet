package users

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/holonet/core/logger"
)

func UpdateUser(w http.ResponseWriter, r *http.Request, id int) {
	var request UserRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var exists bool
	err := dbHandler.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE id = $1 AND deleted_at IS NULL)", id).Scan(&exists)
	if err != nil {
		logger.Error("Failed to check if user exists: %v", err)
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}
	if !exists {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	query := "UPDATE users SET updated_at = NOW()"
	args := []interface{}{}
	argIndex := 1

	if request.Username != "" {
		query += ", username = $" + strconv.Itoa(argIndex)
		args = append(args, request.Username)
		argIndex++
	}

	if request.Email != "" {
		query += ", email = $" + strconv.Itoa(argIndex)
		args = append(args, request.Email)
		argIndex++
	}

	if request.Password != "" {
		passwordHash := "hashed_" + request.Password
		query += ", password_hash = $" + strconv.Itoa(argIndex)
		args = append(args, passwordHash)
		argIndex++
	}

	if request.Status != "" {
		query += ", status = $" + strconv.Itoa(argIndex)
		args = append(args, request.Status)
		argIndex++
	}

	query += " WHERE id = $" + strconv.Itoa(argIndex)
	args = append(args, id)

	_, err = dbHandler.Exec(query, args...)
	if err != nil {
		logger.Error("Failed to update user: %v", err)
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	GetUser(w, r, id)
}
