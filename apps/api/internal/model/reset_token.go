package model

import (
	"time"

	"github.com/google/uuid"
)

// ResetToken đại diện cho bản ghi password_reset_tokens trong database
type ResetToken struct {
	ID        uuid.UUID  `db:"id"`
	UserID    uuid.UUID  `db:"user_id"`
	TokenHash string     `db:"token_hash"`
	ExpiresAt time.Time  `db:"expires_at"`
	UsedAt    *time.Time `db:"used_at"` // nil = chưa sử dụng
	CreatedAt time.Time  `db:"created_at"`
}
