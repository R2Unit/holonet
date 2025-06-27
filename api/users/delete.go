package users

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/holonet/core/logger"
)

func DeleteUser(w http.ResponseWriter, r *http.Request, id int) {
	query := `
		UPDATE users
		SET deleted_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
		RETURNING id
	`

	var returnedID int
	err := dbHandler.QueryRow(query, id).Scan(&returnedID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		logger.Error("Failed to delete user: %v", err)
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "User deleted successfully",
	})
}
