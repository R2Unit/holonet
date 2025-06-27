package database

import (
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"os"
	"time"

	"github.com/holonet/core/logger"
)

func InitAdminUser(db *sql.DB) error {
	// Get admin username from environment variable or use default
	username := getEnvOrDefault("ADMIN_USERNAME", "admin")

	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE username = $1", username).Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check if admin user exists: %w", err)
	}

	if count > 0 {
		logger.Info("Admin user already exists, skipping initialization")
		return nil
	}

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Get admin user details from environment variables or use defaults
	email := getEnvOrDefault("ADMIN_EMAIL", "admin@example.com")
	password := getEnvOrDefault("ADMIN_PASSWORD", "insecure")

	logger.Info("Creating admin user: %s (%s)", username, email)
	var userID int
	passwordHash, err := hashPassword(password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	err = tx.QueryRow(`
		INSERT INTO users (username, email, password_hash, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW())
		RETURNING id
	`, username, email, passwordHash).Scan(&userID)
	if err != nil {
		return fmt.Errorf("failed to create admin user: %w", err)
	}

	logger.Info("Creating Super User group")
	var groupID int
	err = tx.QueryRow(`
		INSERT INTO groups (user_id, name, description, permissions, created_at, updated_at)
		VALUES ($1, 'Super User', 'Administrator group with full permissions', 'admin', NOW(), NOW())
		RETURNING id
	`, userID).Scan(&groupID)
	if err != nil {
		return fmt.Errorf("failed to create Super User group: %w", err)
	}

	logger.Info("Creating token for admin user")
	token, err := generateRandomToken(32)
	if err != nil {
		return fmt.Errorf("failed to generate token: %w", err)
	}

	expiresAt := time.Now().AddDate(1, 0, 0)

	_, err = tx.Exec(`
		INSERT INTO tokens (user_id, token, expires_at, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW())
	`, userID, token, expiresAt)
	if err != nil {
		return fmt.Errorf("failed to create token for admin user: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	logger.Info("Admin user initialized successfully with ID %d", userID)
	logger.Info("Admin token: %s", token)
	return nil
}

func hashPassword(password string) (string, error) {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:]), nil
}

func generateRandomToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}
	return hex.EncodeToString(bytes), nil
}

// getEnvOrDefault returns the value of the environment variable if it exists,
// otherwise it returns the default value
func getEnvOrDefault(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if exists {
		return value
	}
	return defaultValue
}
