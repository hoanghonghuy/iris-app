package repo

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

// RefreshTokenRepo quản lý bảng refresh_tokens.
type RefreshTokenRepo struct {
	pool *pgxpool.Pool
}

func NewRefreshTokenRepo(pool *pgxpool.Pool) *RefreshTokenRepo {
	return &RefreshTokenRepo{pool: pool}
}

func (r *RefreshTokenRepo) Create(ctx context.Context, userID uuid.UUID, tokenHash string, expiresAt time.Time) error {
	const q = `
		INSERT INTO refresh_tokens (user_id, token_hash, expires_at)
		VALUES ($1, $2, $3);
	`
	_, err := r.pool.Exec(ctx, q, userID, tokenHash, expiresAt)
	return err
}

// FindActiveByTokenHash tìm refresh token còn hiệu lực (chưa revoked + chưa hết hạn).
func (r *RefreshTokenRepo) FindActiveByTokenHash(ctx context.Context, tokenHash string) (*model.RefreshToken, error) {
	const q = `
		SELECT id, user_id, token_hash, expires_at, revoked_at, created_at
		FROM refresh_tokens
		WHERE token_hash = $1
			AND revoked_at IS NULL
			AND expires_at > now()
		LIMIT 1;
	`
	t := &model.RefreshToken{}
	err := r.pool.QueryRow(ctx, q, tokenHash).Scan(
		&t.ID, &t.UserID, &t.TokenHash, &t.ExpiresAt, &t.RevokedAt, &t.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return t, nil
}

// RevokeByID thu hồi refresh token (idempotent: chỉ update khi chưa revoked).
func (r *RefreshTokenRepo) RevokeByID(ctx context.Context, tokenID uuid.UUID) error {
	const q = `
		UPDATE refresh_tokens
		SET revoked_at = now()
		WHERE id = $1 AND revoked_at IS NULL;
	`
	tag, err := r.pool.Exec(ctx, q, tokenID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNoRowsUpdated
	}
	return nil
}
