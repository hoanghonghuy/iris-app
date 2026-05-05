package model

import (
	"time"

	"github.com/google/uuid"
)

// RefreshToken đại diện cho bản ghi refresh_tokens trong database.
type RefreshToken struct {
	ID        uuid.UUID  `db:"id"`
	UserID    uuid.UUID  `db:"user_id"`
	TokenHash string     `db:"token_hash"`
	ExpiresAt time.Time  `db:"expires_at"`
	RevokedAt *time.Time `db:"revoked_at"` // nil = chưa thu hồi
	CreatedAt time.Time  `db:"created_at"`
}
