package database

import (
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/holonet/core/logger"
)

func InitAdminUser(db *sql.DB) error {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE username = 'admin'").Scan(&count)
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

	logger.Info("Creating admin user")
	var userID int
	passwordHash, err := hashPassword("insecure")
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	err = tx.QueryRow(`
		INSERT INTO users (username, email, password_hash, created_at, updated_at, deleted_at)
		VALUES ('admin', 'admin@example.com', $1, NOW(), NOW(), NOW())
		RETURNING id
	`, passwordHash).Scan(&userID)
	if err != nil {
		return fmt.Errorf("failed to create admin user: %w", err)
	}

	logger.Info("Creating Super User group")
	var groupID int
	err = tx.QueryRow(`
		INSERT INTO groups (user_id, name, description, permissions, created_at, updated_at, deleted_at)
		VALUES ($1, 'Super User', 'Administrator group with full permissions', 'admin', NOW(), NOW(), NOW())
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
		INSERT INTO tokens (user_id, token, expires_at, created_at, updated_at, deleted_at)
		VALUES ($1, $2, $3, NOW(), NOW(), NOW())
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
