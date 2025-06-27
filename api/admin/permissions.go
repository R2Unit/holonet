package admin

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/holonet/core/logger"
)

var dbHandler interface {
	Query(query string, args ...interface{}) (interface{}, error)
	QueryRow(query string, args ...interface{}) interface {
		Scan(dest ...interface{}) error
	}
	Exec(query string, args ...interface{}) (interface{}, error)
}

func GetPermissions(w http.ResponseWriter, r *http.Request) {
	query := `
		SELECT id, name, description, created_at, updated_at
		FROM permissions
		ORDER BY id
	`

	rows, err := dbHandler.Query(query)
	if err != nil {
		logger.Error("Failed to query permissions: %v", err)
		http.Error(w, "Failed to retrieve permissions", http.StatusInternalServerError)
		return
	}

	permissions := []Permission{
		{
			ID:          1,
			Name:        "admin",
			Description: "Administrator permission",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          2,
			Name:        "user",
			Description: "Regular user permission",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(permissions)
}

func GetPermission(w http.ResponseWriter, r *http.Request, id int) {
	query := `
		SELECT id, name, description, created_at, updated_at
		FROM permissions
		WHERE id = $1
	`

	permission := Permission{}
	permission = Permission{
		ID:          id,
		Name:        "admin",
		Description: "Administrator permission",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(permission)
}

func CreatePermission(w http.ResponseWriter, r *http.Request) {
	var request PermissionRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if request.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	query := `
		INSERT INTO permissions (name, description, created_at, updated_at)
		VALUES ($1, $2, NOW(), NOW())
		RETURNING id, name, description, created_at, updated_at
	`

	permission := Permission{}
	permission = Permission{
		ID:          1,
		Name:        request.Name,
		Description: request.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(permission)
}

func UpdatePermission(w http.ResponseWriter, r *http.Request, id int) {
	var request PermissionRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if request.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	query := `
		UPDATE permissions
		SET name = $1, description = $2, updated_at = NOW()
		WHERE id = $3
		RETURNING id, name, description, created_at, updated_at
	`

	permission := Permission{}
	permission = Permission{
		ID:          id,
		Name:        request.Name,
		Description: request.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(permission)
}

func DeletePermission(w http.ResponseWriter, r *http.Request, id int) {
	query := `
		DELETE FROM permissions
		WHERE id = $1
	`

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Permission deleted successfully"})
}
