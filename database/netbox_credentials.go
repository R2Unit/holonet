package database

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/holonet/core/logger"
)

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

func StoreNetboxCredentials(db *sql.DB, credentials NetboxCredentials) error {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM netbox_credentials WHERE user_id = $1", credentials.UserID).Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check if NetBox credentials exist: %w", err)
	}

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

func encryptValue(value string) (string, error) {
	key := []byte("holonet-netbox-encryption-key")

	hashedKey := sha256.Sum256(key)

	encrypted := make([]byte, len(value))
	for i := 0; i < len(value); i++ {
		encrypted[i] = value[i] ^ hashedKey[i%len(hashedKey)]
	}

	return hex.EncodeToString(encrypted), nil
}

func decryptValue(encryptedHex string) (string, error) {
	encrypted, err := hex.DecodeString(encryptedHex)
	if err != nil {
		return "", fmt.Errorf("failed to decode encrypted value: %w", err)
	}

	key := []byte("holonet-netbox-encryption-key")

	hashedKey := sha256.Sum256(key)

	decrypted := make([]byte, len(encrypted))
	for i := 0; i < len(encrypted); i++ {
		decrypted[i] = encrypted[i] ^ hashedKey[i%len(hashedKey)]
	}

	return string(decrypted), nil
}
