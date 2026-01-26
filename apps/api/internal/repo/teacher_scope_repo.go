package repo

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrForbidden = errors.New("forbidden")

type TeacherScopeRepo struct {
	pool *pgxpool.Pool
}

func NewTeacherScopeRepo(pool *pgxpool.Pool) *TeacherScopeRepo {
	return &TeacherScopeRepo{
		pool: pool,
	}
}

// ListMyClass liệt kê các lớp học mà giáo viên (theo user_id) được phân công giảng dạy.
func (r *TeacherScopeRepo) ListMyClass(ctx context.Context, teacherUserID uuid.UUID) ([]model.Class, error) {
	const q = `
		SELECT c.class_id, c.school_id, c.name, c.school_year
		FROM classes c
		JOIN teacher_classes tc ON tc.class_id = c.class_id
		JOIN teachers t ON t.teacher_id = tc.teacher_id
		WHERE t.user_id = $1
		ORDER BY c.school_year DESC, c.name;
	`

	rows, err := r.pool.Query(ctx, q, teacherUserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var classes []model.Class
	for rows.Next() {
		var c model.Class
		if err := rows.Scan(&c.ID, &c.SchoolID, &c.Name, &c.SchoolYear); err != nil {
			return nil, err
		}
		classes = append(classes, c)
	}
	return classes, rows.Err()
}

// ListMyStudentsInClass liệt kê học sinh trong một lớp nếu giáo viên được phân công dạy lớp đó
func (r *TeacherScopeRepo) ListMyStudentsInClass(ctx context.Context, teacherUserID, classID uuid.UUID) ([]model.Student, error) {
	const q = `
		SELECT s.student_id, s.school_id, s.current_class_id, s.full_name, s.dob, s.gender
		FROM students s
		JOIN teacher_classes tc ON tc.class_id = s.current_class_id
		JOIN teachers t ON t.teacher_id = tc.teacher_id
		WHERE t.user_id = $1 AND s.current_class_id = $2
		ORDER BY s.full_name;
	`
	rows, err := r.pool.Query(ctx, q, teacherUserID, classID)
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

// UpsertAttendance: Giáo viên chỉ có thể điểm danh cho học sinh trong lớp của mình.
func (r *TeacherScopeRepo) UpsertAttendance(ctx context.Context, teacherUserID, studentID uuid.UUID,
	date time.Time, status string, checkInAt, checkOutAt *time.Time, note string) error {
	const q = `
			INSERT INTO attendance_records (student_id, date, status, check_in_at, check_out_at, note, recorded_by)
			SELECT s.student_id, $3, $4, $5, $6, $7, $1
			FROM students s
			JOIN teacher_classes tc ON tc.class_id = s.current_class_id
			JOIN teachers t ON t.teacher_id = tc.teacher_id
			WHERE t.user_id = $1 AND s.student_id = $2
			ON CONFLICT (student_id, date)
			DO UPDATE SET
			  status = EXCLUDED.status,
			  check_in_at = EXCLUDED.check_in_at,
			  check_out_at = EXCLUDED.check_out_at,
			  note = EXCLUDED.note,
			  recorded_by = EXCLUDED.recorded_by,
			  updated_at = now();
		`
	
	tag, err := r.pool.Exec(ctx, q, teacherUserID, studentID, date, status, checkInAt, checkOutAt, note)
	if err != nil {
		return err
	}
	
	if tag.RowsAffected() == 0 { // không có hàng nào được cập nhật, điều kiện WHERE không thỏa mãn
		return ErrForbidden
	}
	return nil
}
