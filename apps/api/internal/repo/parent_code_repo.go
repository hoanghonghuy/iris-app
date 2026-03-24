package repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)


type ParentCodeRepo struct {
	pool *pgxpool.Pool
}

func NewParentCodeRepo(pool *pgxpool.Pool) *ParentCodeRepo {
	return &ParentCodeRepo{pool: pool}
}

// Create tạo parent code mới
func (r *ParentCodeRepo) Create(ctx context.Context, studentID uuid.UUID, code string, maxUsage int) error {
	const q = `
		INSERT INTO student_parent_codes (student_id, code, usage_count, max_usage)
		VALUES ($1, $2, 0, $3);
	`
	_, err := r.pool.Exec(ctx, q, studentID, code, maxUsage)
	return err
}

// FindByStudentID tìm code theo student_id
func (r *ParentCodeRepo) FindByStudentID(ctx context.Context, studentID uuid.UUID) (string, error) {
	const q = `SELECT code FROM student_parent_codes WHERE student_id = $1;`
	var code string
	err := r.pool.QueryRow(ctx, q, studentID).Scan(&code)
	return code, err
}

// FindByCode tìm thông tin theo code
func (r *ParentCodeRepo) FindByCode(ctx context.Context, code string) (*model.StudentParentCode, error) {
	const q = `
		SELECT code_id, student_id, code, usage_count, max_usage, expires_at
		FROM student_parent_codes
		WHERE code = $1;
	`
	info := &model.StudentParentCode{}
	err := r.pool.QueryRow(ctx, q, code).Scan(
		&info.CodeID,
		&info.StudentID,
		&info.Code,
		&info.UsageCount,
		&info.MaxUsage,
		&info.ExpiresAt,
	)
	return info, err
}

// IncrementUsage tăng usage_count (dùng nội bộ khi đã biết còn slot)
func (r *ParentCodeRepo) IncrementUsage(ctx context.Context, code string) error {
	const q = `
		UPDATE student_parent_codes
		SET usage_count = usage_count + 1
		WHERE code = $1;
	`
	_, err := r.pool.Exec(ctx, q, code)
	return err
}

// IncrementUsageIfNotMaxed tăng usage_count chỉ khi chưa đạt giới hạn.
// Trả về ErrParentCodeMaxUsageReached nếu đã full (atomic — tránh race condition).
func (r *ParentCodeRepo) IncrementUsageIfNotMaxed(ctx context.Context, code string) error {
	const q = `
		UPDATE student_parent_codes
		SET usage_count = usage_count + 1
		WHERE code = $1
		  AND usage_count < max_usage
		  AND expires_at > now();
	`
	tag, err := r.pool.Exec(ctx, q, code)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNoRowsUpdated
	}
	return nil
}

