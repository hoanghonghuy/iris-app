package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/hoanghonghuy/iris-app/apps/api/internal/model"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/repo"
)

type teacherScopeServiceTeacherScopeRepo interface {
	teacherScopeStatsRepo
	teacherScopeAttendanceRepo
	teacherScopePostRepo
}

type teacherScopeServiceHealthLogRepo interface {
	teacherScopeHealthLogRepo
}

type teacherScopeServiceTeacherRepo interface {
	teacherScopeProfileRepo
}

type teacherScopeServicePostInteractionRepo interface {
	teacherScopeInteractionRepo
}

type TeacherScopeService struct {
	teacherScopeRepo teacherScopeServiceTeacherScopeRepo
	healthLogRepo    teacherScopeServiceHealthLogRepo
	teacherRepo      teacherScopeServiceTeacherRepo
	postInteractRepo teacherScopeServicePostInteractionRepo
}

func NewTeacherScopeService(teacherScopeRepo *repo.TeacherScopeRepo, healthLogRepo *repo.HealthLogRepo, teacherRepo *repo.TeacherRepo, postInteractRepo *repo.PostInteractionRepo) *TeacherScopeService {
	return &TeacherScopeService{
		teacherScopeRepo: teacherScopeRepo,
		healthLogRepo:    healthLogRepo,
		teacherRepo:      teacherRepo,
		postInteractRepo: postInteractRepo,
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
		if errors.Is(err, repo.ErrNoRowsUpdated) {
			return ErrForbidden
		}
		return fmt.Errorf("failed to mark attendance: %w", err)
	}

	return nil
}

// CancelAttendanceForDate hủy điểm danh của học sinh trong ngày.
func (s *TeacherScopeService) CancelAttendanceForDate(ctx context.Context, teacherUserID, studentID uuid.UUID, date string) error {
	if teacherUserID == uuid.Nil {
		return ErrInvalidUserID
	}

	if studentID == uuid.Nil {
		return ErrInvalidUserID
	}

	parsedDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return ErrInvalidDate
	}

	err = s.teacherScopeRepo.DeleteAttendanceForDate(ctx, teacherUserID, studentID, parsedDate)
	if err != nil {
		if errors.Is(err, repo.ErrNoRowsUpdated) {
			return ErrForbidden
		}
		return fmt.Errorf("failed to cancel attendance: %w", err)
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

// ListAttendanceChangeLogsByStudent liệt kê lịch sử chỉnh sửa điểm danh của một học sinh.
func (s *TeacherScopeService) ListAttendanceChangeLogsByStudent(ctx context.Context, teacherUserID, studentID uuid.UUID,
	from, to time.Time) ([]model.AttendanceChangeLog, error) {

	if teacherUserID == uuid.Nil {
		return nil, ErrInvalidUserID
	}

	if studentID == uuid.Nil {
		return nil, ErrInvalidUserID
	}

	logs, err := s.teacherScopeRepo.ListAttendanceChangeLogsByStudent(ctx, teacherUserID, studentID, from, to)
	if err != nil {
		return nil, fmt.Errorf("failed to list attendance change logs: %w", err)
	}

	return logs, nil
}

// ListAttendanceChangeLogsByClass liệt kê lịch sử chỉnh sửa điểm danh theo lớp có phân trang.
func (s *TeacherScopeService) ListAttendanceChangeLogsByClass(ctx context.Context, teacherUserID, classID uuid.UUID,
	studentID *uuid.UUID, status *string, from, to time.Time, limit, offset int) ([]model.AttendanceChangeLog, int, error) {
	if teacherUserID == uuid.Nil {
		return nil, 0, ErrInvalidUserID
	}

	if classID == uuid.Nil {
		return nil, 0, ErrInvalidClassID
	}

	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	if status != nil {
		if !isValidAttendanceStatus(*status) {
			return nil, 0, ErrInvalidStatus
		}
	}

	logs, total, err := s.teacherScopeRepo.ListAttendanceChangeLogsByClass(ctx, teacherUserID, classID, studentID, status, from, to, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list class attendance change logs: %w", err)
	}

	return logs, total, nil
}

func isValidAttendanceStatus(status string) bool {
	return status == "present" || status == "absent" || status == "late" || status == "excused"
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
	id, err := s.healthLogRepo.CreateByStudentAndTeacher(ctx, teacherUserID, studentID, recordedAt, temperature, symptoms, severity, note)
	if err != nil {
		if errors.Is(err, repo.ErrNoRowsUpdated) {
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

	healthLogs, err := s.healthLogRepo.ListByStudentAndTeacher(ctx, teacherUserID, studentID, from, to)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) || errors.Is(err, repo.ErrNoRowsUpdated) {
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
		if errors.Is(err, repo.ErrNoRowsUpdated) {
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
		if errors.Is(err, repo.ErrNoRowsUpdated) {
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
		if errors.Is(err, repo.ErrNoRowsUpdated) {
			return []model.Post{}, 0, nil
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
		if errors.Is(err, repo.ErrNoRowsUpdated) {
			return []model.Post{}, 0, nil
		}
		return nil, 0, fmt.Errorf("failed to list student posts: %w", err)
	}

	return posts, total, nil
}

// UpdatePost cập nhật nội dung bài đăng của chính giáo viên và lưu lịch sử trước/sau chỉnh sửa.
func (s *TeacherScopeService) UpdatePost(ctx context.Context, teacherUserID, postID uuid.UUID, content string) error {
	if teacherUserID == uuid.Nil {
		return ErrInvalidUserID
	}

	if postID == uuid.Nil {
		return ErrInvalidValue
	}

	trimmedContent := strings.TrimSpace(content)
	if trimmedContent == "" {
		return fmt.Errorf("%w: content cannot be empty", ErrInvalidValue)
	}

	err := s.teacherScopeRepo.UpdatePost(ctx, teacherUserID, postID, trimmedContent)
	if err != nil {
		if errors.Is(err, repo.ErrNoRowsUpdated) {
			return ErrForbidden
		}
		return fmt.Errorf("failed to update post: %w", err)
	}

	return nil
}

// DeletePost xóa bài đăng của chính giáo viên.
func (s *TeacherScopeService) DeletePost(ctx context.Context, teacherUserID, postID uuid.UUID) error {
	if teacherUserID == uuid.Nil {
		return ErrInvalidUserID
	}

	if postID == uuid.Nil {
		return ErrInvalidValue
	}

	err := s.teacherScopeRepo.DeletePost(ctx, teacherUserID, postID)
	if err != nil {
		if errors.Is(err, repo.ErrNoRowsUpdated) {
			return ErrForbidden
		}
		return fmt.Errorf("failed to delete post: %w", err)
	}

	return nil
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

// TogglePostLike bật/tắt like cho bài đăng mà giáo viên có quyền truy cập.
func (s *TeacherScopeService) TogglePostLike(ctx context.Context, teacherUserID, postID uuid.UUID) (bool, int, error) {
	if teacherUserID == uuid.Nil {
		return false, 0, ErrInvalidUserID
	}
	if postID == uuid.Nil {
		return false, 0, ErrInvalidValue
	}

	allowed, err := s.postInteractRepo.TeacherCanAccessPost(ctx, teacherUserID, postID)
	if err != nil {
		return false, 0, fmt.Errorf("failed to verify post access: %w", err)
	}
	if !allowed {
		return false, 0, ErrForbidden
	}

	liked, likeCount, err := s.postInteractRepo.ToggleLike(ctx, teacherUserID, postID)
	if err != nil {
		return false, 0, fmt.Errorf("failed to toggle like: %w", err)
	}

	return liked, likeCount, nil
}

// AddPostComment thêm bình luận cho bài đăng mà giáo viên có quyền truy cập.
func (s *TeacherScopeService) AddPostComment(ctx context.Context, teacherUserID, postID uuid.UUID, content string) (model.PostComment, int, error) {
	if teacherUserID == uuid.Nil {
		return model.PostComment{}, 0, ErrInvalidUserID
	}
	if postID == uuid.Nil {
		return model.PostComment{}, 0, ErrInvalidValue
	}
	trimmedContent := strings.TrimSpace(content)
	if trimmedContent == "" {
		return model.PostComment{}, 0, fmt.Errorf("%w: content cannot be empty", ErrInvalidValue)
	}

	allowed, err := s.postInteractRepo.TeacherCanAccessPost(ctx, teacherUserID, postID)
	if err != nil {
		return model.PostComment{}, 0, fmt.Errorf("failed to verify post access: %w", err)
	}
	if !allowed {
		return model.PostComment{}, 0, ErrForbidden
	}

	comment, err := s.postInteractRepo.AddComment(ctx, teacherUserID, postID, trimmedContent)
	if err != nil {
		return model.PostComment{}, 0, fmt.Errorf("failed to add comment: %w", err)
	}

	commentCount, err := s.postInteractRepo.CountComments(ctx, postID)
	if err != nil {
		return model.PostComment{}, 0, fmt.Errorf("failed to count comments: %w", err)
	}

	return comment, commentCount, nil
}

// ListPostComments liệt kê bình luận của bài đăng mà giáo viên có quyền truy cập.
func (s *TeacherScopeService) ListPostComments(ctx context.Context, teacherUserID, postID uuid.UUID, limit, offset int) ([]model.PostComment, int, error) {
	if teacherUserID == uuid.Nil {
		return nil, 0, ErrInvalidUserID
	}
	if postID == uuid.Nil {
		return nil, 0, ErrInvalidValue
	}
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	allowed, err := s.postInteractRepo.TeacherCanAccessPost(ctx, teacherUserID, postID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to verify post access: %w", err)
	}
	if !allowed {
		return nil, 0, ErrForbidden
	}

	comments, total, err := s.postInteractRepo.ListComments(ctx, postID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list comments: %w", err)
	}

	return comments, total, nil
}

// SharePost ghi nhận một lượt chia sẻ cho bài đăng mà giáo viên có quyền truy cập.
func (s *TeacherScopeService) SharePost(ctx context.Context, teacherUserID, postID uuid.UUID) (int, error) {
	if teacherUserID == uuid.Nil {
		return 0, ErrInvalidUserID
	}
	if postID == uuid.Nil {
		return 0, ErrInvalidValue
	}

	allowed, err := s.postInteractRepo.TeacherCanAccessPost(ctx, teacherUserID, postID)
	if err != nil {
		return 0, fmt.Errorf("failed to verify post access: %w", err)
	}
	if !allowed {
		return 0, ErrForbidden
	}

	shareCount, err := s.postInteractRepo.AddShare(ctx, teacherUserID, postID)
	if err != nil {
		return 0, fmt.Errorf("failed to add share: %w", err)
	}

	return shareCount, nil
}
