package repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/model"
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

func (r *TeacherRepo) List(ctx context.Context) ([]model.Teacher, error) {
	const q = `
		SELECT t.teacher_id, t.user_id, u.email, t.full_name, COALESCE(t.phone,''), t.school_id
		FROM teachers t
		JORN users u ON u.user_id = t.user_id
		ORDER BY t.full_name;
	`
	rows, err := r.pool.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var teachers []model.Teacher
	for rows.Next() {
		var t model.Teacher
		if err := rows.Scan(&t.TeacherID, &t.UserID, &t.Email, &t.FullName, &t.Phone, &t.SchoolID); err != nil {
			return nil, err
		}
		teachers = append(teachers, t)
	}
	return teachers, rows.Err()
}

// lấy thông tin giáo viên theo teacher_id.
func (r *TeacherRepo) GetByTeacherID(ctx context.Context, teacherID uuid.UUID) (*model.Teacher, error) {
	const q = `SELECT teacher_id, user_id, email, full_name, phone, school_id FROM teachers WHERE teacher_id=$1 LIMIT 1;`
	teacher := &model.Teacher{}

	err := r.pool.QueryRow(ctx, q, teacherID).Scan(&teacher.TeacherID, &teacher.UserID, &teacher.Email,
		&teacher.FullName, &teacher.Phone, &teacher.SchoolID)
	if err != nil {
		return nil, err
	}
	return teacher, nil
}
