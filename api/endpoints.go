package api

import (
	"net/http"

	"github.com/holonet/core/api/users"
)

// RegisterEndpoints registers all API endpoints
func RegisterEndpoints() {
	// Register workflow routes
	http.HandleFunc("/api/workflows", tokenAuthMiddleware(handleWorkflows))
	http.HandleFunc("/api/workflows/", tokenAuthMiddleware(handleWorkflowByID))
	http.HandleFunc("/api/workflows/schedule", tokenAuthMiddleware(handleScheduleWorkflow))

	// Register policy routes
	http.HandleFunc("/api/policies", tokenAuthMiddleware(handlePolicies))
	http.HandleFunc("/api/policies/", tokenAuthMiddleware(handlePolicyByID))
	http.HandleFunc("/api/tokens/policy", tokenAuthMiddleware(handleTokenPolicy))

	// Register user routes - RESTful style (original)
	http.HandleFunc("/api/users", users.HandleUsers)
	http.HandleFunc("/api/users/", tokenAuthMiddleware(users.HandleUserByID))

	// Register user routes - Action-based style (new)
	http.HandleFunc("/api/users/list", tokenAuthMiddleware(users.HandleListUsers))
	http.HandleFunc("/api/users/create", tokenAuthMiddleware(users.HandleCreateUser))
	http.HandleFunc("/api/users/get/", tokenAuthMiddleware(users.HandleGetUser))
	http.HandleFunc("/api/users/update/", tokenAuthMiddleware(users.HandleUpdateUser))
	http.HandleFunc("/api/users/delete/", tokenAuthMiddleware(users.HandleDeleteUser))

	// Register admin routes
	http.HandleFunc("/api/admin/token", tokenAuthMiddleware(handleGenerateAdminToken))

	//http.HandleFunc("/api/test", handleTest)
}

//func handleTest(w http.ResponseWriter, r *http.Request) {
//	logger.Info("Test endpoint called with method: %s", r.Method)
//	w.Header().Set("Content-Type", "application/json")
//	logger.Info("Content-Type header set to application/json")
//	err := json.NewEncoder(w).Encode(map[string]string{"message": "Test endpoint working"})
//	if err != nil {
//		logger.Error("Failed to encode response to JSON: %v", err)
//		http.Error(w, "Failed to encode response to JSON", http.StatusInternalServerError)
//		return
//	}
//	logger.Info("Successfully encoded response to JSON")
//}
