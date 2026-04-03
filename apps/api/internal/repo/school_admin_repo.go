package repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SchoolAdminRepo struct {
	pool *pgxpool.Pool
}

func NewSchoolAdminRepo(pool *pgxpool.Pool) *SchoolAdminRepo {
	return &SchoolAdminRepo{
		pool: pool,
	}
}

// Create thêm mới một school admin vào cơ sở dữ liệu.
func (r *SchoolAdminRepo) Create(ctx context.Context, userID, schoolID uuid.UUID, fullName, phone string) (uuid.UUID, error) {
	const q = `
		INSERT INTO school_admins (user_id, school_id, full_name, phone)
		VALUES ($1, $2, $3, $4)
		RETURNING admin_id;
	`
	var id uuid.UUID
	err := r.pool.QueryRow(ctx, q, userID, schoolID, fullName, phone).Scan(&id)
	return id, err
}

// GetByAdminID lấy thông tin school admin theo admin_id.
func (r *SchoolAdminRepo) GetByAdminID(ctx context.Context, adminID uuid.UUID) (*model.SchoolAdmin, error) {
	const q = `
		SELECT sa.admin_id, sa.user_id, u.email, sa.full_name, COALESCE(sa.phone,''), sa.school_id
		FROM school_admins sa
		JOIN users u ON u.user_id = sa.user_id
		WHERE sa.admin_id = $1
		LIMIT 1;
	`
	a := &model.SchoolAdmin{}
	err := r.pool.QueryRow(ctx, q, adminID).Scan(
		&a.AdminID, &a.UserID, &a.Email, &a.FullName, &a.Phone, &a.SchoolID,
	)
	if err != nil {
		return nil, err
	}
	return a, nil
}

// GetByUserID lấy thông tin school admin theo user_id (dùng cho login flow).
func (r *SchoolAdminRepo) GetByUserID(ctx context.Context, userID uuid.UUID) (*model.SchoolAdmin, error) {
	const q = `
		SELECT sa.admin_id, sa.user_id, u.email, sa.full_name, COALESCE(sa.phone,''), sa.school_id
		FROM school_admins sa
		JOIN users u ON u.user_id = sa.user_id
		WHERE sa.user_id = $1
		LIMIT 1;
	`
	a := &model.SchoolAdmin{}
	err := r.pool.QueryRow(ctx, q, userID).Scan(
		&a.AdminID, &a.UserID, &a.Email, &a.FullName, &a.Phone, &a.SchoolID,
	)
	if err != nil {
		return nil, err
	}
	return a, nil
}

// List lấy danh sách tất cả school admins.
func (r *SchoolAdminRepo) List(ctx context.Context, limit, offset int) ([]model.SchoolAdmin, int, error) {
	const q = `
		SELECT sa.admin_id, sa.user_id, u.email, sa.full_name, COALESCE(sa.phone,''), sa.school_id,
		       COUNT(*) OVER() as total_count
		FROM school_admins sa
		JOIN users u ON u.user_id = sa.user_id
		ORDER BY sa.full_name
		LIMIT $1 OFFSET $2;
	`
	rows, err := r.pool.Query(ctx, q, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var out []model.SchoolAdmin
	var total int
	for rows.Next() {
		var x model.SchoolAdmin
		if err := rows.Scan(&x.AdminID, &x.UserID, &x.Email, &x.FullName, &x.Phone, &x.SchoolID, &total); err != nil {
			return nil, 0, err
		}
		out = append(out, x)
	}
	return out, total, rows.Err()
}

// Delete xóa school admin theo admin_id.
func (r *SchoolAdminRepo) Delete(ctx context.Context, adminID uuid.UUID) error {
	const q = `DELETE FROM school_admins WHERE admin_id = $1;`
	tag, err := r.pool.Exec(ctx, q, adminID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}
