package repo

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ResetTokenRepo quản lý bảng password_reset_tokens
type ResetTokenRepo struct {
	pool *pgxpool.Pool
}

func NewResetTokenRepo(pool *pgxpool.Pool) *ResetTokenRepo {
	return &ResetTokenRepo{pool: pool}
}

// Create lưu token hash mới vào database
func (r *ResetTokenRepo) Create(ctx context.Context, userID uuid.UUID, tokenHash string, expiresAt time.Time) error {
	const q = `
		INSERT INTO password_reset_tokens (user_id, token_hash, expires_at)
		VALUES ($1, $2, $3);
	`
	_, err := r.pool.Exec(ctx, q, userID, tokenHash, expiresAt)
	return err
}

// FindByTokenHash tìm token chưa dùng theo hash
func (r *ResetTokenRepo) FindByTokenHash(ctx context.Context, tokenHash string) (*model.ResetToken, error) {
	const q = `
		SELECT id, user_id, token_hash, expires_at, used_at, created_at
		FROM password_reset_tokens
		WHERE token_hash = $1 AND used_at IS NULL
		LIMIT 1;
	`
	t := &model.ResetToken{}
	err := r.pool.QueryRow(ctx, q, tokenHash).Scan(
		&t.ID, &t.UserID, &t.TokenHash, &t.ExpiresAt, &t.UsedAt, &t.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return t, nil
}

// MarkUsed đánh dấu token đã sử dụng (chỉ update khi used_at IS NULL).
// Trả về ErrNoRowsUpdated nếu token đã dùng (race condition guard: mỗi token chỉ có thể được sử dụng thành công một lần duy nhất).
func (r *ResetTokenRepo) MarkUsed(ctx context.Context, tokenID uuid.UUID) error {
	const q = `
		UPDATE password_reset_tokens
		SET used_at = now()
		WHERE id = $1 AND used_at IS NULL;
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

// DeleteExpired dọn dẹp token hết hạn (gọi định kỳ nếu cần)
func (r *ResetTokenRepo) DeleteExpired(ctx context.Context) (int64, error) {
	const q = `
		DELETE FROM password_reset_tokens
		WHERE expires_at < now() OR used_at IS NOT NULL;
	`
	ct, err := r.pool.Exec(ctx, q)
	if err != nil {
		return 0, err
	}
	return ct.RowsAffected(), nil
}
