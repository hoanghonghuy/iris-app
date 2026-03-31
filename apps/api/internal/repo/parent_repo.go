package repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/model"
	"github.com/jackc/pgx/v5"
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

// List lấy danh sách phụ huynh, có thể lọc theo schoolID (nếu không nil).
func (r *ParentRepo) List(ctx context.Context, schoolID *uuid.UUID, limit, offset int) ([]model.Parent, int, error) {
	const qAll = `
		SELECT p.parent_id, p.user_id, u.email, p.full_name, COALESCE(p.phone,''), p.school_id,
		       COALESCE(
		           json_agg(
		               json_build_object('student_id', s.student_id, 'full_name', s.full_name)
		           ) FILTER (WHERE s.student_id IS NOT NULL), 
		           '[]'
		       ) as children,
		       COUNT(*) OVER() as total_count
		FROM parents p
		JOIN users u ON u.user_id = p.user_id
		LEFT JOIN student_parents sp ON p.parent_id = sp.parent_id
		LEFT JOIN students s ON sp.student_id = s.student_id
		GROUP BY p.parent_id, u.email
		ORDER BY p.full_name
		LIMIT $1 OFFSET $2;
	`
	const qBySchool = `
		SELECT p.parent_id, p.user_id, u.email, p.full_name, COALESCE(p.phone,''), p.school_id,
		       COALESCE(
		           json_agg(
		               json_build_object('student_id', s.student_id, 'full_name', s.full_name)
		           ) FILTER (WHERE s.student_id IS NOT NULL), 
		           '[]'
		       ) as children,
		       COUNT(*) OVER() as total_count
		FROM parents p
		JOIN users u ON u.user_id = p.user_id
		LEFT JOIN student_parents sp ON p.parent_id = sp.parent_id
		LEFT JOIN students s ON sp.student_id = s.student_id
		WHERE p.school_id = $3
		GROUP BY p.parent_id, u.email
		ORDER BY p.full_name
		LIMIT $1 OFFSET $2;
	`

	var rows pgx.Rows
	var err error
	if schoolID != nil { // => chỉ lấy phụ huynh thuộc trường cụ thể
		rows, err = r.pool.Query(ctx, qBySchool, limit, offset, *schoolID)
	} else { // => lấy tất cả phụ huynh
		rows, err = r.pool.Query(ctx, qAll, limit, offset)
	}
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var out []model.Parent
	var total int
	for rows.Next() {
		var x model.Parent
		if err := rows.Scan(&x.ParentID, &x.UserID, &x.Email, &x.FullName, &x.Phone, &x.SchoolID, &x.Children, &total); err != nil {
			return nil, 0, err
		}
		out = append(out, x)
	}
	return out, total, rows.Err()
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
		SELECT p.parent_id, p.user_id, u.email, p.full_name, COALESCE(p.phone,''), p.school_id,
		       COALESCE(
		           json_agg(
		               json_build_object('student_id', s.student_id, 'full_name', s.full_name)
		           ) FILTER (WHERE s.student_id IS NOT NULL), 
		           '[]'
		       ) as children
		FROM parents p
		JOIN users u ON u.user_id = p.user_id
		LEFT JOIN student_parents sp ON p.parent_id = sp.parent_id
		LEFT JOIN students s ON sp.student_id = s.student_id
		WHERE p.parent_id = $1
		GROUP BY p.parent_id, u.email
		LIMIT 1;
	`
	parent := &model.Parent{}
	err := r.pool.QueryRow(ctx, q, parentID).Scan(&parent.ParentID, &parent.UserID,
		&parent.Email, &parent.FullName, &parent.Phone, &parent.SchoolID, &parent.Children)
	if err != nil {
		return nil, err
	}
	return parent, nil
}
