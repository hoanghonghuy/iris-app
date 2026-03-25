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

func (r *StudentRepo) ListByClass(ctx context.Context, classID uuid.UUID, limit, offset int) ([]model.Student, int, error) {
	const q = `
		SELECT s.student_id, s.school_id, s.current_class_id, s.full_name, s.dob, s.gender,
		       pc.code, pc.expires_at, pc.usage_count, pc.max_usage,
		       COUNT(*) OVER() AS total_count
		FROM students s
		LEFT JOIN (
			SELECT DISTINCT ON (student_id) student_id, code, expires_at, usage_count, max_usage
			FROM student_parent_codes
			WHERE expires_at > now() AND usage_count < max_usage
			ORDER BY student_id, expires_at DESC
		) pc ON pc.student_id = s.student_id
		WHERE s.current_class_id = $1
		ORDER BY s.full_name
		LIMIT $2 OFFSET $3;
	`
	rows, err := r.pool.Query(ctx, q, classID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var students []model.Student
	var total int
	for rows.Next() {
		var s model.Student
		if err := rows.Scan(
			&s.StudentID, &s.SchoolID, &s.CurrentClassID, &s.FullName, &s.DOB, &s.Gender,
			&s.ActiveParentCode, &s.CodeExpiresAt, &s.CodeUsageCount, &s.CodeMaxUsage,
			&total,
		); err != nil {
			return nil, 0, err
		}
		students = append(students, s)
	}
	return students, total, rows.Err()
}

// GetByStudentID lấy thông tin student theo ID
func (r *StudentRepo) GetByStudentID(ctx context.Context, studentID uuid.UUID) (*model.Student, error) {
	const q = `SELECT student_id, school_id, current_class_id, full_name, dob, gender
			   FROM students WHERE student_id = $1;`
	var student model.Student
	err := r.pool.QueryRow(ctx, q, studentID).Scan(&student.StudentID, &student.SchoolID,
		&student.CurrentClassID, &student.FullName, &student.DOB, &student.Gender)
	if err != nil {
		return nil, err
	}
	return &student, nil
}

// GetSchoolIDByStudentID lấy school_id của student
func (r *StudentRepo) GetSchoolIDByStudentID(ctx context.Context, studentID uuid.UUID) (uuid.UUID, error) {
	const q = `SELECT school_id FROM students WHERE student_id = $1;`
	var schoolID uuid.UUID
	err := r.pool.QueryRow(ctx, q, studentID).Scan(&schoolID)
	if err != nil {
		return uuid.Nil, err
	}
	return schoolID, nil
}

// CountStudentsBySchool đếm tổng số học sinh (nếu schoolID rỗng thì đếm toàn hệ thống)
func (r *StudentRepo) CountStudentsBySchool(ctx context.Context, schoolID *uuid.UUID) (int, error) {
	var q string
	var err error
	var count int

	if schoolID != nil {
		q = `SELECT COUNT(*) FROM students WHERE school_id = $1;`
		err = r.pool.QueryRow(ctx, q, *schoolID).Scan(&count)
	} else {
		q = `SELECT COUNT(*) FROM students;`
		err = r.pool.QueryRow(ctx, q).Scan(&count)
	}

	return count, err
}
