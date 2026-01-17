package model

import "github.com/google/uuid"

// User đại diện cho bản ghi users trong database
type User struct {
	ID           uuid.UUID `db:"user_id"`
	Email        string    `db:"email"`
	PasswordHash string    `db:"password_hash"`
	Status       string    `db:"status"`
}

// UserInfo chứa thông tin user trả cho client
type UserInfo struct {
	ID     uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
	Status string    `json:"status"`
	Roles  []string  `json:"roles"`
}
