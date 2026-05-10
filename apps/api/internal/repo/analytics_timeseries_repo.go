package repo

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// AnalyticsTimeseriesRepo truy vấn aggregate theo ngày cho dashboard chart (additive API).
type AnalyticsTimeseriesRepo struct {
	pool *pgxpool.Pool
}

// NewAnalyticsTimeseriesRepo tạo repo.
func NewAnalyticsTimeseriesRepo(pool *pgxpool.Pool) *AnalyticsTimeseriesRepo {
	return &AnalyticsTimeseriesRepo{pool: pool}
}

// DayFloat điểm (ngày UTC 00:00, giá trị float).
type DayFloat struct {
	Day   time.Time
	Value float64
}

// DayInt điểm (ngày UTC 00:00, giá trị int).
type DayInt struct {
	Day   time.Time
	Value int
}

// StatusDayCount đếm theo ngày và trạng thái (appointments).
type StatusDayCount struct {
	Day    time.Time
	Status string
	Count  int
}

func normalizeDayUTC(d time.Time) time.Time {
	u := d.UTC()
	return time.Date(u.Year(), u.Month(), u.Day(), 0, 0, 0, 0, time.UTC)
}

func scanDayFloatRows(rows pgx.Rows) ([]DayFloat, error) {
	defer rows.Close()
	var out []DayFloat
	for rows.Next() {
		var d time.Time
		var v float64
		if err := rows.Scan(&d, &v); err != nil {
			return nil, err
		}
		out = append(out, DayFloat{Day: normalizeDayUTC(d), Value: v})
	}
	return out, rows.Err()
}

func scanDayIntRows(rows pgx.Rows) ([]DayInt, error) {
	defer rows.Close()
	var out []DayInt
	for rows.Next() {
		var d time.Time
		var v int
		if err := rows.Scan(&d, &v); err != nil {
			return nil, err
		}
		out = append(out, DayInt{Day: normalizeDayUTC(d), Value: v})
	}
	return out, rows.Err()
}

// AdminAttendancePresentRateDaily tỷ lệ học sinh có mặt (present) / tổng học sinh trong scope, theo ngày.
func (r *AnalyticsTimeseriesRepo) AdminAttendancePresentRateDaily(ctx context.Context, schoolID *uuid.UUID, fromDate, toDate time.Time) ([]DayFloat, error) {
	from := fromDate.UTC().Truncate(24 * time.Hour)
	to := toDate.UTC().Truncate(24 * time.Hour)
	if schoolID != nil {
		const q = `
WITH days AS (
  SELECT gs.d::date AS d
  FROM generate_series($1::date, $2::date, interval '1 day') AS gs(d)
),
total AS (
  SELECT COUNT(*)::float AS cnt FROM students s WHERE s.school_id = $3
),
present_daily AS (
  SELECT ar.date AS d,
         COUNT(DISTINCT ar.student_id) AS c
  FROM attendance_records ar
  JOIN students s ON s.student_id = ar.student_id
  WHERE ar.status = 'present'
    AND s.school_id = $3
    AND ar.date BETWEEN $1::date AND $2::date
  GROUP BY ar.date
)
SELECT days.d,
       COALESCE(p.c, 0) / NULLIF((SELECT cnt FROM total), 0) * 100
FROM days
LEFT JOIN present_daily p ON p.d = days.d
ORDER BY days.d;
`
		rows, err := r.pool.Query(ctx, q, from, to, *schoolID)
		if err != nil {
			return nil, err
		}
		return scanDayFloatRows(rows)
	}
	const q = `
WITH days AS (
  SELECT gs.d::date AS d
  FROM generate_series($1::date, $2::date, interval '1 day') AS gs(d)
),
total AS (
  SELECT COUNT(*)::float AS cnt FROM students
),
present_daily AS (
  SELECT ar.date AS d,
         COUNT(DISTINCT ar.student_id) AS c
  FROM attendance_records ar
  WHERE ar.status = 'present'
    AND ar.date BETWEEN $1::date AND $2::date
  GROUP BY ar.date
)
SELECT days.d,
       COALESCE(p.c, 0) / NULLIF((SELECT cnt FROM total), 0) * 100
FROM days
LEFT JOIN present_daily p ON p.d = days.d
ORDER BY days.d;
`
	rows, err := r.pool.Query(ctx, q, from, to)
	if err != nil {
		return nil, err
	}
	return scanDayFloatRows(rows)
}

// AdminHealthAlertsDaily đếm log severity watch/urgent theo ngày (UTC date của recorded_at).
func (r *AnalyticsTimeseriesRepo) AdminHealthAlertsDaily(ctx context.Context, schoolID *uuid.UUID, fromDate, toDate time.Time) ([]DayInt, error) {
	from := fromDate.UTC().Truncate(24 * time.Hour)
	to := toDate.UTC().Truncate(24 * time.Hour)
	if schoolID != nil {
		const q = `
WITH days AS (
  SELECT gs.d::date AS d
  FROM generate_series($1::date, $2::date, interval '1 day') AS gs(d)
),
alerts AS (
  SELECT (h.recorded_at AT TIME ZONE 'UTC')::date AS d,
         COUNT(*)::int AS c
  FROM health_logs h
  JOIN students s ON s.student_id = h.student_id
  WHERE h.severity IN ('watch', 'urgent')
    AND s.school_id = $3
    AND (h.recorded_at AT TIME ZONE 'UTC')::date BETWEEN $1::date AND $2::date
  GROUP BY 1
)
SELECT days.d, COALESCE(alerts.c, 0)::int
FROM days
LEFT JOIN alerts ON alerts.d = days.d
ORDER BY days.d;
`
		rows, err := r.pool.Query(ctx, q, from, to, *schoolID)
		if err != nil {
			return nil, err
		}
		return scanDayIntRows(rows)
	}
	const q = `
WITH days AS (
  SELECT gs.d::date AS d
  FROM generate_series($1::date, $2::date, interval '1 day') AS gs(d)
),
alerts AS (
  SELECT (h.recorded_at AT TIME ZONE 'UTC')::date AS d,
         COUNT(*)::int AS c
  FROM health_logs h
  WHERE h.severity IN ('watch', 'urgent')
    AND (h.recorded_at AT TIME ZONE 'UTC')::date BETWEEN $1::date AND $2::date
  GROUP BY 1
)
SELECT days.d, COALESCE(alerts.c, 0)::int
FROM days
LEFT JOIN alerts ON alerts.d = days.d
ORDER BY days.d;
`
	rows, err := r.pool.Query(ctx, q, from, to)
	if err != nil {
		return nil, err
	}
	return scanDayIntRows(rows)
}

// TeacherAttendanceMarkedDaily đếm học sinh đã có bản ghi điểm danh trong ngày.
func (r *AnalyticsTimeseriesRepo) TeacherAttendanceMarkedDaily(ctx context.Context, teacherUserID uuid.UUID, fromDate, toDate time.Time) ([]DayInt, error) {
	from := fromDate.UTC().Truncate(24 * time.Hour)
	to := toDate.UTC().Truncate(24 * time.Hour)
	const q = `
WITH days AS (
  SELECT gs.d::date AS d
  FROM generate_series($2::date, $3::date, interval '1 day') AS gs(d)
),
marked AS (
  SELECT ar.date AS d,
         COUNT(DISTINCT ar.student_id)::int AS c
  FROM attendance_records ar
  JOIN students s ON s.student_id = ar.student_id
  JOIN teacher_classes tc ON tc.class_id = s.current_class_id
  JOIN teachers t ON t.teacher_id = tc.teacher_id
  WHERE t.user_id = $1
    AND ar.date BETWEEN $2::date AND $3::date
  GROUP BY ar.date
)
SELECT days.d, COALESCE(marked.c, 0)::int
FROM days
LEFT JOIN marked ON marked.d = days.d
ORDER BY days.d;
`
	rows, err := r.pool.Query(ctx, q, teacherUserID, from, to)
	if err != nil {
		return nil, err
	}
	return scanDayIntRows(rows)
}

// TeacherHealthAlertsDaily cảnh báo sức khỏe (watch/urgent) trong scope lớp của giáo viên.
func (r *AnalyticsTimeseriesRepo) TeacherHealthAlertsDaily(ctx context.Context, teacherUserID uuid.UUID, fromDate, toDate time.Time) ([]DayInt, error) {
	from := fromDate.UTC().Truncate(24 * time.Hour)
	to := toDate.UTC().Truncate(24 * time.Hour)
	const q = `
WITH days AS (
  SELECT gs.d::date AS d
  FROM generate_series($2::date, $3::date, interval '1 day') AS gs(d)
),
alerts AS (
  SELECT (h.recorded_at AT TIME ZONE 'UTC')::date AS d,
         COUNT(*)::int AS c
  FROM health_logs h
  JOIN students s ON s.student_id = h.student_id
  JOIN teacher_classes tc ON tc.class_id = s.current_class_id
  JOIN teachers t ON t.teacher_id = tc.teacher_id
  WHERE t.user_id = $1
    AND h.severity IN ('watch', 'urgent')
    AND (h.recorded_at AT TIME ZONE 'UTC')::date BETWEEN $2::date AND $3::date
  GROUP BY 1
)
SELECT days.d, COALESCE(alerts.c, 0)::int
FROM days
LEFT JOIN alerts ON alerts.d = days.d
ORDER BY days.d;
`
	rows, err := r.pool.Query(ctx, q, teacherUserID, from, to)
	if err != nil {
		return nil, err
	}
	return scanDayIntRows(rows)
}

// TeacherAppointmentsByStatusDay đếm lịch hẹn theo ngày (UTC date của start_time slot) và status.
func (r *AnalyticsTimeseriesRepo) TeacherAppointmentsByStatusDay(ctx context.Context, teacherUserID uuid.UUID, fromDate, toDate time.Time) ([]StatusDayCount, error) {
	from := fromDate.UTC().Truncate(24 * time.Hour)
	to := toDate.UTC().Truncate(24 * time.Hour)
	const q = `
SELECT (s.start_time AT TIME ZONE 'UTC')::date AS d,
       a.status,
       COUNT(*)::int
FROM appointments a
JOIN appointment_slots s ON s.slot_id = a.slot_id
JOIN teachers t ON t.teacher_id = s.teacher_id
WHERE t.user_id = $1
  AND (s.start_time AT TIME ZONE 'UTC')::date BETWEEN $2::date AND $3::date
GROUP BY 1, 2
ORDER BY 1, 2;
`
	rows, err := r.pool.Query(ctx, q, teacherUserID, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []StatusDayCount
	for rows.Next() {
		var d time.Time
		var st string
		var c int
		if err := rows.Scan(&d, &st, &c); err != nil {
			return nil, err
		}
		out = append(out, StatusDayCount{Day: normalizeDayUTC(d), Status: st, Count: c})
	}
	return out, rows.Err()
}

// ParentAttendanceByStatusDaily một học sinh — một dòng mỗi ngày có điểm danh.
func (r *AnalyticsTimeseriesRepo) ParentAttendanceByStatusDaily(ctx context.Context, studentID uuid.UUID, fromDate, toDate time.Time) ([]StatusDayCount, error) {
	from := fromDate.UTC().Truncate(24 * time.Hour)
	to := toDate.UTC().Truncate(24 * time.Hour)
	const q = `
SELECT ar.date AS d,
       ar.status,
       1
FROM attendance_records ar
WHERE ar.student_id = $1
  AND ar.date BETWEEN $2::date AND $3::date;
`
	rows, err := r.pool.Query(ctx, q, studentID, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []StatusDayCount
	for rows.Next() {
		var d time.Time
		var st string
		var c int
		if err := rows.Scan(&d, &st, &c); err != nil {
			return nil, err
		}
		out = append(out, StatusDayCount{Day: normalizeDayUTC(d), Status: st, Count: c})
	}
	return out, rows.Err()
}

// ParentHealthAlertsDaily cảnh báo sức khỏe cho một học sinh.
func (r *AnalyticsTimeseriesRepo) ParentHealthAlertsDaily(ctx context.Context, studentID uuid.UUID, fromDate, toDate time.Time) ([]DayInt, error) {
	from := fromDate.UTC().Truncate(24 * time.Hour)
	to := toDate.UTC().Truncate(24 * time.Hour)
	const q = `
WITH days AS (
  SELECT gs.d::date AS d
  FROM generate_series($2::date, $3::date, interval '1 day') AS gs(d)
),
alerts AS (
  SELECT (h.recorded_at AT TIME ZONE 'UTC')::date AS d,
         COUNT(*)::int AS c
  FROM health_logs h
  WHERE h.student_id = $1
    AND h.severity IN ('watch', 'urgent')
    AND (h.recorded_at AT TIME ZONE 'UTC')::date BETWEEN $2::date AND $3::date
  GROUP BY 1
)
SELECT days.d, COALESCE(alerts.c, 0)::int
FROM days
LEFT JOIN alerts ON alerts.d = days.d
ORDER BY days.d;
`
	rows, err := r.pool.Query(ctx, q, studentID, from, to)
	if err != nil {
		return nil, err
	}
	return scanDayIntRows(rows)
}

// ParentAppointmentsByStatusDay lịch hẹn của phụ huynh cho một học sinh.
func (r *AnalyticsTimeseriesRepo) ParentAppointmentsByStatusDay(ctx context.Context, parentUserID, studentID uuid.UUID, fromDate, toDate time.Time) ([]StatusDayCount, error) {
	from := fromDate.UTC().Truncate(24 * time.Hour)
	to := toDate.UTC().Truncate(24 * time.Hour)
	const q = `
SELECT (s.start_time AT TIME ZONE 'UTC')::date AS d,
       a.status,
       COUNT(*)::int
FROM appointments a
JOIN appointment_slots s ON s.slot_id = a.slot_id
JOIN parents p ON p.parent_id = a.parent_id
WHERE p.user_id = $1
  AND a.student_id = $2
  AND (s.start_time AT TIME ZONE 'UTC')::date BETWEEN $3::date AND $4::date
GROUP BY 1, 2
ORDER BY 1, 2;
`
	rows, err := r.pool.Query(ctx, q, parentUserID, studentID, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []StatusDayCount
	for rows.Next() {
		var d time.Time
		var st string
		var c int
		if err := rows.Scan(&d, &st, &c); err != nil {
			return nil, err
		}
		out = append(out, StatusDayCount{Day: normalizeDayUTC(d), Status: st, Count: c})
	}
	return out, rows.Err()
}
