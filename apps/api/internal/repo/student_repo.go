package repo

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type StudentRepo struct {
	pool *pgxpool.Pool
}

func NewStudentRepo(pool *pgxpool.Pool) *StudentRepo {
	return &StudentRepo{
		pool: pool,
	}
}

func (r *StudentRepo) Create(ctx context.Context, schoolID, classID uuid.UUID,
	fullName string, dob time.Time, gender string) (uuid.UUID, error) {
	const q = `
			INSERT INTO students (school_id, current_class_id, full_name, dob, gender)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING student_id;
		`
	var id uuid.UUID
	err := r.pool.QueryRow(ctx, q, schoolID, classID, fullName, dob, gender).Scan(&id)
	return id, err
}

func (r *StudentRepo) ListByClass(ctx context.Context, classID uuid.UUID) ([]model.Student, error) {
	const q = `
		SELECT student_id, school_id, current_class_id, full_name, dob, gender
		FROM students
		WHERE current_class_id = $1
		ORDER BY full_name;
	`
	rows, err := r.pool.Query(ctx, q, classID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []model.Student
	for rows.Next() {
		var s model.Student
		if err := rows.Scan(&s.ID, &s.SchoolID, &s.CurrentClassID, &s.FullName, &s.DOB, &s.Gender); err != nil {
			return nil, err
		}
		students = append(students, s)
	}
	return students, rows.Err()
}
