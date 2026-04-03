package repo

import (
	"context"
	"time"

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
func (r *ParentCodeRepo) Create(ctx context.Context, studentID uuid.UUID, code string, maxUsage int, expiresAt time.Time) error {
	const q = `
		INSERT INTO student_parent_codes (student_id, code, usage_count, max_usage, expires_at)
		VALUES ($1, $2, 0, $3, $4);
	`
	_, err := r.pool.Exec(ctx, q, studentID, code, maxUsage, expiresAt)
	return err
}

// DeleteByStudentID query to delete parent codes by studentID
func (r *ParentCodeRepo) DeleteByStudentID(ctx context.Context, studentID uuid.UUID) error {
	const q = `DELETE FROM student_parent_codes WHERE student_id = $1;`
	_, err := r.pool.Exec(ctx, q, studentID)
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

// IncrementUsage tăng số lần sử dụng mà không kiểm tra max_usage hay expires_at. Chỉ sử dụng nội bộ khi đã xác nhận điều kiện.
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
// Trả về ErrNoRowsUpdated nếu đã full
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

// RegisterParentTxParams chứa các tham số cần thiết để đăng ký tài khoản Parent trong transaction.
type RegisterParentTxParams struct {
	UserID       uuid.UUID
	Email        string
	PasswordHash string
	FullName     string
	Phone        string
	SchoolID     uuid.UUID
	StudentID    uuid.UUID
	Code         string
	GoogleSub    string // Optional
}

// RegisterParentTx thực thi toàn bộ luồng đăng ký tài khoản Parent trong một transaction.
// Bao gồm: Tạo User, Assign Role PARENT, Tạo Parent, Link StudentParent,
// Set Google Sub (nếu có), và tăng Usage Code.
func (r *ParentCodeRepo) RegisterParentTx(ctx context.Context, p RegisterParentTxParams) (uuid.UUID, error) {
	tx, err := r.pool.Begin(ctx) // tx: transaction
	if err != nil {
		return uuid.Nil, err
	}
	defer tx.Rollback(ctx)

	// 1. Tạo User
	const qUser = `
		INSERT INTO users (user_id, email, password_hash, status)
		VALUES ($1, $2, $3, 'active');
	`
	if _, err := tx.Exec(ctx, qUser, p.UserID, p.Email, p.PasswordHash); err != nil {
		return uuid.Nil, err
	}

	// 2. Assign ROLE PARENT
	const qRole = `
		INSERT INTO user_roles (user_id, role_id)
		SELECT $1, role_id FROM roles WHERE name = 'PARENT';
	`
	if _, err := tx.Exec(ctx, qRole, p.UserID); err != nil {
		return uuid.Nil, err
	}

	// 3. Tạo Parent record
	const qParent = `
		INSERT INTO parents (user_id, school_id, full_name, phone)
		VALUES ($1, $2, $3, $4)
		RETURNING parent_id;
	`
	var parentID uuid.UUID
	if err := tx.QueryRow(ctx, qParent, p.UserID, p.SchoolID, p.FullName, p.Phone).Scan(&parentID); err != nil {
		return uuid.Nil, err
	}

	// 4. Liên kết Student - Parent
	const qLink = `
		INSERT INTO student_parents (student_id, parent_id, relationship)
		VALUES ($1, $2, 'parent')
		ON CONFLICT (student_id, parent_id) DO UPDATE
		SET relationship = EXCLUDED.relationship;
	`
	if _, err := tx.Exec(ctx, qLink, p.StudentID, parentID); err != nil {
		return uuid.Nil, err
	}

	// 5. Link Google Sub (nếu có)
	if p.GoogleSub != "" {
		const qSub = `
			UPDATE users 
			SET google_sub = $2 
			WHERE user_id = $1;
		`
		if _, err := tx.Exec(ctx, qSub, p.UserID, p.GoogleSub); err != nil {
			return uuid.Nil, err
		}
	}

	// 6. Tăng usage count và kiểm tra max_usage trong một thao tác duy nhất (đảm bảo an toàn khi nhiều người dùng dùng chung mã)
	const qInc = `
		UPDATE student_parent_codes
		SET usage_count = usage_count + 1
		WHERE code = $1
		  AND usage_count < max_usage
		  AND expires_at > now();
	`
	tag, err := tx.Exec(ctx, qInc, p.Code)
	if err != nil {
		return uuid.Nil, err
	}
	if tag.RowsAffected() == 0 {
		return uuid.Nil, ErrNoRowsUpdated // Parent code hết hạn hoặc đã đạt max usage
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return uuid.Nil, err
	}

	return parentID, nil
}
