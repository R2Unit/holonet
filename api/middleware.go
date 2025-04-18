package api

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func tokenAuthMiddleware(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := getBearerToken(r)
		if token == "" {
			http.Error(w, "Missing or invalid Authorization header", http.StatusUnauthorized)
			return
		}

		valid, err := validateToken(token, dbHandler.DB)
		if err != nil {
			log.Printf("Error validating token: %v", err)
			http.Error(w, "Error validating token", http.StatusInternalServerError)
			return
		}
		if !valid {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		handler(w, r)
	}
}

func getBearerToken(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return ""
	}
	return parts[1]
}

func validateToken(token string, db *sql.DB) (bool, error) {
	query := `
		SELECT COUNT(*) 
		FROM tokens 
		WHERE token = $1 AND expires_at > NOW();
	`
	var count int
	err := db.QueryRow(query, token).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("error querying token: %w", err)
	}
	return count > 0, nil
}
