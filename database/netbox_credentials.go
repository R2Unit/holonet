package database

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/holonet/core/logger"
)

// NetboxCredentials represents the credentials for a NetBox instance
type NetboxCredentials struct {
	ID             int
	UserID         int
	NetboxUsername string
	NetboxPassword string
	NetboxToken    string
	NetboxGroup    string
	NetboxHost     string
	IsEncrypted    bool
	LastVerifiedAt time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      time.Time
}

// StoreNetboxCredentials stores the NetBox credentials in the database
func StoreNetboxCredentials(db *sql.DB, credentials NetboxCredentials) error {
	// Check if credentials already exist for this user
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM netbox_credentials WHERE user_id = $1", credentials.UserID).Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check if NetBox credentials exist: %w", err)
	}

	// Encrypt the credentials if they're not already encrypted
	if !credentials.IsEncrypted {
		var err error
		credentials.NetboxPassword, err = encryptValue(credentials.NetboxPassword)
		if err != nil {
			return fmt.Errorf("failed to encrypt NetBox password: %w", err)
		}
		credentials.NetboxToken, err = encryptValue(credentials.NetboxToken)
		if err != nil {
			return fmt.Errorf("failed to encrypt NetBox token: %w", err)
		}
		credentials.IsEncrypted = true
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

	if count > 0 {
		// Update existing credentials
		_, err = tx.Exec(`
			UPDATE netbox_credentials
			SET netbox_username = $1, netbox_password = $2, netbox_token = $3, netbox_group = $4, 
				netbox_host = $5, is_encrypted = $6, last_verified_at = $7, updated_at = NOW()
			WHERE user_id = $8
		`, credentials.NetboxUsername, credentials.NetboxPassword, credentials.NetboxToken,
			credentials.NetboxGroup, credentials.NetboxHost, credentials.IsEncrypted,
			credentials.LastVerifiedAt, credentials.UserID)
		if err != nil {
			return fmt.Errorf("failed to update NetBox credentials: %w", err)
		}
		logger.Info("Updated NetBox credentials for user ID %d", credentials.UserID)
	} else {
		// Insert new credentials
		_, err = tx.Exec(`
			INSERT INTO netbox_credentials (user_id, netbox_username, netbox_password, netbox_token, 
				netbox_group, netbox_host, is_encrypted, last_verified_at, created_at, updated_at, deleted_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW(), NOW())
		`, credentials.UserID, credentials.NetboxUsername, credentials.NetboxPassword,
			credentials.NetboxToken, credentials.NetboxGroup, credentials.NetboxHost,
			credentials.IsEncrypted, credentials.LastVerifiedAt)
		if err != nil {
			return fmt.Errorf("failed to insert NetBox credentials: %w", err)
		}
		logger.Info("Inserted NetBox credentials for user ID %d", credentials.UserID)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// GetNetboxCredentials retrieves the NetBox credentials from the database
func GetNetboxCredentials(db *sql.DB, userID int) (*NetboxCredentials, error) {
	var credentials NetboxCredentials
	err := db.QueryRow(`
		SELECT id, user_id, netbox_username, netbox_password, netbox_token, netbox_group, 
			netbox_host, is_encrypted, last_verified_at, created_at, updated_at, deleted_at
		FROM netbox_credentials
		WHERE user_id = $1
	`, userID).Scan(
		&credentials.ID, &credentials.UserID, &credentials.NetboxUsername,
		&credentials.NetboxPassword, &credentials.NetboxToken, &credentials.NetboxGroup,
		&credentials.NetboxHost, &credentials.IsEncrypted, &credentials.LastVerifiedAt,
		&credentials.CreatedAt, &credentials.UpdatedAt, &credentials.DeletedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get NetBox credentials: %w", err)
	}

	// Decrypt the credentials if they're encrypted
	if credentials.IsEncrypted {
		var err error
		credentials.NetboxPassword, err = decryptValue(credentials.NetboxPassword)
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt NetBox password: %w", err)
		}
		credentials.NetboxToken, err = decryptValue(credentials.NetboxToken)
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt NetBox token: %w", err)
		}
		credentials.IsEncrypted = false
	}

	return &credentials, nil
}

// Simple XOR-based encryption/decryption
// Note: This is not secure for production use, but meets the requirement of
// "encrypted if possible without any packages"
func encryptValue(value string) (string, error) {
	// Use a fixed key for simplicity
	// In a real application, this should be a secure, randomly generated key
	key := []byte("holonet-netbox-encryption-key")

	// Hash the key to get a consistent length
	hashedKey := sha256.Sum256(key)

	// XOR the value with the key
	encrypted := make([]byte, len(value))
	for i := 0; i < len(value); i++ {
		encrypted[i] = value[i] ^ hashedKey[i%len(hashedKey)]
	}

	// Return the encrypted value as a hex string
	return hex.EncodeToString(encrypted), nil
}

// Decrypt a value that was encrypted with encryptValue
func decryptValue(encryptedHex string) (string, error) {
	// Decode the hex string
	encrypted, err := hex.DecodeString(encryptedHex)
	if err != nil {
		return "", fmt.Errorf("failed to decode encrypted value: %w", err)
	}

	// Use the same key as in encryptValue
	key := []byte("holonet-netbox-encryption-key")

	// Hash the key to get a consistent length
	hashedKey := sha256.Sum256(key)

	// XOR the encrypted value with the key to get the original value
	decrypted := make([]byte, len(encrypted))
	for i := 0; i < len(encrypted); i++ {
		decrypted[i] = encrypted[i] ^ hashedKey[i%len(hashedKey)]
	}

	return string(decrypted), nil
}
