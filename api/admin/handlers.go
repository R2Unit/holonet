package admin

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/holonet/core/logger"
)

func HandlePermissions(w http.ResponseWriter, r *http.Request) {
	logger.Info("HandlePermissions called with method: %s", r.Method)
	switch r.Method {
	case http.MethodGet:
		logger.Info("Calling GetPermissions")
		GetPermissions(w, r)
	case http.MethodPost:
		logger.Info("Calling CreatePermission")
		CreatePermission(w, r)
	default:
		logger.Info("Method not allowed: %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func HandlePermissionByID(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[len("/api/admin/permissions/"):]

	id, err := strconv.Atoi(path)
	if err != nil {
		http.Error(w, "Invalid permission ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		GetPermission(w, r, id)
	case http.MethodPut:
		UpdatePermission(w, r, id)
	case http.MethodDelete:
		DeletePermission(w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func HandleTokens(w http.ResponseWriter, r *http.Request) {
	logger.Info("HandleTokens called with method: %s", r.Method)
	switch r.Method {
	case http.MethodGet:
		logger.Info("Calling GetTokens")
		GetTokens(w, r)
	case http.MethodPost:
		logger.Info("Calling CreateToken")
		CreateToken(w, r)
	default:
		logger.Info("Method not allowed: %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func HandleTokenByID(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[len("/api/admin/tokens/"):]

	if strings.HasSuffix(path, "/revoke") {
		idStr := path[:len(path)-len("/revoke")]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid token ID", http.StatusBadRequest)
			return
		}

		if r.Method == http.MethodPost {
			RevokeToken(w, r, id)
			return
		}
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id, err := strconv.Atoi(path)
	if err != nil {
		http.Error(w, "Invalid token ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		GetToken(w, r, id)
	case http.MethodPut:
		UpdateToken(w, r, id)
	case http.MethodDelete:
		DeleteToken(w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func HandleUserPermissions(w http.ResponseWriter, r *http.Request) {
	logger.Info("HandleUserPermissions called with method: %s", r.Method)
	switch r.Method {
	case http.MethodGet:
		logger.Info("Calling GetUserPermissions")
		GetUserPermissions(w, r)
	case http.MethodPost:
		logger.Info("Calling AssignPermission")
		AssignPermission(w, r)
	default:
		logger.Info("Method not allowed: %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func HandleUserPermissionByID(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[len("/api/admin/user-permissions/"):]

	id, err := strconv.Atoi(path)
	if err != nil {
		http.Error(w, "Invalid user permission ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodDelete:
		RemovePermission(w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
