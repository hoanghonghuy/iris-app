package repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ParentRepo struct {
	pool *pgxpool.Pool
}

func NewParentRepo(pool *pgxpool.Pool) *ParentRepo {
	return &ParentRepo{
		pool: pool,
	}
}

func (r *ParentRepo) List(ctx context.Context) ([]model.Parent, error) {
	const q = `
		SELECT p.parent_id, p.user_id, u.email, p.full_name, COALESCE(p.phone,''), p.school_id
		FROM parents p
		JOIN users u ON u.user_id = p.user_id
		ORDER BY p.full_name;
	`
	rows, err := r.pool.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []model.Parent
	for rows.Next() {
		var x model.Parent
		if err := rows.Scan(&x.ParentID, &x.UserID, &x.Email, &x.FullName, &x.Phone, &x.SchoolID); err != nil {
			return nil, err
		}
		out = append(out, x)
	}
	return out, rows.Err()
}

// Create thêm mới một parent vào cơ sở dữ liệu.
// Trả về parent_id vừa được tạo hoặc lỗi nếu có.
func (r *ParentRepo) Create(ctx context.Context, userID, schoolID uuid.UUID, fullName, phone string) (uuid.UUID, error) {
	const q = `
		INSERT INTO parents (user_id, school_id, full_name, phone)
		VALUES ($1, $2, $3, $4)
		RETURNING parent_id;
	`
	var id uuid.UUID
	err := r.pool.QueryRow(ctx, q, userID, schoolID, fullName, phone).Scan(&id)
	return id, err
}

// GetByParentID lấy thông tin phụ huynh theo parent_id.
func (r *ParentRepo) GetByParentID(ctx context.Context, parentID uuid.UUID) (*model.Parent, error) {
	const q = `
		SELECT p.parent_id, p.user_id, u.email, p.full_name, COALESCE(p.phone,''), p.school_id
		FROM parents p
		JOIN users u ON u.user_id = p.user_id
		WHERE p.parent_id = $1
		LIMIT 1;
	`
	parent := &model.Parent{}
	err := r.pool.QueryRow(ctx, q, parentID).Scan(&parent.ParentID, &parent.UserID,
		&parent.Email, &parent.FullName, &parent.Phone, &parent.SchoolID)
	if err != nil {
		return nil, err
	}
	return parent, nil
}
