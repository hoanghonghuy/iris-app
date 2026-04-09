package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/hoanghonghuy/iris-app/apps/api/internal/model"
)

type AppointmentRepo struct {
	pool *pgxpool.Pool
}

type scanableRow interface {
	Scan(...any) error
}

func NewAppointmentRepo(pool *pgxpool.Pool) *AppointmentRepo {
	return &AppointmentRepo{pool: pool}
}

// scanAppointment gom logic Scan cho các mutation trả về cùng cấu trúc appointments.
func scanAppointment(row scanableRow) (model.Appointment, error) {
	var a model.Appointment
	err := row.Scan(
		&a.AppointmentID,
		&a.SlotID,
		&a.ParentID,
		&a.StudentID,
		&a.Status,
		&a.Note,
		&a.CancelReason,
		&a.ConfirmedAt,
		&a.CompletedAt,
		&a.CancelledAt,
		&a.CreatedAt,
		&a.UpdatedAt,
	)
	if err != nil {
		return model.Appointment{}, err
	}

	return a, nil
}

// listAppointmentsByUser dùng chung phần query/filter/paging cho teacher và parent,
// chỉ khác điều kiện WHERE theo user.
func (r *AppointmentRepo) listAppointmentsByUser(
	ctx context.Context,
	userWhereSQL string,
	userID uuid.UUID,
	status string,
	from, to *time.Time,
	limit, offset int,
) ([]model.Appointment, int, error) {
	q := `
		SELECT a.appointment_id, a.slot_id, a.parent_id, p.full_name, a.student_id, st.full_name,
		       s.teacher_id, t.full_name, s.class_id, c.name,
		       a.status, COALESCE(a.note, ''), COALESCE(a.cancel_reason, ''), s.start_time, s.end_time,
		       a.confirmed_at, a.completed_at, a.cancelled_at, a.created_at, a.updated_at,
		       COUNT(*) OVER() AS total_count
		FROM appointments a
		JOIN appointment_slots s ON s.slot_id = a.slot_id
		JOIN teachers t ON t.teacher_id = s.teacher_id
		JOIN classes c ON c.class_id = s.class_id
		JOIN parents p ON p.parent_id = a.parent_id
		JOIN students st ON st.student_id = a.student_id
		WHERE ` + userWhereSQL

	args := []any{userID}
	argPos := 2
	if status != "" {
		q += fmt.Sprintf(" AND a.status = $%d", argPos)
		args = append(args, status)
		argPos++
	}
	if from != nil {
		q += fmt.Sprintf(" AND s.start_time >= $%d", argPos)
		args = append(args, *from)
		argPos++
	}
	if to != nil {
		q += fmt.Sprintf(" AND s.start_time <= $%d", argPos)
		args = append(args, *to)
		argPos++
	}
	q += fmt.Sprintf(" ORDER BY s.start_time ASC LIMIT $%d OFFSET $%d", argPos, argPos+1)
	args = append(args, limit, offset)

	rows, err := r.pool.Query(ctx, q, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	items := make([]model.Appointment, 0)
	total := 0
	for rows.Next() {
		var a model.Appointment
		if err := rows.Scan(
			&a.AppointmentID, &a.SlotID, &a.ParentID, &a.ParentName, &a.StudentID, &a.StudentName,
			&a.TeacherID, &a.TeacherName, &a.ClassID, &a.ClassName,
			&a.Status, &a.Note, &a.CancelReason, &a.StartTime, &a.EndTime,
			&a.ConfirmedAt, &a.CompletedAt, &a.CancelledAt, &a.CreatedAt, &a.UpdatedAt, &total,
		); err != nil {
			return nil, 0, err
		}
		items = append(items, a)
	}

	return items, total, rows.Err()
}

func (r *AppointmentRepo) CreateSlot(ctx context.Context, teacherUserID, classID uuid.UUID, startTime, endTime time.Time, note string) (model.AppointmentSlot, error) {
	const q = `
		INSERT INTO appointment_slots (teacher_id, class_id, start_time, end_time, note)
		SELECT t.teacher_id, $2, $3, $4, $5
		FROM teachers t
		JOIN teacher_classes tc ON tc.teacher_id = t.teacher_id
		WHERE t.user_id = $1 AND tc.class_id = $2
		RETURNING slot_id, teacher_id, class_id, start_time, end_time, note, is_active, created_at, updated_at;
	`

	var slot model.AppointmentSlot
	err := r.pool.QueryRow(ctx, q, teacherUserID, classID, startTime, endTime, note).Scan(
		&slot.SlotID,
		&slot.TeacherID,
		&slot.ClassID,
		&slot.StartTime,
		&slot.EndTime,
		&slot.Note,
		&slot.IsActive,
		&slot.CreatedAt,
		&slot.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return model.AppointmentSlot{}, ErrNoRowsUpdated
		}
		return model.AppointmentSlot{}, err
	}

	return slot, nil
}

func (r *AppointmentRepo) CreateAppointment(ctx context.Context, parentUserID, studentID, slotID uuid.UUID, note string) (model.Appointment, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return model.Appointment{}, err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	const lockEligibleSlotQ = `
		SELECT p.parent_id
		FROM parents p
		JOIN student_parents sp ON sp.parent_id = p.parent_id
		JOIN students st ON st.student_id = sp.student_id
		JOIN appointment_slots s ON s.slot_id = $3
		WHERE p.user_id = $1
		  AND sp.student_id = $2
		  AND st.current_class_id = s.class_id
		  AND s.is_active = true
		  AND s.start_time >= now()
		FOR UPDATE OF s;
	`

	var parentID uuid.UUID
	err = tx.QueryRow(ctx, lockEligibleSlotQ, parentUserID, studentID, slotID).Scan(&parentID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return model.Appointment{}, ErrNoRowsUpdated
		}
		return model.Appointment{}, err
	}

	const activeAppointmentConflictQ = `
		SELECT 1
		FROM appointments
		WHERE slot_id = $1
		  AND status IN ('pending', 'confirmed', 'completed', 'no_show')
		LIMIT 1;
	`

	var activeConflict int
	err = tx.QueryRow(ctx, activeAppointmentConflictQ, slotID).Scan(&activeConflict)
	if err == nil {
		return model.Appointment{}, ErrAppointmentSlotUnavailable
	}
	if err != pgx.ErrNoRows {
		return model.Appointment{}, err
	}

	const createAppointmentQ = `
		INSERT INTO appointments (slot_id, parent_id, student_id, status, note)
		VALUES ($1, $2, $3, 'pending', $4)
		RETURNING appointment_id, slot_id, parent_id, student_id, status, COALESCE(note, ''), COALESCE(cancel_reason, ''),
		          confirmed_at, completed_at, cancelled_at, created_at, updated_at;
	`

	a, err := scanAppointment(tx.QueryRow(ctx, createAppointmentQ, slotID, parentID, studentID, note))
	if err != nil {
		return model.Appointment{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return model.Appointment{}, err
	}

	return a, nil
}

func (r *AppointmentRepo) UpdateAppointmentStatusByTeacher(ctx context.Context, teacherUserID, appointmentID uuid.UUID, status, cancelReason string) (model.Appointment, error) {
	const q = `
		UPDATE appointments a
		SET status = $3,
		    cancel_reason = CASE WHEN $3 = 'cancelled' THEN $4 ELSE a.cancel_reason END,
		    confirmed_at = CASE WHEN $3 = 'confirmed' THEN now() ELSE a.confirmed_at END,
		    completed_at = CASE WHEN $3 = 'completed' THEN now() ELSE a.completed_at END,
		    cancelled_at = CASE WHEN $3 = 'cancelled' THEN now() ELSE a.cancelled_at END,
		    updated_at = now()
		FROM appointment_slots s
		JOIN teachers t ON t.teacher_id = s.teacher_id
		WHERE a.slot_id = s.slot_id
		  AND a.appointment_id = $2
		  AND t.user_id = $1
		RETURNING a.appointment_id, a.slot_id, a.parent_id, a.student_id, a.status, COALESCE(a.note, ''), COALESCE(a.cancel_reason, ''),
		          a.confirmed_at, a.completed_at, a.cancelled_at, a.created_at, a.updated_at;
	`

	a, err := scanAppointment(r.pool.QueryRow(ctx, q, teacherUserID, appointmentID, status, cancelReason))
	if err != nil {
		if err == pgx.ErrNoRows {
			return model.Appointment{}, ErrNoRowsUpdated
		}
		return model.Appointment{}, err
	}

	return a, nil
}

func (r *AppointmentRepo) CancelAppointmentByParent(ctx context.Context, parentUserID, appointmentID uuid.UUID, cancelReason string) (model.Appointment, error) {
	const q = `
		UPDATE appointments a
		SET status = 'cancelled',
		    cancel_reason = $3,
		    cancelled_at = now(),
		    updated_at = now()
		FROM parents p
		WHERE a.parent_id = p.parent_id
		  AND p.user_id = $1
		  AND a.appointment_id = $2
		  AND a.status IN ('pending', 'confirmed')
		RETURNING a.appointment_id, a.slot_id, a.parent_id, a.student_id, a.status, COALESCE(a.note, ''), COALESCE(a.cancel_reason, ''),
		          a.confirmed_at, a.completed_at, a.cancelled_at, a.created_at, a.updated_at;
	`

	a, err := scanAppointment(r.pool.QueryRow(ctx, q, parentUserID, appointmentID, cancelReason))
	if err != nil {
		if err == pgx.ErrNoRows {
			return model.Appointment{}, ErrNoRowsUpdated
		}
		return model.Appointment{}, err
	}

	return a, nil
}

func (r *AppointmentRepo) ListTeacherAppointments(ctx context.Context, teacherUserID uuid.UUID, status string, from, to *time.Time, limit, offset int) ([]model.Appointment, int, error) {
	return r.listAppointmentsByUser(ctx, "t.user_id = $1", teacherUserID, status, from, to, limit, offset)
}

func (r *AppointmentRepo) ListParentAppointments(ctx context.Context, parentUserID uuid.UUID, status string, from, to *time.Time, limit, offset int) ([]model.Appointment, int, error) {
	return r.listAppointmentsByUser(ctx, "p.user_id = $1", parentUserID, status, from, to, limit, offset)
}

func (r *AppointmentRepo) ListAvailableSlotsForParent(ctx context.Context, parentUserID, studentID uuid.UUID, from, to *time.Time, limit, offset int) ([]model.AppointmentSlot, int, error) {
	q := `
		SELECT s.slot_id, s.teacher_id, t.full_name, s.class_id, c.name, s.start_time, s.end_time,
		       s.note, s.is_active, s.created_at, s.updated_at,
		       COUNT(*) OVER() AS total_count
		FROM appointment_slots s
		JOIN classes c ON c.class_id = s.class_id
		JOIN teachers t ON t.teacher_id = s.teacher_id
		JOIN students st ON st.student_id = $2 AND st.current_class_id = s.class_id
		JOIN student_parents sp ON sp.student_id = st.student_id
		JOIN parents p ON p.parent_id = sp.parent_id
		WHERE p.user_id = $1
		  AND s.is_active = true
		  AND s.start_time >= now()
		  AND NOT EXISTS (
			SELECT 1 FROM appointments a WHERE a.slot_id = s.slot_id AND a.status IN ('pending', 'confirmed', 'completed', 'no_show')
		  )
	`

	args := []any{parentUserID, studentID}
	argPos := 3
	if from != nil {
		q += fmt.Sprintf(" AND s.start_time >= $%d", argPos)
		args = append(args, *from)
		argPos++
	}
	if to != nil {
		q += fmt.Sprintf(" AND s.start_time <= $%d", argPos)
		args = append(args, *to)
		argPos++
	}
	q += fmt.Sprintf(" ORDER BY s.start_time ASC LIMIT $%d OFFSET $%d", argPos, argPos+1)
	args = append(args, limit, offset)

	rows, err := r.pool.Query(ctx, q, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	items := make([]model.AppointmentSlot, 0)
	total := 0
	for rows.Next() {
		var s model.AppointmentSlot
		if err := rows.Scan(
			&s.SlotID, &s.TeacherID, &s.TeacherName, &s.ClassID, &s.ClassName, &s.StartTime, &s.EndTime,
			&s.Note, &s.IsActive, &s.CreatedAt, &s.UpdatedAt, &total,
		); err != nil {
			return nil, 0, err
		}
		items = append(items, s)
	}

	return items, total, rows.Err()
}

func (r *AppointmentRepo) CountParentUpcomingAppointments(ctx context.Context, parentUserID uuid.UUID) (int, error) {
	const q = `
		SELECT COUNT(*)
		FROM appointments a
		JOIN parents p ON p.parent_id = a.parent_id
		JOIN appointment_slots s ON s.slot_id = a.slot_id
		WHERE p.user_id = $1
		  AND a.status IN ('pending', 'confirmed')
		  AND s.start_time >= now();
	`
	var count int
	err := r.pool.QueryRow(ctx, q, parentUserID).Scan(&count)
	return count, err
}

func (r *AppointmentRepo) CountTodayPendingAppointmentsBySchool(ctx context.Context, schoolID *uuid.UUID) (int, error) {
	var (
		q     string
		args  []any
		count int
	)

	if schoolID != nil {
		q = `
			SELECT COUNT(*)
			FROM appointments a
			JOIN appointment_slots s ON s.slot_id = a.slot_id
			JOIN classes c ON c.class_id = s.class_id
			WHERE a.status = 'pending'
			  AND s.start_time::date = CURRENT_DATE
			  AND c.school_id = $1;
		`
		args = append(args, *schoolID)
	} else {
		q = `
			SELECT COUNT(*)
			FROM appointments a
			JOIN appointment_slots s ON s.slot_id = a.slot_id
			WHERE a.status = 'pending'
			  AND s.start_time::date = CURRENT_DATE;
		`
	}

	err := r.pool.QueryRow(ctx, q, args...).Scan(&count)
	return count, err
}

func (r *AppointmentRepo) CountTeacherPendingAppointments(ctx context.Context, teacherUserID uuid.UUID) (int, error) {
	const q = `
		SELECT COUNT(*)
		FROM appointments a
		JOIN appointment_slots s ON s.slot_id = a.slot_id
		JOIN teachers t ON t.teacher_id = s.teacher_id
		WHERE t.user_id = $1
		  AND a.status = 'pending'
		  AND s.start_time >= now();
	`
	var count int
	err := r.pool.QueryRow(ctx, q, teacherUserID).Scan(&count)
	return count, err
}
