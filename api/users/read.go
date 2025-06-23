package users

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/holonet/core/logger"
)

var dbHandler *sql.DB

// SetDBHandler sets the database handler for the users package
func SetDBHandler(db *sql.DB) {
	dbHandler = db
}

// GetUsers returns a list of all users
func GetUsers(w http.ResponseWriter, r *http.Request) {
	logger.Info("GetUsers function called")

	query := `
		SELECT id, username, email, last_login, created_at, updated_at, status
		FROM users
		WHERE deleted_at IS NULL
		ORDER BY id
	`

	rows, err := dbHandler.Query(query)
	if err != nil {
		logger.Error("Failed to list users: %v", err)
		http.Error(w, "Failed to list users", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	users := []User{}
	for rows.Next() {
		user := User{}
		var lastLogin sql.NullTime
		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&lastLogin,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.Status,
		)
		if err != nil {
			logger.Error("Failed to scan user: %v", err)
			http.Error(w, "Failed to list users", http.StatusInternalServerError)
			return
		}
		if lastLogin.Valid {
			user.LastLogin = lastLogin.Time
		}
		users = append(users, user)
	}

	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")
	logger.Info("Content-Type header set to application/json")

	// Return a JSON array with all users
	err = json.NewEncoder(w).Encode(users)
	if err != nil {
		logger.Error("Failed to encode users to JSON: %v", err)
		http.Error(w, "Failed to encode users to JSON", http.StatusInternalServerError)
		return
	}
	logger.Info("Successfully encoded users to JSON")
}

// GetUser returns a specific user by ID
func GetUser(w http.ResponseWriter, r *http.Request, id int) {
	query := `
		SELECT id, username, email, last_login, created_at, updated_at, status
		FROM users
		WHERE id = $1 AND deleted_at IS NULL
	`

	user := User{}
	var lastLogin sql.NullTime
	err := dbHandler.QueryRow(query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&lastLogin,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Status,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		logger.Error("Failed to get user: %v", err)
		http.Error(w, "Failed to get user", http.StatusInternalServerError)
		return
	}

	if lastLogin.Valid {
		user.LastLogin = lastLogin.Time
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
