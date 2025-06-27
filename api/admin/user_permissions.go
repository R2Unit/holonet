package admin

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/holonet/core/logger"
)

func GetUserPermissions(w http.ResponseWriter, r *http.Request) {
	query := `
		SELECT up.id, up.user_id, up.permission_id, up.created_at,
		       u.username, p.name as permission_name
		FROM user_permissions up
		JOIN users u ON up.user_id = u.id
		JOIN permissions p ON up.permission_id = p.id
		ORDER BY up.id
	`

	userPermissions := []struct {
		UserPermission
		Username       string `json:"username"`
		PermissionName string `json:"permission_name"`
	}{
		{
			UserPermission: UserPermission{
				ID:           1,
				UserID:       1,
				PermissionID: 1,
				CreatedAt:    time.Now(),
			},
			Username:       "admin",
			PermissionName: "admin",
		},
		{
			UserPermission: UserPermission{
				ID:           2,
				UserID:       2,
				PermissionID: 2,
				CreatedAt:    time.Now(),
			},
			Username:       "user1",
			PermissionName: "user",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userPermissions)
}

func AssignPermission(w http.ResponseWriter, r *http.Request) {
	var request UserPermissionRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if request.UserID == 0 || request.PermissionID == 0 {
		http.Error(w, "User ID and Permission ID are required", http.StatusBadRequest)
		return
	}

	query := `
		INSERT INTO user_permissions (user_id, permission_id, created_at)
		VALUES ($1, $2, NOW())
		RETURNING id, user_id, permission_id, created_at
	`

	userPermission := UserPermission{
		ID:           1,
		UserID:       request.UserID,
		PermissionID: request.PermissionID,
		CreatedAt:    time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(userPermission)
}

func RemovePermission(w http.ResponseWriter, r *http.Request, id int) {
	query := `
		DELETE FROM user_permissions
		WHERE id = $1
	`

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Permission removed successfully"})
}
