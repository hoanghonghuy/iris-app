package repo

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type appointmentRepoFixture struct {
	schoolID          uuid.UUID
	classID           uuid.UUID
	teacherID         uuid.UUID
	teacherUserID     uuid.UUID
	parentID          uuid.UUID
	parentUserID      uuid.UUID
	otherParentUserID uuid.UUID
	studentID         uuid.UUID
}

func mustNewAppointmentTestPool(t *testing.T) *pgxpool.Pool {
	t.Helper()
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		t.Skip("DATABASE_URL is not set; skipping appointment repo integration tests")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		t.Fatalf("pgxpool.New() error = %v", err)
	}
	t.Cleanup(pool.Close)

	if err := pool.Ping(ctx); err != nil {
		t.Fatalf("pool.Ping() error = %v", err)
	}

	return pool
}

func seedAppointmentFixture(t *testing.T, pool *pgxpool.Pool) appointmentRepoFixture {
	t.Helper()
	ctx := context.Background()

	fx := appointmentRepoFixture{
		schoolID:          uuid.New(),
		classID:           uuid.New(),
		teacherID:         uuid.New(),
		teacherUserID:     uuid.New(),
		parentID:          uuid.New(),
		parentUserID:      uuid.New(),
		otherParentUserID: uuid.New(),
		studentID:         uuid.New(),
	}

	teacherEmail := fmt.Sprintf("teacher-%s@example.com", uuid.NewString())
	parentEmail := fmt.Sprintf("parent-%s@example.com", uuid.NewString())
	otherParentEmail := fmt.Sprintf("parent-other-%s@example.com", uuid.NewString())

	if _, err := pool.Exec(ctx, `INSERT INTO schools (school_id, name) VALUES ($1, $2)`, fx.schoolID, "Test School"); err != nil {
		t.Fatalf("insert school error = %v", err)
	}
	if _, err := pool.Exec(ctx, `INSERT INTO classes (class_id, school_id, name, school_year) VALUES ($1, $2, $3, $4)`, fx.classID, fx.schoolID, "A1", "2026"); err != nil {
		t.Fatalf("insert class error = %v", err)
	}

	if _, err := pool.Exec(ctx, `INSERT INTO users (user_id, email, password_hash) VALUES ($1, $2, $3), ($4, $5, $6), ($7, $8, $9)`,
		fx.teacherUserID, teacherEmail, "hash",
		fx.parentUserID, parentEmail, "hash",
		fx.otherParentUserID, otherParentEmail, "hash",
	); err != nil {
		t.Fatalf("insert users error = %v", err)
	}

	if _, err := pool.Exec(ctx, `INSERT INTO teachers (teacher_id, user_id, school_id, full_name) VALUES ($1, $2, $3, $4)`, fx.teacherID, fx.teacherUserID, fx.schoolID, "Teacher A"); err != nil {
		t.Fatalf("insert teacher error = %v", err)
	}
	if _, err := pool.Exec(ctx, `INSERT INTO parents (parent_id, user_id, school_id, full_name) VALUES ($1, $2, $3, $4)`, fx.parentID, fx.parentUserID, fx.schoolID, "Parent A"); err != nil {
		t.Fatalf("insert parent error = %v", err)
	}
	if _, err := pool.Exec(ctx, `INSERT INTO parents (parent_id, user_id, school_id, full_name) VALUES ($1, $2, $3, $4)`, uuid.New(), fx.otherParentUserID, fx.schoolID, "Parent B"); err != nil {
		t.Fatalf("insert other parent error = %v", err)
	}

	if _, err := pool.Exec(ctx, `INSERT INTO students (student_id, school_id, current_class_id, full_name) VALUES ($1, $2, $3, $4)`, fx.studentID, fx.schoolID, fx.classID, "Student A"); err != nil {
		t.Fatalf("insert student error = %v", err)
	}

	if _, err := pool.Exec(ctx, `INSERT INTO teacher_classes (teacher_id, class_id) VALUES ($1, $2)`, fx.teacherID, fx.classID); err != nil {
		t.Fatalf("insert teacher_classes error = %v", err)
	}
	if _, err := pool.Exec(ctx, `INSERT INTO student_parents (student_id, parent_id, relationship) VALUES ($1, $2, $3)`, fx.studentID, fx.parentID, "father"); err != nil {
		t.Fatalf("insert student_parents error = %v", err)
	}

	t.Cleanup(func() {
		cleanupCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_, _ = pool.Exec(cleanupCtx, `DELETE FROM users WHERE user_id = ANY($1::uuid[])`, []uuid.UUID{fx.teacherUserID, fx.parentUserID, fx.otherParentUserID})
		_, _ = pool.Exec(cleanupCtx, `DELETE FROM schools WHERE school_id = $1`, fx.schoolID)
	})

	return fx
}

func insertAppointmentSlotForFixture(t *testing.T, pool *pgxpool.Pool, teacherID, classID uuid.UUID) uuid.UUID {
	t.Helper()
	ctx := context.Background()
	var slotID uuid.UUID
	start := time.Now().UTC().Add(24 * time.Hour)
	end := start.Add(30 * time.Minute)
	if err := pool.QueryRow(ctx, `
		INSERT INTO appointment_slots (teacher_id, class_id, start_time, end_time, note)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING slot_id
	`, teacherID, classID, start, end, "integration-slot").Scan(&slotID); err != nil {
		t.Fatalf("insert appointment slot error = %v", err)
	}
	return slotID
}

func insertAppointmentSlotForFixtureAt(t *testing.T, pool *pgxpool.Pool, teacherID, classID uuid.UUID, start, end time.Time, note string) uuid.UUID {
	t.Helper()
	ctx := context.Background()
	var slotID uuid.UUID
	if err := pool.QueryRow(ctx, `
		INSERT INTO appointment_slots (teacher_id, class_id, start_time, end_time, note)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING slot_id
	`, teacherID, classID, start, end, note).Scan(&slotID); err != nil {
		t.Fatalf("insert appointment slot at time error = %v", err)
	}
	return slotID
}

func TestAppointmentRepoCreateAppointment_TransactionConflict(t *testing.T) {
	pool := mustNewAppointmentTestPool(t)
	fx := seedAppointmentFixture(t, pool)
	repo := NewAppointmentRepo(pool)

	slotID := insertAppointmentSlotForFixture(t, pool, fx.teacherID, fx.classID)

	created, err := repo.CreateAppointment(context.Background(), fx.parentUserID, fx.studentID, slotID, "first")
	if err != nil {
		t.Fatalf("CreateAppointment() first call error = %v", err)
	}
	if created.Status != "pending" {
		t.Fatalf("status = %q, want %q", created.Status, "pending")
	}

	_, err = repo.CreateAppointment(context.Background(), fx.parentUserID, fx.studentID, slotID, "second")
	if !errors.Is(err, ErrAppointmentSlotUnavailable) {
		t.Fatalf("error = %v, want %v", err, ErrAppointmentSlotUnavailable)
	}

	var count int
	if err := pool.QueryRow(context.Background(), `SELECT COUNT(*) FROM appointments WHERE slot_id = $1`, slotID).Scan(&count); err != nil {
		t.Fatalf("count appointments error = %v", err)
	}
	if count != 1 {
		t.Fatalf("appointments count = %d, want %d", count, 1)
	}
}

func TestAppointmentRepoCreateAppointment_RejectsOutOfScopeParent(t *testing.T) {
	pool := mustNewAppointmentTestPool(t)
	fx := seedAppointmentFixture(t, pool)
	repo := NewAppointmentRepo(pool)

	slotID := insertAppointmentSlotForFixture(t, pool, fx.teacherID, fx.classID)

	_, err := repo.CreateAppointment(context.Background(), fx.otherParentUserID, fx.studentID, slotID, "x")
	if !errors.Is(err, ErrNoRowsUpdated) {
		t.Fatalf("error = %v, want %v", err, ErrNoRowsUpdated)
	}
}

func TestAppointmentRepoUpdateAndCancelAppointment_OwnershipAndStatusFlow(t *testing.T) {
	pool := mustNewAppointmentTestPool(t)
	fx := seedAppointmentFixture(t, pool)
	repo := NewAppointmentRepo(pool)

	slotID := insertAppointmentSlotForFixture(t, pool, fx.teacherID, fx.classID)
	created, err := repo.CreateAppointment(context.Background(), fx.parentUserID, fx.studentID, slotID, "x")
	if err != nil {
		t.Fatalf("CreateAppointment() error = %v", err)
	}

	updated, err := repo.UpdateAppointmentStatusByTeacher(context.Background(), fx.teacherUserID, created.AppointmentID, "confirmed", "")
	if err != nil {
		t.Fatalf("UpdateAppointmentStatusByTeacher() error = %v", err)
	}
	if updated.Status != "confirmed" {
		t.Fatalf("status = %q, want %q", updated.Status, "confirmed")
	}

	_, err = repo.UpdateAppointmentStatusByTeacher(context.Background(), uuid.New(), created.AppointmentID, "completed", "")
	if !errors.Is(err, ErrNoRowsUpdated) {
		t.Fatalf("error = %v, want %v", err, ErrNoRowsUpdated)
	}

	cancelled, err := repo.CancelAppointmentByParent(context.Background(), fx.parentUserID, created.AppointmentID, "need cancel")
	if err != nil {
		t.Fatalf("CancelAppointmentByParent() error = %v", err)
	}
	if cancelled.Status != "cancelled" {
		t.Fatalf("status = %q, want %q", cancelled.Status, "cancelled")
	}
	if cancelled.CancelReason != "need cancel" {
		t.Fatalf("cancel_reason = %q, want %q", cancelled.CancelReason, "need cancel")
	}

	_, err = repo.CancelAppointmentByParent(context.Background(), fx.parentUserID, created.AppointmentID, "again")
	if !errors.Is(err, ErrNoRowsUpdated) {
		t.Fatalf("error = %v, want %v", err, ErrNoRowsUpdated)
	}
}

func TestAppointmentRepoListAppointments_FilterPaginationAndScope(t *testing.T) {
	pool := mustNewAppointmentTestPool(t)
	fx := seedAppointmentFixture(t, pool)
	repo := NewAppointmentRepo(pool)
	ctx := context.Background()

	base := time.Now().UTC().Add(24 * time.Hour).Truncate(time.Second)
	slot1 := insertAppointmentSlotForFixtureAt(t, pool, fx.teacherID, fx.classID, base.Add(1*time.Hour), base.Add(90*time.Minute), "slot-1")
	slot2 := insertAppointmentSlotForFixtureAt(t, pool, fx.teacherID, fx.classID, base.Add(2*time.Hour), base.Add(150*time.Minute), "slot-2")
	slot3 := insertAppointmentSlotForFixtureAt(t, pool, fx.teacherID, fx.classID, base.Add(3*time.Hour), base.Add(210*time.Minute), "slot-3")

	a1, err := repo.CreateAppointment(ctx, fx.parentUserID, fx.studentID, slot1, "a1")
	if err != nil {
		t.Fatalf("CreateAppointment(slot1) error = %v", err)
	}
	a2, err := repo.CreateAppointment(ctx, fx.parentUserID, fx.studentID, slot2, "a2")
	if err != nil {
		t.Fatalf("CreateAppointment(slot2) error = %v", err)
	}
	a3, err := repo.CreateAppointment(ctx, fx.parentUserID, fx.studentID, slot3, "a3")
	if err != nil {
		t.Fatalf("CreateAppointment(slot3) error = %v", err)
	}

	if _, err := repo.UpdateAppointmentStatusByTeacher(ctx, fx.teacherUserID, a2.AppointmentID, "confirmed", ""); err != nil {
		t.Fatalf("UpdateAppointmentStatusByTeacher() error = %v", err)
	}
	if _, err := repo.CancelAppointmentByParent(ctx, fx.parentUserID, a3.AppointmentID, "cancelled-by-parent"); err != nil {
		t.Fatalf("CancelAppointmentByParent() error = %v", err)
	}

	teacherItems, teacherTotal, err := repo.ListTeacherAppointments(ctx, fx.teacherUserID, "", nil, nil, 10, 0)
	if err != nil {
		t.Fatalf("ListTeacherAppointments(all) error = %v", err)
	}
	if teacherTotal != 3 {
		t.Fatalf("teacher total = %d, want %d", teacherTotal, 3)
	}
	if len(teacherItems) != 3 {
		t.Fatalf("teacher items len = %d, want %d", len(teacherItems), 3)
	}
	if teacherItems[0].AppointmentID != a1.AppointmentID || teacherItems[1].AppointmentID != a2.AppointmentID || teacherItems[2].AppointmentID != a3.AppointmentID {
		t.Fatalf("teacher items order by start_time is not as expected")
	}

	confirmedItems, confirmedTotal, err := repo.ListTeacherAppointments(ctx, fx.teacherUserID, "confirmed", nil, nil, 10, 0)
	if err != nil {
		t.Fatalf("ListTeacherAppointments(confirmed) error = %v", err)
	}
	if confirmedTotal != 1 || len(confirmedItems) != 1 || confirmedItems[0].AppointmentID != a2.AppointmentID {
		t.Fatalf("confirmed filter mismatch: total=%d len=%d", confirmedTotal, len(confirmedItems))
	}

	from := base.Add(2 * time.Hour)
	fromItems, fromTotal, err := repo.ListTeacherAppointments(ctx, fx.teacherUserID, "", &from, nil, 10, 0)
	if err != nil {
		t.Fatalf("ListTeacherAppointments(from) error = %v", err)
	}
	if fromTotal != 2 || len(fromItems) != 2 {
		t.Fatalf("from filter mismatch: total=%d len=%d", fromTotal, len(fromItems))
	}
	if fromItems[0].AppointmentID != a2.AppointmentID || fromItems[1].AppointmentID != a3.AppointmentID {
		t.Fatalf("from filter order mismatch")
	}

	pagedItems, pagedTotal, err := repo.ListTeacherAppointments(ctx, fx.teacherUserID, "", nil, nil, 1, 1)
	if err != nil {
		t.Fatalf("ListTeacherAppointments(pagination) error = %v", err)
	}
	if pagedTotal != 3 || len(pagedItems) != 1 || pagedItems[0].AppointmentID != a2.AppointmentID {
		t.Fatalf("pagination mismatch: total=%d len=%d", pagedTotal, len(pagedItems))
	}

	parentPendingItems, parentPendingTotal, err := repo.ListParentAppointments(ctx, fx.parentUserID, "pending", nil, nil, 10, 0)
	if err != nil {
		t.Fatalf("ListParentAppointments(pending) error = %v", err)
	}
	if parentPendingTotal != 1 || len(parentPendingItems) != 1 || parentPendingItems[0].AppointmentID != a1.AppointmentID {
		t.Fatalf("parent pending filter mismatch: total=%d len=%d", parentPendingTotal, len(parentPendingItems))
	}

	otherParentItems, otherParentTotal, err := repo.ListParentAppointments(ctx, fx.otherParentUserID, "", nil, nil, 10, 0)
	if err != nil {
		t.Fatalf("ListParentAppointments(other parent) error = %v", err)
	}
	if otherParentTotal != 0 || len(otherParentItems) != 0 {
		t.Fatalf("expected no appointments for other parent, got total=%d len=%d", otherParentTotal, len(otherParentItems))
	}
}

func TestAppointmentRepoListAvailableSlotsForParent_FilterPaginationAndScope(t *testing.T) {
	pool := mustNewAppointmentTestPool(t)
	fx := seedAppointmentFixture(t, pool)
	repo := NewAppointmentRepo(pool)
	ctx := context.Background()

	base := time.Now().UTC().Add(24 * time.Hour).Truncate(time.Second)
	slot1 := insertAppointmentSlotForFixtureAt(t, pool, fx.teacherID, fx.classID, base.Add(1*time.Hour), base.Add(90*time.Minute), "slot-1")
	slot2 := insertAppointmentSlotForFixtureAt(t, pool, fx.teacherID, fx.classID, base.Add(2*time.Hour), base.Add(150*time.Minute), "slot-2")
	slot3 := insertAppointmentSlotForFixtureAt(t, pool, fx.teacherID, fx.classID, base.Add(3*time.Hour), base.Add(210*time.Minute), "slot-3")
	slot4 := insertAppointmentSlotForFixtureAt(t, pool, fx.teacherID, fx.classID, base.Add(4*time.Hour), base.Add(270*time.Minute), "slot-4")

	if _, err := repo.CreateAppointment(ctx, fx.parentUserID, fx.studentID, slot2, "booked"); err != nil {
		t.Fatalf("CreateAppointment(slot2) error = %v", err)
	}
	if _, err := pool.Exec(ctx, `UPDATE appointment_slots SET is_active = false WHERE slot_id = $1`, slot4); err != nil {
		t.Fatalf("deactivate slot4 error = %v", err)
	}

	allItems, allTotal, err := repo.ListAvailableSlotsForParent(ctx, fx.parentUserID, fx.studentID, nil, nil, 10, 0)
	if err != nil {
		t.Fatalf("ListAvailableSlotsForParent(all) error = %v", err)
	}
	if allTotal != 2 || len(allItems) != 2 {
		t.Fatalf("all available mismatch: total=%d len=%d", allTotal, len(allItems))
	}
	if allItems[0].SlotID != slot1 || allItems[1].SlotID != slot3 {
		t.Fatalf("available slots order mismatch")
	}

	pagedItems, pagedTotal, err := repo.ListAvailableSlotsForParent(ctx, fx.parentUserID, fx.studentID, nil, nil, 1, 1)
	if err != nil {
		t.Fatalf("ListAvailableSlotsForParent(pagination) error = %v", err)
	}
	if pagedTotal != 2 || len(pagedItems) != 1 || pagedItems[0].SlotID != slot3 {
		t.Fatalf("available slots pagination mismatch: total=%d len=%d", pagedTotal, len(pagedItems))
	}

	from := base.Add(150 * time.Minute)
	to := base.Add(220 * time.Minute)
	windowItems, windowTotal, err := repo.ListAvailableSlotsForParent(ctx, fx.parentUserID, fx.studentID, &from, &to, 10, 0)
	if err != nil {
		t.Fatalf("ListAvailableSlotsForParent(time window) error = %v", err)
	}
	if windowTotal != 1 || len(windowItems) != 1 || windowItems[0].SlotID != slot3 {
		t.Fatalf("available slots time window mismatch: total=%d len=%d", windowTotal, len(windowItems))
	}

	otherParentItems, otherParentTotal, err := repo.ListAvailableSlotsForParent(ctx, fx.otherParentUserID, fx.studentID, nil, nil, 10, 0)
	if err != nil {
		t.Fatalf("ListAvailableSlotsForParent(other parent) error = %v", err)
	}
	if otherParentTotal != 0 || len(otherParentItems) != 0 {
		t.Fatalf("expected no available slots for other parent, got total=%d len=%d", otherParentTotal, len(otherParentItems))
	}
}
