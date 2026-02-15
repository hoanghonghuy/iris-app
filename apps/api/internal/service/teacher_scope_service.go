package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/hoanghonghuy/iris-app/apps/api/internal/model"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/repo"
)

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

	classes, err := s.teacherScopeRepo.ListMyClasses(ctx, teacherUserID)
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
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrForbidden
		}
		return nil, fmt.Errorf("failed to list students: %w", err)
	}

	return students, nil
}

// UpsertAttendance marks or updates attendance for a student
// Teacher can only mark attendance for students in their assigned classes
func (s *TeacherScopeService) UpsertAttendance(ctx context.Context, teacherUserID, studentID uuid.UUID,
	date string, status string, checkInAt, checkOutAt *time.Time, note string) error {
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
	err = s.teacherScopeRepo.UpsertAttendance(ctx, teacherUserID, studentID, parsedDate, status, checkInAt, checkOutAt, note)
	if err != nil {
		if errors.Is(err, repo.ErrForbidden) {
			return ErrForbidden
		}
		return fmt.Errorf("failed to mark attendance: %w", err)
	}

	return nil
}

// ListAttendanceByStudent liệt kê lịch sử điểm danh của một học sinh.
// Giáo viên chỉ có thể xem điểm danh của học sinh trong các lớp được phân công.
func (s *TeacherScopeService) ListAttendanceByStudent(ctx context.Context, teacherUserID, studentID uuid.UUID,
	from, to time.Time) ([]model.AttendanceRecord, error) {

	if teacherUserID == uuid.Nil {
		return nil, ErrInvalidUserID
	}

	if studentID == uuid.Nil {
		return nil, ErrInvalidUserID
	}

	records, err := s.teacherScopeRepo.ListAttendanceByStudent(ctx, teacherUserID, studentID, from, to)
	if err != nil {
		return nil, fmt.Errorf("failed to list attendance: %w", err)
	}

	return records, nil
}

// CreateHealthLog tạo nhật ký sức khỏe mới cho học sinh
func (s *TeacherScopeService) CreateHealthLog(ctx context.Context, teacherUserID, studentID uuid.UUID,
	recordedAt *time.Time, temperature *float64, symptoms, note string, severity *string) (uuid.UUID, error) {
	// Validate teacherUserID
	if teacherUserID == uuid.Nil {
		return uuid.Nil, ErrInvalidUserID
	}

	// Validate studentID
	if studentID == uuid.Nil {
		return uuid.Nil, ErrInvalidUserID
	}

	// Validate severity
	if severity != nil {
		sev := *severity
		if sev != "normal" && sev != "watch" && sev != "urgent" {
			return uuid.Nil, fmt.Errorf("%w: severity must be normal|watch|urgent", ErrInvalidValue)
		}
	}

	// xác minh giáo viên có quyền truy cập
	id, err := s.teacherScopeRepo.CreateHealthLog(ctx, teacherUserID, studentID, recordedAt, temperature, symptoms, severity, note)
	if err != nil {
		if errors.Is(err, repo.ErrForbidden) {
			return uuid.Nil, ErrForbidden
		}
		return uuid.Nil, fmt.Errorf("failed to create health log: %w", err)
	}

	return id, nil
}

// ListHealthLogs liệt kê nhật ký sức khỏe của một học sinh.
func (s *TeacherScopeService) ListHealthLogs(ctx context.Context, teacherUserID, studentID uuid.UUID,
	from, to time.Time) ([]model.HealthLog, error) {
	// Validate teacherUserID
	if teacherUserID == uuid.Nil {
		return nil, ErrInvalidUserID
	}

	// Validate studentID
	if studentID == uuid.Nil {
		return nil, ErrInvalidUserID
	}

	healthLogs, err := s.teacherScopeRepo.ListHealthLogsByStudent(ctx, teacherUserID, studentID, from, to)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) || errors.Is(err, repo.ErrForbidden) {
			return []model.HealthLog{}, nil
		}
		return nil, fmt.Errorf("failed to list health logs: %w", err)
	}

	return healthLogs, nil
}

// UpdateMyProfile updates teacher's own profile (teacher only - can only update phone)
func (s *TeacherScopeService) UpdateMyProfile(ctx context.Context, teacherUserID uuid.UUID, phone string) error {
	// Validate teacherUserID
	if teacherUserID == uuid.Nil {
		return ErrInvalidUserID
	}

	// Get teacherID from teacherUserID
	teacher, err := s.teacherRepo.GetByUserID(ctx, teacherUserID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrTeacherNotFound
		}
		return fmt.Errorf("failed to get teacher: %w", err)
	}

	// Teacher can only update phone (cannot update school_id, teacher_id, user_id)
	err = s.teacherRepo.UpdatePhone(ctx, teacher.TeacherID, phone)
	if err != nil {
		return fmt.Errorf("failed to update profile: %w", err)
	}

	return nil
}

// CreateClassPost tạo bài đăng cho một lớp học
func (s *TeacherScopeService) CreateClassPost(ctx context.Context, teacherUserID, classID uuid.UUID,
	postType, content string) (uuid.UUID, error) {
	// Validate teacherUserID
	if teacherUserID == uuid.Nil {
		return uuid.Nil, ErrInvalidUserID
	}

	// Validate classID
	if classID == uuid.Nil {
		return uuid.Nil, ErrInvalidClassID
	}

	// Validate postType
	if !isValidPostType(postType) {
		return uuid.Nil, fmt.Errorf("%w: type must be announcement|activity|daily_note|health_note", ErrInvalidValue)
	}

	// Validate content
	if content == "" {
		return uuid.Nil, fmt.Errorf("%w: content cannot be empty", ErrInvalidValue)
	}

	id, err := s.teacherScopeRepo.CreateClassPost(ctx, teacherUserID, classID, postType, content)
	if err != nil {
		if errors.Is(err, repo.ErrForbidden) {
			return uuid.Nil, ErrForbidden
		}
		return uuid.Nil, fmt.Errorf("failed to create class post: %w", err)
	}

	return id, nil
}

// CreateStudentPost tạo bài đăng cho một học sinh
func (s *TeacherScopeService) CreateStudentPost(ctx context.Context, teacherUserID, studentID uuid.UUID,
	postType, content string) (uuid.UUID, error) {
	// Validate teacherUserID
	if teacherUserID == uuid.Nil {
		return uuid.Nil, ErrInvalidUserID
	}

	// Validate studentID
	if studentID == uuid.Nil {
		return uuid.Nil, ErrInvalidUserID
	}

	// Validate postType
	if !isValidPostType(postType) {
		return uuid.Nil, fmt.Errorf("%w: type must be announcement|activity|daily_note|health_note", ErrInvalidValue)
	}

	// Validate content
	if content == "" {
		return uuid.Nil, fmt.Errorf("%w: content cannot be empty", ErrInvalidValue)
	}

	id, err := s.teacherScopeRepo.CreateStudentPost(ctx, teacherUserID, studentID, postType, content)
	if err != nil {
		if errors.Is(err, repo.ErrForbidden) {
			return uuid.Nil, ErrForbidden
		}
		return uuid.Nil, fmt.Errorf("failed to create student post: %w", err)
	}

	return id, nil
}

// ListClassPosts liệt kê bài đăng của một lớp học
func (s *TeacherScopeService) ListClassPosts(ctx context.Context, teacherUserID, classID uuid.UUID,
	limit, offset int) ([]model.Post, int, error) {
	// Validate teacherUserID
	if teacherUserID == uuid.Nil {
		return nil, 0, ErrInvalidUserID
	}

	// Validate classID
	if classID == uuid.Nil {
		return nil, 0, ErrInvalidClassID
	}

	// Default limit
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	posts, total, err := s.teacherScopeRepo.ListClassPosts(ctx, teacherUserID, classID, limit, offset)
	if err != nil {
		if errors.Is(err, repo.ErrForbidden) {
			return nil, 0, ErrForbidden
		}
		return nil, 0, fmt.Errorf("failed to list class posts: %w", err)
	}

	return posts, total, nil
}

// ListStudentPosts liệt kê bài đăng của một học sinh
func (s *TeacherScopeService) ListStudentPosts(ctx context.Context, teacherUserID, studentID uuid.UUID,
	limit, offset int) ([]model.Post, int, error) {
	// Validate teacherUserID
	if teacherUserID == uuid.Nil {
		return nil, 0, ErrInvalidUserID
	}

	// Validate studentID
	if studentID == uuid.Nil {
		return nil, 0, ErrInvalidUserID
	}

	// Default limit
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	posts, total, err := s.teacherScopeRepo.ListStudentPosts(ctx, teacherUserID, studentID, limit, offset)
	if err != nil {
		if errors.Is(err, repo.ErrForbidden) {
			return nil, 0, ErrForbidden
		}
		return nil, 0, fmt.Errorf("failed to list student posts: %w", err)
	}

	return posts, total, nil
}

// isValidPostType kiểm tra postType có hợp lệ không
func isValidPostType(postType string) bool {
	switch postType {
	case "announcement", "activity", "daily_note", "health_note":
		return true
	default:
		return false
	}
}
