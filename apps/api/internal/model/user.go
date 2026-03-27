package model

import (
	"time"

	"github.com/google/uuid"
)

// User đại diện cho bản ghi users trong database
type User struct {
	UserID       uuid.UUID `db:"user_id"`
	Email        string    `db:"email"`
	PasswordHash string    `db:"password_hash"`
	Status       string    `db:"status"`
	GoogleSub    string    `db:"google_sub"`
}

// UserInfo chứa thông tin user trả cho client
type UserInfo struct {
	UserID   uuid.UUID `json:"user_id"`
	Email    string    `json:"email"`
	FullName string    `json:"full_name,omitempty"`
	Status   string    `json:"status"`
	Roles    []string  `json:"roles"`
}

// UserWithToken dùng cho activation flow
type UserWithToken struct {
	UserID          uuid.UUID
	Email           string
	PasswordHash    string
	Status          string
	ActivationToken string
	TokenExpiresAt  time.Time
}
