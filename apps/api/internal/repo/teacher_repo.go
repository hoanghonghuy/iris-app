package repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TeacherRepo struct {
	pool *pgxpool.Pool
}

func NewTeacherRepo(pool *pgxpool.Pool) *TeacherRepo {
	return &TeacherRepo{
		pool: pool,
	}
}

func (r *TeacherRepo) List(ctx context.Context, schoolID *uuid.UUID, limit, offset int) ([]model.Teacher, int, error) {
	const qAll = `
		SELECT t.teacher_id, t.user_id, u.email, t.full_name, COALESCE(t.phone,''), t.school_id,
		       COUNT(*) OVER() as total_count
		FROM teachers t
		JOIN users u ON u.user_id = t.user_id
		ORDER BY t.full_name
		LIMIT $1 OFFSET $2;
	`
	const qBySchool = `
		SELECT t.teacher_id, t.user_id, u.email, t.full_name, COALESCE(t.phone,''), t.school_id,
		       COUNT(*) OVER() as total_count
		FROM teachers t
		JOIN users u ON u.user_id = t.user_id
		WHERE t.school_id = $3
		ORDER BY t.full_name
		LIMIT $1 OFFSET $2;
	`

	var rows pgx.Rows
	var err error
	if schoolID != nil {
		rows, err = r.pool.Query(ctx, qBySchool, limit, offset, *schoolID)
	} else {
		rows, err = r.pool.Query(ctx, qAll, limit, offset)
	}
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var teachers []model.Teacher
	var total int
	for rows.Next() {
		var t model.Teacher
		if err := rows.Scan(&t.TeacherID, &t.UserID, &t.Email, &t.FullName, &t.Phone, &t.SchoolID, &total); err != nil {
			return nil, 0, err
		}
		teachers = append(teachers, t)
	}
	return teachers, total, rows.Err()
}

// lấy thông tin giáo viên theo teacher_id.
func (r *TeacherRepo) GetByTeacherID(ctx context.Context, teacherID uuid.UUID) (*model.Teacher, error) {
	const q = `
		SELECT t.teacher_id, t.user_id, u.email, t.full_name, COALESCE(t.phone,''), t.school_id
		FROM teachers t
		JOIN users u ON u.user_id = t.user_id
		WHERE t.teacher_id = $1
		LIMIT 1;
	`
	teacher := &model.Teacher{}

	err := r.pool.QueryRow(ctx, q, teacherID).Scan(&teacher.TeacherID, &teacher.UserID, &teacher.Email,
		&teacher.FullName, &teacher.Phone, &teacher.SchoolID)
	if err != nil {
		return nil, err
	}
	return teacher, nil
}

// Update cập nhật thông tin teacher (admin có thể update tất cả fields)
func (r *TeacherRepo) Update(ctx context.Context, teacherID uuid.UUID, fullName, phone string, schoolID uuid.UUID) error {
	const q = `
		UPDATE teachers
		SET full_name = $2, phone = $3, school_id = $4, updated_at = now()
		WHERE teacher_id = $1;
	`
	_, err := r.pool.Exec(ctx, q, teacherID, fullName, phone, schoolID)
	return err
}

// UpdatePhone chỉ cập nhật phone (teacher chỉ có thể update phone của chính mình)
func (r *TeacherRepo) UpdatePhone(ctx context.Context, teacherID uuid.UUID, phone string) error {
	const q = `
		UPDATE teachers
		SET phone = $2, updated_at = now()
		WHERE teacher_id = $1;
	`
	_, err := r.pool.Exec(ctx, q, teacherID, phone)
	return err
}

// GetByUserID lấy thông tin teacher theo user_id
func (r *TeacherRepo) GetByUserID(ctx context.Context, userID uuid.UUID) (*model.Teacher, error) {
	const q = `
		SELECT t.teacher_id, t.user_id, u.email, t.full_name, COALESCE(t.phone,''), t.school_id
		FROM teachers t
		JOIN users u ON u.user_id = t.user_id
		WHERE t.user_id = $1
		LIMIT 1;
	`
	teacher := &model.Teacher{}

	err := r.pool.QueryRow(ctx, q, userID).Scan(&teacher.TeacherID, &teacher.UserID, &teacher.Email,
		&teacher.FullName, &teacher.Phone, &teacher.SchoolID)
	if err != nil {
		return nil, err
	}
	return teacher, nil
}
