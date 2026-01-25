//go:build !test

package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/hoanghonghuy/iris-app/apps/api/internal/model"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/repo"
)

// UpdateMyProfileRequest represents the request to update teacher's own profile (teacher only - can only update phone)
type UpdateMyProfileRequest struct {
	Phone string `json:"phone"`
}

type TeacherScopeService struct {
	teacherScopeRepo *repo.TeacherScopeRepo
	teacherRepo      *repo.TeacherRepo
}

func NewTeacherScopeService(teacherScopeRepo *repo.TeacherScopeRepo, teacherRepo *repo.TeacherRepo) *TeacherScopeService {
	return &TeacherScopeService{
		teacherScopeRepo: teacherScopeRepo,
		teacherRepo:      teacherRepo,
	}
}

// ListMyClasses returns list of classes that the teacher (by user_id) is assigned to teach
func (s *TeacherScopeService) ListMyClasses(ctx context.Context, teacherUserID uuid.UUID) ([]model.Class, error) {
	// Validate teacherUserID
	if teacherUserID == uuid.Nil {
		return nil, ErrInvalidUserID
	}

	classes, err := s.teacherScopeRepo.ListMyClass(ctx, teacherUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to list classes: %w", err)
	}

	return classes, nil
}

// ListMyStudentsInClass returns list of students in a class if the teacher is assigned to that class
func (s *TeacherScopeService) ListMyStudentsInClass(ctx context.Context, teacherUserID, classID uuid.UUID) ([]model.Student, error) {
	// Validate teacherUserID
	if teacherUserID == uuid.Nil {
		return nil, ErrInvalidUserID
	}

	// Validate classID
	if classID == uuid.Nil {
		return nil, ErrInvalidClassID
	}

	students, err := s.teacherScopeRepo.ListMyStudentsInClass(ctx, teacherUserID, classID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrForbidden
		}
		return nil, fmt.Errorf("failed to list students: %w", err)
	}

	return students, nil
}

// UpsertAttendance marks or updates attendance for a student
// Teacher can only mark attendance for students in their assigned classes
func (s *TeacherScopeService) UpsertAttendance(ctx context.Context, teacherUserID, studentID uuid.UUID, date string, status, note string) error {
	// Validate teacherUserID
	if teacherUserID == uuid.Nil {
		return ErrInvalidUserID
	}

	// Validate studentID
	if studentID == uuid.Nil {
		return ErrInvalidUserID
	}

	// Validate date format
	parsedDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return ErrInvalidDate
	}

	// Validate status
	validStatuses := map[string]bool{
		"present": true,
		"absent":  true,
		"late":    true,
		"excused": true,
	}
	if !validStatuses[status] {
		return ErrInvalidStatus
	}

	// Call repo to upsert attendance
	err = s.teacherScopeRepo.UpsertAttendance(ctx, teacherUserID, studentID, parsedDate, status, note)
	if err != nil {
		if err == repo.ErrForbidden {
			return ErrForbidden
		}
		return fmt.Errorf("failed to mark attendance: %w", err)
	}

	return nil
}

// UpdateMyProfile updates teacher's own profile (teacher only - can only update phone)
func (s *TeacherScopeService) UpdateMyProfile(ctx context.Context, teacherUserID uuid.UUID, req UpdateMyProfileRequest) error {
	// Validate teacherUserID
	if teacherUserID == uuid.Nil {
		return ErrInvalidUserID
	}

	// Get teacherID from teacherUserID
	teacher, err := s.teacherRepo.GetByUserID(ctx, teacherUserID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return ErrTeacherNotFound
		}
		return fmt.Errorf("failed to get teacher: %w", err)
	}

	// Teacher can only update phone (cannot update school_id, teacher_id, user_id)
	err = s.teacherRepo.UpdatePhone(ctx, teacher.TeacherID, req.Phone)
	if err != nil {
		return fmt.Errorf("failed to update profile: %w", err)
	}

	return nil
}
