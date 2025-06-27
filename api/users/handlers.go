package users

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/holonet/core/logger"
)

func HandleUsers(w http.ResponseWriter, r *http.Request) {
	logger.Info("HandleUsers called with method: %s", r.Method)
	switch r.Method {
	case http.MethodGet:
		logger.Info("Calling GetUsers")
		GetUsers(w, r)
	case http.MethodPost:
		logger.Info("Calling CreateUser")
		CreateUser(w, r)
	default:
		logger.Info("Method not allowed: %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func HandleUserByID(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[len("/api/users/"):]

	if strings.HasSuffix(path, "/suspend") {
		idStr := path[:len(path)-len("/suspend")]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		if r.Method == http.MethodPost {
			SuspendUser(w, r, id)
			return
		}
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if strings.HasSuffix(path, "/unsuspend") {
		idStr := path[:len(path)-len("/unsuspend")]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		if r.Method == http.MethodPost {
			UnsuspendUser(w, r, id)
			return
		}
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id, err := strconv.Atoi(path)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		GetUser(w, r, id)
	case http.MethodPut:
		UpdateUser(w, r, id)
	case http.MethodDelete:
		DeleteUser(w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
