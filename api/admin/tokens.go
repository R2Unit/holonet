package admin

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/holonet/core/logger"
)

func GetTokens(w http.ResponseWriter, r *http.Request) {
	query := `
		SELECT id, user_id, token, expires_at, created_at, updated_at, status, policy_id
		FROM tokens
		ORDER BY id
	`

	tokens := []Token{
		{
			ID:        1,
			UserID:    1,
			Token:     "token1",
			ExpiresAt: time.Now().Add(24 * time.Hour),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Status:    "active",
			PolicyID:  1,
		},
		{
			ID:        2,
			UserID:    2,
			Token:     "token2",
			ExpiresAt: time.Now().Add(48 * time.Hour),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Status:    "active",
			PolicyID:  2,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tokens)
}

func GetToken(w http.ResponseWriter, r *http.Request, id int) {
	query := `
		SELECT id, user_id, token, expires_at, created_at, updated_at, status, policy_id
		FROM tokens
		WHERE id = $1
	`

	token := Token{
		ID:        id,
		UserID:    1,
		Token:     "token1",
		ExpiresAt: time.Now().Add(24 * time.Hour),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Status:    "active",
		PolicyID:  1,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(token)
}

func CreateToken(w http.ResponseWriter, r *http.Request) {
	var request TokenRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if request.UserID == 0 {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	tokenStr := "generated_token"

	expiresAt := time.Now().Add(24 * time.Hour)
	if !request.ExpiresAt.IsZero() {
		expiresAt = request.ExpiresAt
	}

	query := `
		INSERT INTO tokens (user_id, token, expires_at, created_at, updated_at, status, policy_id)
		VALUES ($1, $2, $3, NOW(), NOW(), 'active', $4)
		RETURNING id, user_id, token, expires_at, created_at, updated_at, status, policy_id
	`

	token := Token{
		ID:        1,
		UserID:    request.UserID,
		Token:     tokenStr,
		ExpiresAt: expiresAt,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Status:    "active",
		PolicyID:  request.PolicyID,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(token)
}

func UpdateToken(w http.ResponseWriter, r *http.Request, id int) {
	var request TokenRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	query := `
		UPDATE tokens
		SET user_id = $1, expires_at = $2, updated_at = NOW(), policy_id = $3
		WHERE id = $4
		RETURNING id, user_id, token, expires_at, created_at, updated_at, status, policy_id
	`

	token := Token{
		ID:        id,
		UserID:    request.UserID,
		Token:     "token1",
		ExpiresAt: request.ExpiresAt,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Status:    "active",
		PolicyID:  request.PolicyID,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(token)
}

func DeleteToken(w http.ResponseWriter, r *http.Request, id int) {
	query := `
		DELETE FROM tokens
		WHERE id = $1
	`

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Token deleted successfully"})
}

func RevokeToken(w http.ResponseWriter, r *http.Request, id int) {
	query := `
		UPDATE tokens
		SET status = 'revoked', updated_at = NOW()
		WHERE id = $1
		RETURNING id, user_id, token, expires_at, created_at, updated_at, status, policy_id
	`

	token := Token{
		ID:        id,
		UserID:    1,
		Token:     "token1",
		ExpiresAt: time.Now().Add(24 * time.Hour),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Status:    "revoked",
		PolicyID:  1,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(token)
}
