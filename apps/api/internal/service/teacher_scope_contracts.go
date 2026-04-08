package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/model"
)

// teacherScopeStatsRepo groups read-model methods reused by teacher analytics flows.
type teacherScopeStatsRepo interface {
	ListMyClasses(ctx context.Context, teacherUserID uuid.UUID) ([]model.Class, error)
	CountMyStudents(ctx context.Context, teacherUserID uuid.UUID) (int, error)
	CountMyPosts(ctx context.Context, teacherUserID uuid.UUID) (int, error)
}

// teacherScopeAttendanceRepo groups attendance and attendance-change-log methods.
type teacherScopeAttendanceRepo interface {
	ListMyStudentsInClass(ctx context.Context, teacherUserID, classID uuid.UUID) ([]model.Student, error)
	UpsertAttendance(ctx context.Context, teacherUserID, studentID uuid.UUID, date time.Time, status string, checkInAt, checkOutAt *time.Time, note string) error
	DeleteAttendanceForDate(ctx context.Context, teacherUserID, studentID uuid.UUID, date time.Time) error
	ListAttendanceByStudent(ctx context.Context, teacherUserID, studentID uuid.UUID, from, to time.Time) ([]model.AttendanceRecord, error)
	ListAttendanceChangeLogsByStudent(ctx context.Context, teacherUserID, studentID uuid.UUID, from, to time.Time) ([]model.AttendanceChangeLog, error)
	ListAttendanceChangeLogsByClass(ctx context.Context, teacherUserID, classID uuid.UUID, studentID *uuid.UUID, status *string, from, to time.Time, limit, offset int) ([]model.AttendanceChangeLog, int, error)
}

// teacherScopePostRepo groups post management methods under teacher scope.
type teacherScopePostRepo interface {
	CreateClassPost(ctx context.Context, teacherUserID, classID uuid.UUID, postType, content string) (uuid.UUID, error)
	CreateStudentPost(ctx context.Context, teacherUserID, studentID uuid.UUID, postType, content string) (uuid.UUID, error)
	ListClassPosts(ctx context.Context, teacherUserID, classID uuid.UUID, limit, offset int) ([]model.Post, int, error)
	ListStudentPosts(ctx context.Context, teacherUserID, studentID uuid.UUID, limit, offset int) ([]model.Post, int, error)
	UpdatePost(ctx context.Context, authorUserID, postID uuid.UUID, newContent string) error
	DeletePost(ctx context.Context, authorUserID, postID uuid.UUID) error
}

// teacherScopeHealthLogRepo groups health log methods under teacher scope.
type teacherScopeHealthLogRepo interface {
	CreateByStudentAndTeacher(ctx context.Context, teacherUserID, studentID uuid.UUID, recordedAt *time.Time, temperature *float64, symptoms string, severity *string, note string) (uuid.UUID, error)
	ListByStudentAndTeacher(ctx context.Context, teacherUserID, studentID uuid.UUID, from, to time.Time) ([]model.HealthLog, error)
}

// teacherScopeProfileRepo groups teacher profile methods needed by teacher-scope service.
type teacherScopeProfileRepo interface {
	GetByUserID(ctx context.Context, userID uuid.UUID) (*model.Teacher, error)
	UpdatePhone(ctx context.Context, teacherID uuid.UUID, phone string) error
}

// teacherScopeInteractionRepo groups post interaction methods used by teacher scope.
type teacherScopeInteractionRepo interface {
	TeacherCanAccessPost(ctx context.Context, teacherUserID, postID uuid.UUID) (bool, error)
	ToggleLike(ctx context.Context, userID, postID uuid.UUID) (bool, int, error)
	AddComment(ctx context.Context, userID, postID uuid.UUID, content string) (model.PostComment, error)
	CountComments(ctx context.Context, postID uuid.UUID) (int, error)
	ListComments(ctx context.Context, postID uuid.UUID, limit, offset int) ([]model.PostComment, int, error)
	AddShare(ctx context.Context, userID, postID uuid.UUID) (int, error)
}
