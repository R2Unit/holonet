package api

import (
	"net/http"

	"github.com/holonet/core/api/admin"
	"github.com/holonet/core/api/users"
)

func RegisterEndpoints() {
	http.HandleFunc("/api/workflows", tokenAuthMiddleware(handleWorkflows))
	http.HandleFunc("/api/workflows/", tokenAuthMiddleware(handleWorkflowByID))
	http.HandleFunc("/api/workflows/schedule", tokenAuthMiddleware(handleScheduleWorkflow))

	http.HandleFunc("/api/policies", tokenAuthMiddleware(handlePolicies))
	http.HandleFunc("/api/policies/", tokenAuthMiddleware(handlePolicyByID))
	http.HandleFunc("/api/tokens/policy", tokenAuthMiddleware(handleTokenPolicy))

	http.HandleFunc("/api/users", users.HandleUsers)
	http.HandleFunc("/api/users/", tokenAuthMiddleware(users.HandleUserByID))

	http.HandleFunc("/api/users/list", tokenAuthMiddleware(users.HandleListUsers))
	http.HandleFunc("/api/users/create", tokenAuthMiddleware(users.HandleCreateUser))
	http.HandleFunc("/api/users/get/", tokenAuthMiddleware(users.HandleGetUser))
	http.HandleFunc("/api/users/update/", tokenAuthMiddleware(users.HandleUpdateUser))
	http.HandleFunc("/api/users/delete/", tokenAuthMiddleware(users.HandleDeleteUser))

	http.HandleFunc("/api/admin/token", tokenAuthMiddleware(handleGenerateAdminToken))

	http.HandleFunc("/api/admin/permissions", tokenAuthMiddleware(admin.HandlePermissions))
	http.HandleFunc("/api/admin/permissions/", tokenAuthMiddleware(admin.HandlePermissionByID))

	http.HandleFunc("/api/admin/permissions/list", tokenAuthMiddleware(admin.HandleListPermissions))
	http.HandleFunc("/api/admin/permissions/create", tokenAuthMiddleware(admin.HandleCreatePermission))
	http.HandleFunc("/api/admin/permissions/get/", tokenAuthMiddleware(admin.HandleGetPermission))
	http.HandleFunc("/api/admin/permissions/update/", tokenAuthMiddleware(admin.HandleUpdatePermission))
	http.HandleFunc("/api/admin/permissions/delete/", tokenAuthMiddleware(admin.HandleDeletePermission))

	http.HandleFunc("/api/admin/tokens", tokenAuthMiddleware(admin.HandleTokens))
	http.HandleFunc("/api/admin/tokens/", tokenAuthMiddleware(admin.HandleTokenByID))

	http.HandleFunc("/api/admin/tokens/list", tokenAuthMiddleware(admin.HandleListTokens))
	http.HandleFunc("/api/admin/tokens/create", tokenAuthMiddleware(admin.HandleCreateToken))
	http.HandleFunc("/api/admin/tokens/get/", tokenAuthMiddleware(admin.HandleGetToken))
	http.HandleFunc("/api/admin/tokens/update/", tokenAuthMiddleware(admin.HandleUpdateToken))
	http.HandleFunc("/api/admin/tokens/delete/", tokenAuthMiddleware(admin.HandleDeleteToken))
	http.HandleFunc("/api/admin/tokens/revoke/", tokenAuthMiddleware(admin.HandleRevokeToken))

	http.HandleFunc("/api/admin/user-permissions", tokenAuthMiddleware(admin.HandleUserPermissions))
	http.HandleFunc("/api/admin/user-permissions/", tokenAuthMiddleware(admin.HandleUserPermissionByID))

	http.HandleFunc("/api/admin/user-permissions/list", tokenAuthMiddleware(admin.HandleListUserPermissions))
	http.HandleFunc("/api/admin/user-permissions/assign", tokenAuthMiddleware(admin.HandleAssignPermission))
	http.HandleFunc("/api/admin/user-permissions/remove/", tokenAuthMiddleware(admin.HandleRemovePermission))
}
