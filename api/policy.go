package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/holonet/core/logger"
)

type TokenPolicy struct {
	ID                int    `json:"id"`
	Name              string `json:"name"`
	Description       string `json:"description"`
	RateLimitPerMin   int    `json:"rate_limit_per_min"`
	MaxRequestsPerDay int    `json:"max_requests_per_day"`
	Active            bool   `json:"active"`
}

func RegisterPolicyRoutes() {
	http.HandleFunc("/api/policies", tokenAuthMiddleware(handlePolicies))
	http.HandleFunc("/api/policies/", tokenAuthMiddleware(handlePolicyByID))
	http.HandleFunc("/api/tokens/policy", tokenAuthMiddleware(handleTokenPolicy))
}

func handlePolicies(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getPolicies(w, r)
	case http.MethodPost:
		createPolicy(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handlePolicyByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/api/policies/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid policy ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		getPolicy(w, r, id)
	case http.MethodPut:
		updatePolicy(w, r, id)
	case http.MethodDelete:
		deletePolicy(w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleTokenPolicy(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		TokenID  int `json:"token_id"`
		PolicyID int `json:"policy_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	query := `
		UPDATE tokens
		SET policy_id = $1, updated_at = NOW()
		WHERE id = $2
		RETURNING id
	`

	var id int
	err := dbHandler.DB.QueryRow(query, request.PolicyID, request.TokenID).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Token not found", http.StatusNotFound)
			return
		}
		logger.Error("Failed to update token policy: %v", err)
		http.Error(w, "Failed to update token policy", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Token policy updated successfully",
	})
}

func getPolicies(w http.ResponseWriter, r *http.Request) {
	query := `
		SELECT id, name, description, rate_limit_per_min, max_requests_per_day, active
		FROM token_policies
		ORDER BY id
	`

	rows, err := dbHandler.DB.Query(query)
	if err != nil {
		logger.Error("Failed to list policies: %v", err)
		http.Error(w, "Failed to list policies", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	policies := []TokenPolicy{}
	for rows.Next() {
		policy := TokenPolicy{}
		err := rows.Scan(
			&policy.ID,
			&policy.Name,
			&policy.Description,
			&policy.RateLimitPerMin,
			&policy.MaxRequestsPerDay,
			&policy.Active,
		)
		if err != nil {
			logger.Error("Failed to scan policy: %v", err)
			http.Error(w, "Failed to list policies", http.StatusInternalServerError)
			return
		}
		policies = append(policies, policy)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(policies)
}

func getPolicy(w http.ResponseWriter, r *http.Request, id int) {
	query := `
		SELECT id, name, description, rate_limit_per_min, max_requests_per_day, active
		FROM token_policies
		WHERE id = $1
	`

	policy := TokenPolicy{}
	err := dbHandler.DB.QueryRow(query, id).Scan(
		&policy.ID,
		&policy.Name,
		&policy.Description,
		&policy.RateLimitPerMin,
		&policy.MaxRequestsPerDay,
		&policy.Active,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Policy not found", http.StatusNotFound)
			return
		}
		logger.Error("Failed to get policy: %v", err)
		http.Error(w, "Failed to get policy", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(policy)
}

func createPolicy(w http.ResponseWriter, r *http.Request) {
	var policy TokenPolicy
	if err := json.NewDecoder(r.Body).Decode(&policy); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	query := `
		INSERT INTO token_policies (name, description, rate_limit_per_min, max_requests_per_day, active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
		RETURNING id
	`

	err := dbHandler.DB.QueryRow(
		query,
		policy.Name,
		policy.Description,
		policy.RateLimitPerMin,
		policy.MaxRequestsPerDay,
		policy.Active,
	).Scan(&policy.ID)

	if err != nil {
		logger.Error("Failed to create policy: %v", err)
		http.Error(w, "Failed to create policy", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(policy)
}

func updatePolicy(w http.ResponseWriter, r *http.Request, id int) {
	var policy TokenPolicy
	if err := json.NewDecoder(r.Body).Decode(&policy); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	policy.ID = id

	query := `
		UPDATE token_policies
		SET name = $1, description = $2, rate_limit_per_min = $3, max_requests_per_day = $4, active = $5, updated_at = NOW()
		WHERE id = $6
		RETURNING id
	`

	err := dbHandler.DB.QueryRow(
		query,
		policy.Name,
		policy.Description,
		policy.RateLimitPerMin,
		policy.MaxRequestsPerDay,
		policy.Active,
		policy.ID,
	).Scan(&policy.ID)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Policy not found", http.StatusNotFound)
			return
		}
		logger.Error("Failed to update policy: %v", err)
		http.Error(w, "Failed to update policy", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(policy)
}

func deletePolicy(w http.ResponseWriter, r *http.Request, id int) {
	checkQuery := `
		SELECT COUNT(*)
		FROM tokens
		WHERE policy_id = $1
	`

	var count int
	err := dbHandler.DB.QueryRow(checkQuery, id).Scan(&count)
	if err != nil {
		logger.Error("Failed to check policy usage: %v", err)
		http.Error(w, "Failed to delete policy", http.StatusInternalServerError)
		return
	}

	if count > 0 {
		http.Error(w, "Cannot delete policy that is in use by tokens", http.StatusBadRequest)
		return
	}

	query := `
		DELETE FROM token_policies
		WHERE id = $1
		RETURNING id
	`

	var deletedID int
	err = dbHandler.DB.QueryRow(query, id).Scan(&deletedID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Policy not found", http.StatusNotFound)
			return
		}
		logger.Error("Failed to delete policy: %v", err)
		http.Error(w, "Failed to delete policy", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Policy deleted successfully",
	})
}
