package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/model"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/repo"
	"github.com/jackc/pgx/v5"
)

type fakeTeacherScopeServiceTeacherScopeRepo struct {
	listMyClassesCalls int
	listMyClassesRes   []model.Class
	listMyClassesErr   error

	listMyStudentsInClassCalls int
	listMyStudentsInClassRes   []model.Student
	listMyStudentsInClassErr   error

	upsertAttendanceCalls int
	upsertAttendanceArg   struct {
		teacherUserID uuid.UUID
		studentID     uuid.UUID
		date          time.Time
		status        string
		note          string
	}
	upsertAttendanceErr error

	deleteAttendanceForDateCalls int
	deleteAttendanceForDateArg   struct {
		teacherUserID uuid.UUID
		studentID     uuid.UUID
		date          time.Time
	}
	deleteAttendanceForDateErr error

	listAttendanceByStudentRes []model.AttendanceRecord
	listAttendanceByStudentErr error

	listAttendanceChangeLogsByStudentRes []model.AttendanceChangeLog
	listAttendanceChangeLogsByStudentErr error

	listAttendanceChangeLogsByClassCalls int
	listAttendanceChangeLogsByClassArg   struct {
		limit  int
		offset int
		status *string
	}
	listAttendanceChangeLogsByClassRes   []model.AttendanceChangeLog
	listAttendanceChangeLogsByClassTotal int
	listAttendanceChangeLogsByClassErr   error

	createClassPostErr error
	createClassPostRes uuid.UUID

	createStudentPostErr error
	createStudentPostRes uuid.UUID

	listClassPostsCalls int
	listClassPostsArg   struct {
		limit  int
		offset int
	}
	listClassPostsRes   []model.Post
	listClassPostsTotal int
	listClassPostsErr   error

	listStudentPostsCalls int
	listStudentPostsArg   struct {
		limit  int
		offset int
	}
	listStudentPostsRes   []model.Post
	listStudentPostsTotal int
	listStudentPostsErr   error

	updatePostErr error
	deletePostErr error

	countMyStudentsRes int
	countMyStudentsErr error
	countMyPostsRes    int
	countMyPostsErr    error
}

func (f *fakeTeacherScopeServiceTeacherScopeRepo) ListMyClasses(_ context.Context, _ uuid.UUID) ([]model.Class, error) {
	f.listMyClassesCalls++
	return f.listMyClassesRes, f.listMyClassesErr
}

func (f *fakeTeacherScopeServiceTeacherScopeRepo) ListMyStudentsInClass(_ context.Context, _, _ uuid.UUID) ([]model.Student, error) {
	f.listMyStudentsInClassCalls++
	return f.listMyStudentsInClassRes, f.listMyStudentsInClassErr
}

func (f *fakeTeacherScopeServiceTeacherScopeRepo) UpsertAttendance(_ context.Context, teacherUserID, studentID uuid.UUID, date time.Time, status string, _ *time.Time, _ *time.Time, note string) error {
	f.upsertAttendanceCalls++
	f.upsertAttendanceArg.teacherUserID = teacherUserID
	f.upsertAttendanceArg.studentID = studentID
	f.upsertAttendanceArg.date = date
	f.upsertAttendanceArg.status = status
	f.upsertAttendanceArg.note = note
	return f.upsertAttendanceErr
}

func (f *fakeTeacherScopeServiceTeacherScopeRepo) DeleteAttendanceForDate(_ context.Context, teacherUserID, studentID uuid.UUID, date time.Time) error {
	f.deleteAttendanceForDateCalls++
	f.deleteAttendanceForDateArg.teacherUserID = teacherUserID
	f.deleteAttendanceForDateArg.studentID = studentID
	f.deleteAttendanceForDateArg.date = date
	return f.deleteAttendanceForDateErr
}

func (f *fakeTeacherScopeServiceTeacherScopeRepo) ListAttendanceByStudent(_ context.Context, _, _ uuid.UUID, _, _ time.Time) ([]model.AttendanceRecord, error) {
	return f.listAttendanceByStudentRes, f.listAttendanceByStudentErr
}

func (f *fakeTeacherScopeServiceTeacherScopeRepo) ListAttendanceChangeLogsByStudent(_ context.Context, _, _ uuid.UUID, _, _ time.Time) ([]model.AttendanceChangeLog, error) {
	return f.listAttendanceChangeLogsByStudentRes, f.listAttendanceChangeLogsByStudentErr
}

func (f *fakeTeacherScopeServiceTeacherScopeRepo) ListAttendanceChangeLogsByClass(_ context.Context, _ uuid.UUID, _ uuid.UUID, _ *uuid.UUID, status *string, _, _ time.Time, limit, offset int) ([]model.AttendanceChangeLog, int, error) {
	f.listAttendanceChangeLogsByClassCalls++
	f.listAttendanceChangeLogsByClassArg.limit = limit
	f.listAttendanceChangeLogsByClassArg.offset = offset
	f.listAttendanceChangeLogsByClassArg.status = status
	return f.listAttendanceChangeLogsByClassRes, f.listAttendanceChangeLogsByClassTotal, f.listAttendanceChangeLogsByClassErr
}

func (f *fakeTeacherScopeServiceTeacherScopeRepo) CreateClassPost(_ context.Context, _, _ uuid.UUID, _, _ string) (uuid.UUID, error) {
	return f.createClassPostRes, f.createClassPostErr
}

func (f *fakeTeacherScopeServiceTeacherScopeRepo) CreateStudentPost(_ context.Context, _, _ uuid.UUID, _, _ string) (uuid.UUID, error) {
	return f.createStudentPostRes, f.createStudentPostErr
}

func (f *fakeTeacherScopeServiceTeacherScopeRepo) ListClassPosts(_ context.Context, _, _ uuid.UUID, limit, offset int) ([]model.Post, int, error) {
	f.listClassPostsCalls++
	f.listClassPostsArg.limit = limit
	f.listClassPostsArg.offset = offset
	return f.listClassPostsRes, f.listClassPostsTotal, f.listClassPostsErr
}

func (f *fakeTeacherScopeServiceTeacherScopeRepo) ListStudentPosts(_ context.Context, _, _ uuid.UUID, limit, offset int) ([]model.Post, int, error) {
	f.listStudentPostsCalls++
	f.listStudentPostsArg.limit = limit
	f.listStudentPostsArg.offset = offset
	return f.listStudentPostsRes, f.listStudentPostsTotal, f.listStudentPostsErr
}

func (f *fakeTeacherScopeServiceTeacherScopeRepo) UpdatePost(_ context.Context, _, _ uuid.UUID, _ string) error {
	return f.updatePostErr
}

func (f *fakeTeacherScopeServiceTeacherScopeRepo) DeletePost(_ context.Context, _, _ uuid.UUID) error {
	return f.deletePostErr
}

func (f *fakeTeacherScopeServiceTeacherScopeRepo) CountMyStudents(_ context.Context, _ uuid.UUID) (int, error) {
	return f.countMyStudentsRes, f.countMyStudentsErr
}

func (f *fakeTeacherScopeServiceTeacherScopeRepo) CountMyPosts(_ context.Context, _ uuid.UUID) (int, error) {
	return f.countMyPostsRes, f.countMyPostsErr
}

type fakeTeacherScopeServiceHealthLogRepo struct {
	createRes uuid.UUID
	createErr error
	listRes   []model.HealthLog
	listErr   error
}

func (f *fakeTeacherScopeServiceHealthLogRepo) CreateByStudentAndTeacher(_ context.Context, _, _ uuid.UUID, _ *time.Time, _ *float64, _ string, _ *string, _ string) (uuid.UUID, error) {
	return f.createRes, f.createErr
}

func (f *fakeTeacherScopeServiceHealthLogRepo) ListByStudentAndTeacher(_ context.Context, _, _ uuid.UUID, _, _ time.Time) ([]model.HealthLog, error) {
	return f.listRes, f.listErr
}

type fakeTeacherScopeServiceTeacherRepo struct {
	getByUserRes *model.Teacher
	getByUserErr error
	updateErr    error
}

func (f *fakeTeacherScopeServiceTeacherRepo) GetByUserID(_ context.Context, _ uuid.UUID) (*model.Teacher, error) {
	return f.getByUserRes, f.getByUserErr
}

func (f *fakeTeacherScopeServiceTeacherRepo) UpdatePhone(_ context.Context, _ uuid.UUID, _ string) error {
	return f.updateErr
}

type fakeTeacherScopeServicePostInteractionRepo struct {
	canAccessRes  bool
	canAccessErr  error
	toggleRes     bool
	toggleCount   int
	toggleErr     error
	addComment    model.PostComment
	addCommentErr error
	countComments int
	countErr      error
	listComments  []model.PostComment
	listTotal     int
	listErr       error
	addShare      int
	addShareErr   error
}

func (f *fakeTeacherScopeServicePostInteractionRepo) TeacherCanAccessPost(_ context.Context, _, _ uuid.UUID) (bool, error) {
	return f.canAccessRes, f.canAccessErr
}

func (f *fakeTeacherScopeServicePostInteractionRepo) ToggleLike(_ context.Context, _, _ uuid.UUID) (bool, int, error) {
	return f.toggleRes, f.toggleCount, f.toggleErr
}

func (f *fakeTeacherScopeServicePostInteractionRepo) AddComment(_ context.Context, _, _ uuid.UUID, _ string) (model.PostComment, error) {
	return f.addComment, f.addCommentErr
}

func (f *fakeTeacherScopeServicePostInteractionRepo) CountComments(_ context.Context, _ uuid.UUID) (int, error) {
	return f.countComments, f.countErr
}

func (f *fakeTeacherScopeServicePostInteractionRepo) ListComments(_ context.Context, _ uuid.UUID, _, _ int) ([]model.PostComment, int, error) {
	return f.listComments, f.listTotal, f.listErr
}

func (f *fakeTeacherScopeServicePostInteractionRepo) AddShare(_ context.Context, _, _ uuid.UUID) (int, error) {
	return f.addShare, f.addShareErr
}

func newTeacherScopeServiceForTest() (*TeacherScopeService, *fakeTeacherScopeServiceTeacherScopeRepo, *fakeTeacherScopeServiceHealthLogRepo, *fakeTeacherScopeServiceTeacherRepo, *fakeTeacherScopeServicePostInteractionRepo) {
	teacherScopeRepo := &fakeTeacherScopeServiceTeacherScopeRepo{}
	healthLogRepo := &fakeTeacherScopeServiceHealthLogRepo{}
	teacherRepo := &fakeTeacherScopeServiceTeacherRepo{}
	postRepo := &fakeTeacherScopeServicePostInteractionRepo{}
	svc := &TeacherScopeService{teacherScopeRepo: teacherScopeRepo, healthLogRepo: healthLogRepo, teacherRepo: teacherRepo, postInteractRepo: postRepo}
	return svc, teacherScopeRepo, healthLogRepo, teacherRepo, postRepo
}

func TestTeacherScopeServiceListMyClasses(t *testing.T) {
	svc, repoFake, _, _, _ := newTeacherScopeServiceForTest()

	_, err := svc.ListMyClasses(context.Background(), uuid.Nil)
	if !errors.Is(err, ErrInvalidUserID) {
		t.Fatalf("expected ErrInvalidUserID, got %v", err)
	}
	if repoFake.listMyClassesCalls != 0 {
		t.Fatalf("repo should not be called for invalid user")
	}

	repoFake.listMyClassesErr = errors.New("db fail")
	_, err = svc.ListMyClasses(context.Background(), uuid.New())
	if err == nil || err.Error() != "failed to list classes: db fail" {
		t.Fatalf("unexpected error: %v", err)
	}

	repoFake.listMyClassesErr = nil
	repoFake.listMyClassesRes = []model.Class{{ClassID: uuid.New()}, {ClassID: uuid.New()}}
	got, err := svc.ListMyClasses(context.Background(), uuid.New())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 classes, got %d", len(got))
	}
}

func TestTeacherScopeServiceUpsertAttendance(t *testing.T) {
	svc, repoFake, _, _, _ := newTeacherScopeServiceForTest()
	teacherUserID := uuid.New()
	studentID := uuid.New()

	if err := svc.UpsertAttendance(context.Background(), uuid.Nil, studentID, "2026-01-01", "present", nil, nil, "ok"); !errors.Is(err, ErrInvalidUserID) {
		t.Fatalf("expected ErrInvalidUserID, got %v", err)
	}
	if err := svc.UpsertAttendance(context.Background(), teacherUserID, uuid.Nil, "2026-01-01", "present", nil, nil, "ok"); !errors.Is(err, ErrInvalidUserID) {
		t.Fatalf("expected ErrInvalidUserID for nil student, got %v", err)
	}
	if err := svc.UpsertAttendance(context.Background(), teacherUserID, studentID, "bad", "present", nil, nil, "ok"); !errors.Is(err, ErrInvalidDate) {
		t.Fatalf("expected ErrInvalidDate, got %v", err)
	}
	if err := svc.UpsertAttendance(context.Background(), teacherUserID, studentID, "2026-01-01", "bad", nil, nil, "ok"); !errors.Is(err, ErrInvalidStatus) {
		t.Fatalf("expected ErrInvalidStatus, got %v", err)
	}

	repoFake.upsertAttendanceErr = repo.ErrNoRowsUpdated
	if err := svc.UpsertAttendance(context.Background(), teacherUserID, studentID, "2026-01-02", "present", nil, nil, "note"); !errors.Is(err, ErrForbidden) {
		t.Fatalf("expected ErrForbidden, got %v", err)
	}

	repoFake.upsertAttendanceErr = errors.New("boom")
	err := svc.UpsertAttendance(context.Background(), teacherUserID, studentID, "2026-01-02", "present", nil, nil, "note")
	if err == nil || err.Error() != "failed to mark attendance: boom" {
		t.Fatalf("unexpected wrapped error: %v", err)
	}

	repoFake.upsertAttendanceErr = nil
	if err := svc.UpsertAttendance(context.Background(), teacherUserID, studentID, "2026-01-03", "late", nil, nil, "ok"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if repoFake.upsertAttendanceCalls == 0 {
		t.Fatal("expected repo upsert to be called")
	}
	if repoFake.upsertAttendanceArg.status != "late" || repoFake.upsertAttendanceArg.note != "ok" {
		t.Fatalf("unexpected forwarded values: %#v", repoFake.upsertAttendanceArg)
	}
	if repoFake.upsertAttendanceArg.date.Format("2006-01-02") != "2026-01-03" {
		t.Fatalf("date not parsed/forwarded correctly: %s", repoFake.upsertAttendanceArg.date.Format("2006-01-02"))
	}
}

func TestTeacherScopeServiceCancelAttendanceForDate(t *testing.T) {
	svc, repoFake, _, _, _ := newTeacherScopeServiceForTest()
	teacherUserID := uuid.New()
	studentID := uuid.New()

	if err := svc.CancelAttendanceForDate(context.Background(), uuid.Nil, studentID, "2026-01-01"); !errors.Is(err, ErrInvalidUserID) {
		t.Fatalf("expected ErrInvalidUserID, got %v", err)
	}
	if err := svc.CancelAttendanceForDate(context.Background(), teacherUserID, uuid.Nil, "2026-01-01"); !errors.Is(err, ErrInvalidUserID) {
		t.Fatalf("expected ErrInvalidUserID, got %v", err)
	}
	if err := svc.CancelAttendanceForDate(context.Background(), teacherUserID, studentID, "bad"); !errors.Is(err, ErrInvalidDate) {
		t.Fatalf("expected ErrInvalidDate, got %v", err)
	}

	repoFake.deleteAttendanceForDateErr = repo.ErrNoRowsUpdated
	if err := svc.CancelAttendanceForDate(context.Background(), teacherUserID, studentID, "2026-01-05"); !errors.Is(err, ErrForbidden) {
		t.Fatalf("expected ErrForbidden, got %v", err)
	}

	repoFake.deleteAttendanceForDateErr = errors.New("x")
	err := svc.CancelAttendanceForDate(context.Background(), teacherUserID, studentID, "2026-01-05")
	if err == nil || err.Error() != "failed to cancel attendance: x" {
		t.Fatalf("unexpected error: %v", err)
	}

	repoFake.deleteAttendanceForDateErr = nil
	if err := svc.CancelAttendanceForDate(context.Background(), teacherUserID, studentID, "2026-01-05"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if repoFake.deleteAttendanceForDateArg.date.Format("2006-01-02") != "2026-01-05" {
		t.Fatalf("unexpected forwarded date: %s", repoFake.deleteAttendanceForDateArg.date.Format("2006-01-02"))
	}
}

func TestTeacherScopeServiceListAttendanceChangeLogsByClass(t *testing.T) {
	svc, repoFake, _, _, _ := newTeacherScopeServiceForTest()
	teacherUserID := uuid.New()
	classID := uuid.New()

	_, _, err := svc.ListAttendanceChangeLogsByClass(context.Background(), uuid.Nil, classID, nil, nil, time.Now(), time.Now(), 10, 0)
	if !errors.Is(err, ErrInvalidUserID) {
		t.Fatalf("expected ErrInvalidUserID, got %v", err)
	}

	_, _, err = svc.ListAttendanceChangeLogsByClass(context.Background(), teacherUserID, uuid.Nil, nil, nil, time.Now(), time.Now(), 10, 0)
	if !errors.Is(err, ErrInvalidClassID) {
		t.Fatalf("expected ErrInvalidClassID, got %v", err)
	}

	badStatus := "bad"
	_, _, err = svc.ListAttendanceChangeLogsByClass(context.Background(), teacherUserID, classID, nil, &badStatus, time.Now(), time.Now(), 10, 0)
	if !errors.Is(err, ErrInvalidStatus) {
		t.Fatalf("expected ErrInvalidStatus, got %v", err)
	}

	repoFake.listAttendanceChangeLogsByClassErr = errors.New("repo")
	_, _, err = svc.ListAttendanceChangeLogsByClass(context.Background(), teacherUserID, classID, nil, nil, time.Now(), time.Now(), 10, 2)
	if err == nil || err.Error() != "failed to list class attendance change logs: repo" {
		t.Fatalf("unexpected error: %v", err)
	}

	repoFake.listAttendanceChangeLogsByClassErr = nil
	repoFake.listAttendanceChangeLogsByClassRes = []model.AttendanceChangeLog{{ChangeID: uuid.New()}}
	repoFake.listAttendanceChangeLogsByClassTotal = 1
	logs, total, err := svc.ListAttendanceChangeLogsByClass(context.Background(), teacherUserID, classID, nil, nil, time.Now(), time.Now(), 0, 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(logs) != 1 || total != 1 {
		t.Fatalf("unexpected logs/total: len=%d total=%d", len(logs), total)
	}
	if repoFake.listAttendanceChangeLogsByClassArg.limit != 20 {
		t.Fatalf("expected default limit 20, got %d", repoFake.listAttendanceChangeLogsByClassArg.limit)
	}

	_, _, err = svc.ListAttendanceChangeLogsByClass(context.Background(), teacherUserID, classID, nil, nil, time.Now(), time.Now(), 999, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if repoFake.listAttendanceChangeLogsByClassArg.limit != 100 {
		t.Fatalf("expected capped limit 100, got %d", repoFake.listAttendanceChangeLogsByClassArg.limit)
	}
}

func TestTeacherScopeServiceCreateHealthLogAndListHealthLogs(t *testing.T) {
	svc, _, healthRepo, _, _ := newTeacherScopeServiceForTest()
	teacherUserID := uuid.New()
	studentID := uuid.New()
	invalidSeverity := "invalid"

	_, err := svc.CreateHealthLog(context.Background(), uuid.Nil, studentID, nil, nil, "", "", nil)
	if !errors.Is(err, ErrInvalidUserID) {
		t.Fatalf("expected ErrInvalidUserID, got %v", err)
	}
	_, err = svc.CreateHealthLog(context.Background(), teacherUserID, uuid.Nil, nil, nil, "", "", nil)
	if !errors.Is(err, ErrInvalidUserID) {
		t.Fatalf("expected ErrInvalidUserID, got %v", err)
	}
	_, err = svc.CreateHealthLog(context.Background(), teacherUserID, studentID, nil, nil, "", "", &invalidSeverity)
	if !errors.Is(err, ErrInvalidValue) {
		t.Fatalf("expected ErrInvalidValue for severity, got %v", err)
	}

	healthRepo.createErr = repo.ErrNoRowsUpdated
	_, err = svc.CreateHealthLog(context.Background(), teacherUserID, studentID, nil, nil, "", "", nil)
	if !errors.Is(err, ErrForbidden) {
		t.Fatalf("expected ErrForbidden, got %v", err)
	}

	healthRepo.createErr = nil
	healthRepo.createRes = uuid.New()
	id, err := svc.CreateHealthLog(context.Background(), teacherUserID, studentID, nil, nil, "", "", nil)
	if err != nil || id == uuid.Nil {
		t.Fatalf("unexpected result id=%v err=%v", id, err)
	}

	from := time.Now().Add(-24 * time.Hour)
	to := time.Now()
	healthRepo.listErr = pgx.ErrNoRows
	logs, err := svc.ListHealthLogs(context.Background(), teacherUserID, studentID, from, to)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(logs) != 0 {
		t.Fatalf("expected empty logs, got %d", len(logs))
	}

	healthRepo.listErr = nil
	healthRepo.listRes = []model.HealthLog{{HealthLogID: uuid.New()}}
	logs, err = svc.ListHealthLogs(context.Background(), teacherUserID, studentID, from, to)
	if err != nil || len(logs) != 1 {
		t.Fatalf("unexpected list result len=%d err=%v", len(logs), err)
	}
}

func TestTeacherScopeServiceUpdateMyProfile(t *testing.T) {
	svc, _, _, teacherRepo, _ := newTeacherScopeServiceForTest()
	teacherUserID := uuid.New()

	if err := svc.UpdateMyProfile(context.Background(), uuid.Nil, "090"); !errors.Is(err, ErrInvalidUserID) {
		t.Fatalf("expected ErrInvalidUserID, got %v", err)
	}

	teacherRepo.getByUserErr = pgx.ErrNoRows
	if err := svc.UpdateMyProfile(context.Background(), teacherUserID, "090"); !errors.Is(err, ErrTeacherNotFound) {
		t.Fatalf("expected ErrTeacherNotFound, got %v", err)
	}

	teacherRepo.getByUserErr = nil
	teacherRepo.getByUserRes = &model.Teacher{TeacherID: uuid.New()}
	teacherRepo.updateErr = errors.New("update fail")
	err := svc.UpdateMyProfile(context.Background(), teacherUserID, "090")
	if err == nil || err.Error() != "failed to update profile: update fail" {
		t.Fatalf("unexpected error: %v", err)
	}

	teacherRepo.updateErr = nil
	if err := svc.UpdateMyProfile(context.Background(), teacherUserID, "090"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTeacherScopeServicePostCRUDAndLists(t *testing.T) {
	svc, repoFake, _, _, _ := newTeacherScopeServiceForTest()
	teacherUserID := uuid.New()
	classID := uuid.New()
	studentID := uuid.New()
	postID := uuid.New()

	_, err := svc.CreateClassPost(context.Background(), teacherUserID, classID, "bad", "hello")
	if !errors.Is(err, ErrInvalidValue) {
		t.Fatalf("expected ErrInvalidValue, got %v", err)
	}
	_, err = svc.CreateStudentPost(context.Background(), teacherUserID, studentID, "announcement", "")
	if !errors.Is(err, ErrInvalidValue) {
		t.Fatalf("expected ErrInvalidValue for empty content, got %v", err)
	}

	repoFake.createClassPostErr = repo.ErrNoRowsUpdated
	_, err = svc.CreateClassPost(context.Background(), teacherUserID, classID, "announcement", "x")
	if !errors.Is(err, ErrForbidden) {
		t.Fatalf("expected ErrForbidden, got %v", err)
	}

	repoFake.createClassPostErr = nil
	repoFake.createClassPostRes = uuid.New()
	id, err := svc.CreateClassPost(context.Background(), teacherUserID, classID, "announcement", "x")
	if err != nil || id == uuid.Nil {
		t.Fatalf("unexpected create class post result id=%v err=%v", id, err)
	}

	repoFake.listClassPostsErr = repo.ErrNoRowsUpdated
	posts, total, err := svc.ListClassPosts(context.Background(), teacherUserID, classID, 0, 5)
	if err != nil {
		t.Fatalf("unexpected error on no rows updated: %v", err)
	}
	if len(posts) != 0 || total != 0 {
		t.Fatalf("expected empty posts/0 total, got len=%d total=%d", len(posts), total)
	}
	if repoFake.listClassPostsArg.limit != 20 {
		t.Fatalf("expected default limit 20, got %d", repoFake.listClassPostsArg.limit)
	}

	repoFake.listClassPostsErr = nil
	repoFake.listClassPostsRes = []model.Post{{PostID: uuid.New()}}
	repoFake.listClassPostsTotal = 1
	posts, total, err = svc.ListClassPosts(context.Background(), teacherUserID, classID, 999, 0)
	if err != nil || len(posts) != 1 || total != 1 {
		t.Fatalf("unexpected list class posts result len=%d total=%d err=%v", len(posts), total, err)
	}
	if repoFake.listClassPostsArg.limit != 100 {
		t.Fatalf("expected capped limit 100, got %d", repoFake.listClassPostsArg.limit)
	}

	repoFake.updatePostErr = repo.ErrNoRowsUpdated
	if err := svc.UpdatePost(context.Background(), teacherUserID, postID, "content"); !errors.Is(err, ErrForbidden) {
		t.Fatalf("expected ErrForbidden, got %v", err)
	}
	if err := svc.UpdatePost(context.Background(), teacherUserID, postID, "  \t\n  "); !errors.Is(err, ErrInvalidValue) {
		t.Fatalf("expected ErrInvalidValue for blank content, got %v", err)
	}

	repoFake.deletePostErr = repo.ErrNoRowsUpdated
	if err := svc.DeletePost(context.Background(), teacherUserID, postID); !errors.Is(err, ErrForbidden) {
		t.Fatalf("expected ErrForbidden, got %v", err)
	}

	repoFake.listStudentPostsErr = nil
	repoFake.listStudentPostsRes = []model.Post{{PostID: uuid.New()}}
	repoFake.listStudentPostsTotal = 1
	studentPosts, studentTotal, err := svc.ListStudentPosts(context.Background(), teacherUserID, studentID, 0, 0)
	if err != nil || len(studentPosts) != 1 || studentTotal != 1 {
		t.Fatalf("unexpected list student posts result len=%d total=%d err=%v", len(studentPosts), studentTotal, err)
	}
	if repoFake.listStudentPostsArg.limit != 20 {
		t.Fatalf("expected default limit 20, got %d", repoFake.listStudentPostsArg.limit)
	}
}

func TestTeacherScopeServicePostInteractions(t *testing.T) {
	svc, _, _, _, postRepo := newTeacherScopeServiceForTest()
	teacherUserID := uuid.New()
	postID := uuid.New()

	postRepo.canAccessErr = errors.New("access check failed")
	_, _, err := svc.TogglePostLike(context.Background(), teacherUserID, postID)
	if err == nil || err.Error() != "failed to verify post access: access check failed" {
		t.Fatalf("unexpected error: %v", err)
	}

	postRepo.canAccessErr = nil
	postRepo.canAccessRes = false
	_, _, err = svc.TogglePostLike(context.Background(), teacherUserID, postID)
	if !errors.Is(err, ErrForbidden) {
		t.Fatalf("expected ErrForbidden, got %v", err)
	}

	postRepo.canAccessRes = true
	postRepo.toggleErr = nil
	postRepo.toggleRes = true
	postRepo.toggleCount = 4
	liked, likeCount, err := svc.TogglePostLike(context.Background(), teacherUserID, postID)
	if err != nil || !liked || likeCount != 4 {
		t.Fatalf("unexpected toggle result liked=%v count=%d err=%v", liked, likeCount, err)
	}

	postRepo.addComment = model.PostComment{CommentID: uuid.New()}
	postRepo.countComments = 3
	comment, commentCount, err := svc.AddPostComment(context.Background(), teacherUserID, postID, "  hello  ")
	if err != nil || comment.CommentID == uuid.Nil || commentCount != 3 {
		t.Fatalf("unexpected add comment result comment=%#v count=%d err=%v", comment, commentCount, err)
	}

	_, _, err = svc.AddPostComment(context.Background(), teacherUserID, postID, "   ")
	if !errors.Is(err, ErrInvalidValue) {
		t.Fatalf("expected ErrInvalidValue, got %v", err)
	}

	postRepo.listComments = []model.PostComment{{CommentID: uuid.New()}}
	postRepo.listTotal = 1
	comments, total, err := svc.ListPostComments(context.Background(), teacherUserID, postID, 0, 0)
	if err != nil || len(comments) != 1 || total != 1 {
		t.Fatalf("unexpected list comments result len=%d total=%d err=%v", len(comments), total, err)
	}

	postRepo.addShare = 6
	shareCount, err := svc.SharePost(context.Background(), teacherUserID, postID)
	if err != nil || shareCount != 6 {
		t.Fatalf("unexpected share result count=%d err=%v", shareCount, err)
	}
}
