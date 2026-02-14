package repo

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/model"
	"github.com/jackc/pgx/v5"
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

// CreateHealthLog tạo mới bản ghi sức khỏe cho học sinh nếu giáo viên được phân công dạy lớp của học sinh đó.
func (r *TeacherScopeRepo) CreateHealthLog(ctx context.Context, teacherUserID, studentID uuid.UUID,
	recordedAt *time.Time, temperature *float64, symptoms string, severity *string, note string) (uuid.UUID, error) {
	const q = `
		INSERT INTO health_logs (student_id, recorded_at, temperature, symptoms, severity, note, recorded_by)
		SELECT $2, COALESCE($3, now()), $4, $5, $6, $7, $1
		FROM students s
		JOIN teacher_classes tc ON tc.class_id = s.current_class_id
		JOIN teachers t ON t.teacher_id = tc.teacher_id
		WHERE t.user_id = $1 AND s.student_id = $2
		RETURNING health_log_id;
	`

	var id uuid.UUID
	err := r.pool.QueryRow(ctx, q, teacherUserID, studentID, recordedAt, temperature, symptoms, severity, note).Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return uuid.Nil, ErrForbidden
		}
		return uuid.Nil, err
	}

	return id, nil
}

// ListHealthLogsByStudent liệt kê bản ghi sức khỏe cho học sinh nếu giáo viên được phân công dạy lớp của học sinh đó.
func (r *TeacherScopeRepo) ListHealthLogsByStudent(ctx context.Context, teacherUserID, studentID uuid.UUID,
	from, to time.Time) ([]model.HealthLog, error) {
	const q = `
		SELECT hl.health_log_id, hl.student_id, hl.recorded_at, hl.temperature,
			COALESCE(hl.symptoms,''), hl.severity, COALESCE(hl.note,''), hl.recorded_by
		FROM health_logs hl
		JOIN students s ON s.student_id = hl.student_id
		JOIN teacher_classes tc ON tc.class_id = s.current_class_id
		JOIN teachers t ON t.teacher_id = tc.teacher_id
		WHERE t.user_id = $1
			AND hl.student_id = $2
			AND hl.recorded_at BETWEEN $3 AND $4
		ORDER BY hl.recorded_at DESC;
	`
	rows, err := r.pool.Query(ctx, q, teacherUserID, studentID, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []model.HealthLog
	for rows.Next() {
		var x model.HealthLog
		if err := rows.Scan(
			&x.HealthLogID, &x.StudentID, &x.RecordedAt, &x.Temperature,
			&x.Symptoms, &x.Severity, &x.Note, &x.RecordedBy,
		); err != nil {
			return nil, err
		}
		out = append(out, x)
	}
	return out, rows.Err()
}

// CreateClassPost tạo bài đăng cho một lớp học nếu giáo viên được phân công dạy lớp đó.
func (r *TeacherScopeRepo) CreateClassPost(ctx context.Context, teacherUserID, classID uuid.UUID,
	postType, content string) (uuid.UUID, error) {
	const q = `
		INSERT INTO posts (author_user_id, scope_type, class_id, type, content)
		SELECT $1, 'class', $2, $3, $4
		FROM teacher_classes tc
		JOIN teachers t ON t.teacher_id = tc.teacher_id
		WHERE t.user_id = $1 AND tc.class_id = $2
		RETURNING post_id;
	`

	var id uuid.UUID
	err := r.pool.QueryRow(ctx, q, teacherUserID, classID, postType, content).Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return uuid.Nil, ErrForbidden
		}
		return uuid.Nil, err
	}

	return id, nil
}

// CreateStudentPost tạo bài đăng cho một học sinh nếu giáo viên được phân công dạy lớp của học sinh đó.
func (r *TeacherScopeRepo) CreateStudentPost(ctx context.Context, teacherUserID, studentID uuid.UUID,
	postType, content string) (uuid.UUID, error) {
	const q = `
		INSERT INTO posts (author_user_id, scope_type, student_id, type, content)
		SELECT $1, 'student', $2, $3, $4
		FROM students s
		JOIN teacher_classes tc ON tc.class_id = s.current_class_id
		JOIN teachers t ON t.teacher_id = tc.teacher_id
		WHERE t.user_id = $1 AND s.student_id = $2
		RETURNING post_id;
	`

	var id uuid.UUID
	err := r.pool.QueryRow(ctx, q, teacherUserID, studentID, postType, content).Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return uuid.Nil, ErrForbidden
		}
		return uuid.Nil, err
	}

	return id, nil
}

// ListClassPosts liệt kê bài đăng của một lớp nếu giáo viên được phân công dạy lớp đó.
func (r *TeacherScopeRepo) ListClassPosts(ctx context.Context, teacherUserID, classID uuid.UUID,
	limit, offset int) ([]model.Post, error) {
	const q = `
		SELECT p.post_id, p.author_user_id, p.scope_type, p.school_id, p.class_id, p.student_id,
			p.type, p.content, p.created_at, p.updated_at
		FROM posts p
		JOIN teacher_classes tc ON tc.class_id = p.class_id
		JOIN teachers t ON t.teacher_id = tc.teacher_id
		WHERE t.user_id = $1 AND p.class_id = $2 AND p.scope_type = 'class'
		ORDER BY p.created_at DESC
		LIMIT $3 OFFSET $4;
	`

	rows, err := r.pool.Query(ctx, q, teacherUserID, classID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []model.Post
	for rows.Next() {
		var p model.Post
		if err := rows.Scan(
			&p.PostID, &p.AuthorUserID, &p.ScopeType, &p.SchoolID, &p.ClassID, &p.StudentID,
			&p.Type, &p.Content, &p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, rows.Err()
}

// ListAttendanceByStudent liệt kê lịch sử điểm danh của một học sinh nếu giáo viên được phân công dạy lớp của học sinh đó.
func (r *TeacherScopeRepo) ListAttendanceByStudent(ctx context.Context, teacherUserID, studentID uuid.UUID,
	from, to time.Time) ([]model.AttendanceRecord, error) {
	const q = `
		SELECT ar.attendance_id, ar.student_id, ar.date, ar.status,
			ar.check_in_at, ar.check_out_at, COALESCE(ar.note,''), ar.recorded_by
		FROM attendance_records ar
		JOIN students s ON s.student_id = ar.student_id
		JOIN teacher_classes tc ON tc.class_id = s.current_class_id
		JOIN teachers t ON t.teacher_id = tc.teacher_id
		WHERE t.user_id = $1
			AND ar.student_id = $2
			AND ar.date BETWEEN $3 AND $4
		ORDER BY ar.date DESC;
	`
	rows, err := r.pool.Query(ctx, q, teacherUserID, studentID, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []model.AttendanceRecord
	for rows.Next() {
		var x model.AttendanceRecord
		if err := rows.Scan(
			&x.AttendanceID, &x.StudentID, &x.Date, &x.Status,
			&x.CheckInAt, &x.CheckOutAt, &x.Note, &x.RecordedBy,
		); err != nil {
			return nil, err
		}
		out = append(out, x)
	}
	return out, rows.Err()
}

// ListStudentPosts liệt kê bài đăng của một học sinh nếu giáo viên được phân công dạy lớp của học sinh đó.
func (r *TeacherScopeRepo) ListStudentPosts(ctx context.Context, teacherUserID, studentID uuid.UUID,
	limit, offset int) ([]model.Post, error) {
	const q = `
		SELECT p.post_id, p.author_user_id, p.scope_type, p.school_id, p.class_id, p.student_id,
			p.type, p.content, p.created_at, p.updated_at
		FROM posts p
		JOIN students s ON s.student_id = p.student_id
		JOIN teacher_classes tc ON tc.class_id = s.current_class_id
		JOIN teachers t ON t.teacher_id = tc.teacher_id
		WHERE t.user_id = $1 AND p.student_id = $2 AND p.scope_type = 'student'
		ORDER BY p.created_at DESC
		LIMIT $3 OFFSET $4;
	`

	rows, err := r.pool.Query(ctx, q, teacherUserID, studentID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []model.Post
	for rows.Next() {
		var p model.Post
		if err := rows.Scan(
			&p.PostID, &p.AuthorUserID, &p.ScopeType, &p.SchoolID, &p.ClassID, &p.StudentID,
			&p.Type, &p.Content, &p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, rows.Err()
}
