package users

import (
	"net/http"
	"strconv"

	"github.com/holonet/core/logger"
)

// HandleListUsers handles GET requests to /api/users/list
func HandleListUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	logger.Info("HandleListUsers called")
	GetUsers(w, r)
}

// HandleCreateUser handles POST requests to /api/users/create
func HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	logger.Info("HandleCreateUser called")
	CreateUser(w, r)
}

// HandleGetUser handles GET requests to /api/users/get/{id}
func HandleGetUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract the ID from the URL
	idStr := r.URL.Path[len("/api/users/get/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	logger.Info("HandleGetUser called with ID: %d", id)
	GetUser(w, r, id)
}

// HandleUpdateUser handles PUT requests to /api/users/update/{id}
func HandleUpdateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract the ID from the URL
	idStr := r.URL.Path[len("/api/users/update/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	logger.Info("HandleUpdateUser called with ID: %d", id)
	UpdateUser(w, r, id)
}

// HandleDeleteUser handles DELETE requests to /api/users/delete/{id}
func HandleDeleteUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract the ID from the URL
	idStr := r.URL.Path[len("/api/users/delete/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	logger.Info("HandleDeleteUser called with ID: %d", id)
	DeleteUser(w, r, id)
}
