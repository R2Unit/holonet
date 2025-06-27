package admin

import (
	"time"
)

type Permission struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type PermissionRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Token struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Status    string    `json:"status"`
	PolicyID  int       `json:"policy_id,omitempty"`
}

type TokenRequest struct {
	UserID    int       `json:"user_id"`
	ExpiresAt time.Time `json:"expires_at,omitempty"`
	PolicyID  int       `json:"policy_id,omitempty"`
}

type UserPermission struct {
	ID           int       `json:"id"`
	UserID       int       `json:"user_id"`
	PermissionID int       `json:"permission_id"`
	CreatedAt    time.Time `json:"created_at"`
}

type UserPermissionRequest struct {
	UserID       int `json:"user_id"`
	PermissionID int `json:"permission_id"`
}
