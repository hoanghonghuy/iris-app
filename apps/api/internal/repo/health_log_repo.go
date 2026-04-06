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

type HealthLogRepo struct {
	pool *pgxpool.Pool
}

func NewHealthLogRepo(pool *pgxpool.Pool) *HealthLogRepo {
	return &HealthLogRepo{
		pool: pool,
	}
}

// Create tạo mới bản ghi sức khỏe cho học sinh
func (r *HealthLogRepo) Create(ctx context.Context, studentID, recordedBy uuid.UUID,
	temperature *float64, symptoms string, severity *string, note string) (uuid.UUID, error) {
	const q = `
		INSERT INTO health_logs (student_id, temperature, symptoms, severity, note, recorded_by)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING health_log_id;
	`

	var id uuid.UUID
	err := r.pool.QueryRow(ctx, q, studentID, temperature, symptoms, severity, note, recordedBy).Scan(&id)
	return id, err
}

// CreateByStudentAndTeacher tạo mới bản ghi sức khỏe nếu giáo viên được phân công dạy học sinh đó.
func (r *HealthLogRepo) CreateByStudentAndTeacher(ctx context.Context, teacherUserID, studentID uuid.UUID,
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
			return uuid.Nil, ErrNoRowsUpdated
		}
		return uuid.Nil, err
	}

	return id, nil
}

// ListByStudentAndTeacher liệt kê bản ghi sức khỏe cho học sinh nếu giáo viên được phân công dạy lớp của học sinh đó
func (r *HealthLogRepo) ListByStudentAndTeacher(ctx context.Context, teacherUserID, studentID uuid.UUID,
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

// ListByStudentForAdminReport liệt kê tất cả bản ghi sức khỏe của 1 học sinh.
// Reserved use-case: API báo cáo/tổng hợp cho admin trong giai đoạn tiếp theo.
func (r *HealthLogRepo) ListByStudentForAdminReport(ctx context.Context, studentID uuid.UUID,
	from, to time.Time) ([]model.HealthLog, error) {
	const q = `
		SELECT health_log_id, student_id, recorded_at, temperature,
			COALESCE(symptoms,''), severity, COALESCE(note,''), recorded_by
		FROM health_logs
		WHERE student_id = $1
			AND recorded_at BETWEEN $2 AND $3
		ORDER BY recorded_at DESC;
	`
	rows, err := r.pool.Query(ctx, q, studentID, from, to)
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

// ListByTeacherForAdminReport liệt kê bản ghi sức khỏe theo giáo viên.
// Reserved use-case: báo cáo tổng hợp theo giáo viên/lớp cho admin.
func (r *HealthLogRepo) ListByTeacherForAdminReport(ctx context.Context, teacherUserID uuid.UUID,
	from, to time.Time) ([]model.HealthLog, error) {
	const q = `
		SELECT hl.health_log_id, hl.student_id, hl.recorded_at, hl.temperature,
			COALESCE(hl.symptoms,''), hl.severity, COALESCE(hl.note,''), hl.recorded_by
		FROM health_logs hl
		JOIN students s ON s.student_id = hl.student_id
		JOIN teacher_classes tc ON tc.class_id = s.current_class_id
		JOIN teachers t ON t.teacher_id = tc.teacher_id
		WHERE t.user_id = $1
			AND hl.recorded_at BETWEEN $2 AND $3
		ORDER BY hl.recorded_at DESC;
	`
	rows, err := r.pool.Query(ctx, q, teacherUserID, from, to)
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
