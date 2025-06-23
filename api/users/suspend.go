package users

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/holonet/core/logger"
)

// SuspendUser suspends a user by setting their status to "suspended"
func SuspendUser(w http.ResponseWriter, r *http.Request, id int) {
	// Check if user exists
	var exists bool
	err := dbHandler.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE id = $1 AND deleted_at IS NULL)", id).Scan(&exists)
	if err != nil {
		logger.Error("Failed to check if user exists: %v", err)
		http.Error(w, "Failed to suspend user", http.StatusInternalServerError)
		return
	}
	if !exists {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Update user status to suspended
	query := `
		UPDATE users
		SET status = 'suspended', updated_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
		RETURNING id
	`

	var returnedID int
	err = dbHandler.QueryRow(query, id).Scan(&returnedID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		logger.Error("Failed to suspend user: %v", err)
		http.Error(w, "Failed to suspend user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "User suspended successfully",
	})
}

// UnsuspendUser reactivates a suspended user by setting their status to "active"
func UnsuspendUser(w http.ResponseWriter, r *http.Request, id int) {
	// Check if user exists
	var exists bool
	err := dbHandler.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE id = $1 AND deleted_at IS NULL)", id).Scan(&exists)
	if err != nil {
		logger.Error("Failed to check if user exists: %v", err)
		http.Error(w, "Failed to unsuspend user", http.StatusInternalServerError)
		return
	}
	if !exists {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Update user status to active
	query := `
		UPDATE users
		SET status = 'active', updated_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
		RETURNING id
	`

	var returnedID int
	err = dbHandler.QueryRow(query, id).Scan(&returnedID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		logger.Error("Failed to unsuspend user: %v", err)
		http.Error(w, "Failed to unsuspend user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "User unsuspended successfully",
	})
}
