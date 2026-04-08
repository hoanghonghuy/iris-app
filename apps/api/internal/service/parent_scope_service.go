package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/hoanghonghuy/iris-app/apps/api/internal/model"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/repo"
)

type ParentScopeService struct {
	parentScopeRepo  parentScopeRepo
	postInteractRepo postInteractionRepo
	appointmentRepo  appointmentRepo
}

type parentScopeRepo interface {
	IsParentOfStudent(ctx context.Context, parentUserID, studentID uuid.UUID) (bool, error)
	ListMyChildren(ctx context.Context, parentUserID uuid.UUID) ([]model.Student, error)
	ListMyChildClassPosts(ctx context.Context, parentUserID, studentID uuid.UUID, limit, offset int) ([]model.Post, int, error)
	ListMyChildStudentPosts(ctx context.Context, parentUserID, studentID uuid.UUID, limit, offset int) ([]model.Post, int, error)
	ListAllMyChildPosts(ctx context.Context, parentUserID, studentID uuid.UUID, limit, offset int) ([]model.Post, int, error)
	GetMyFeed(ctx context.Context, parentUserID uuid.UUID, limit, offset int) ([]model.Post, int, error)
	CountMyChildren(ctx context.Context, parentUserID uuid.UUID) (int, error)
	CountMyRecentPosts(ctx context.Context, parentUserID uuid.UUID, since time.Time) (int, error)
	CountMyRecentHealthAlerts(ctx context.Context, parentUserID uuid.UUID, since time.Time) (int, error)
}

type postInteractionRepo interface {
	ParentCanAccessPost(ctx context.Context, parentUserID, postID uuid.UUID) (bool, error)
	ToggleLike(ctx context.Context, userID, postID uuid.UUID) (bool, int, error)
	AddComment(ctx context.Context, userID, postID uuid.UUID, content string) (model.PostComment, error)
	CountComments(ctx context.Context, postID uuid.UUID) (int, error)
	ListComments(ctx context.Context, postID uuid.UUID, limit, offset int) ([]model.PostComment, int, error)
	AddShare(ctx context.Context, userID, postID uuid.UUID) (int, error)
}

type appointmentRepo interface {
	CountParentUpcomingAppointments(ctx context.Context, parentUserID uuid.UUID) (int, error)
}

// normalizeParentScopeLimit chuẩn hóa limit mặc định và giới hạn max cho các API listing.
func normalizeParentScopeLimit(limit int) int {
	if limit <= 0 {
		return 20
	}
	if limit > 100 {
		return 100
	}
	return limit
}

// ensureParentStudentAccess xác thực parent-user có quan hệ hợp lệ với student trước khi truy cập dữ liệu theo student.
func (s *ParentScopeService) ensureParentStudentAccess(ctx context.Context, parentUserID, studentID uuid.UUID) error {
	if parentUserID == uuid.Nil {
		return ErrInvalidUserID
	}

	if studentID == uuid.Nil {
		return ErrInvalidUserID
	}

	// Kiểm tra xem user có phải là parent của student không
	isParent, err := s.parentScopeRepo.IsParentOfStudent(ctx, parentUserID, studentID)
	if err != nil {
		return fmt.Errorf("failed to verify parent-child relationship: %w", err)
	}
	if !isParent {
		return ErrForbidden
	}

	return nil
}

func NewParentScopeService(parentScopeRepo *repo.ParentScopeRepo, postInteractRepo *repo.PostInteractionRepo, appointmentRepo *repo.AppointmentRepo) *ParentScopeService {
	return &ParentScopeService{
		parentScopeRepo:  parentScopeRepo,
		postInteractRepo: postInteractRepo,
		appointmentRepo:  appointmentRepo,
	}
}

// ListMyChildren liệt kê các học sinh (con) của phụ huynh
func (s *ParentScopeService) ListMyChildren(ctx context.Context, parentUserID uuid.UUID) ([]model.Student, error) {
	// Validate parentUserID
	if parentUserID == uuid.Nil {
		return nil, ErrInvalidUserID
	}

	students, err := s.parentScopeRepo.ListMyChildren(ctx, parentUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to list children: %w", err)
	}

	return students, nil
}

// ListMyChildClassPosts liệt kê bài đăng của lớp con của phụ huynh đang học
func (s *ParentScopeService) ListMyChildClassPosts(ctx context.Context, parentUserID, studentID uuid.UUID,
	limit, offset int) ([]model.Post, int, error) {
	if err := s.ensureParentStudentAccess(ctx, parentUserID, studentID); err != nil {
		return nil, 0, err
	}

	limit = normalizeParentScopeLimit(limit)

	posts, total, err := s.parentScopeRepo.ListMyChildClassPosts(ctx, parentUserID, studentID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list class posts: %w", err)
	}

	return posts, total, nil
}

// ListMyChildStudentPosts liệt kê bài đăng riêng của con phụ huynh (student scope)
func (s *ParentScopeService) ListMyChildStudentPosts(ctx context.Context, parentUserID, studentID uuid.UUID,
	limit, offset int) ([]model.Post, int, error) {
	if err := s.ensureParentStudentAccess(ctx, parentUserID, studentID); err != nil {
		return nil, 0, err
	}

	limit = normalizeParentScopeLimit(limit)

	posts, total, err := s.parentScopeRepo.ListMyChildStudentPosts(ctx, parentUserID, studentID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list student posts: %w", err)
	}

	return posts, total, nil
}

// ListAllMyChildPosts liệt kê tất cả bài đăng liên quan đến con của phụ huynh (cả class và student scope)
func (s *ParentScopeService) ListAllMyChildPosts(ctx context.Context, parentUserID, studentID uuid.UUID,
	limit, offset int) ([]model.Post, int, error) {
	if err := s.ensureParentStudentAccess(ctx, parentUserID, studentID); err != nil {
		return nil, 0, err
	}

	limit = normalizeParentScopeLimit(limit)

	posts, total, err := s.parentScopeRepo.ListAllMyChildPosts(ctx, parentUserID, studentID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list all child posts: %w", err)
	}

	return posts, total, nil
}

// GetMyFeed lấy tất cả bài đăng liên quan đến tất cả con của phụ huynh (aggregated feed)
func (s *ParentScopeService) GetMyFeed(ctx context.Context, parentUserID uuid.UUID,
	limit, offset int) ([]model.Post, int, error) {

	if parentUserID == uuid.Nil {
		return nil, 0, ErrInvalidUserID
	}

	limit = normalizeParentScopeLimit(limit)

	posts, total, err := s.parentScopeRepo.GetMyFeed(ctx, parentUserID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get feed: %w", err)
	}

	return posts, total, nil
}

// TogglePostLike bật/tắt like cho bài đăng mà phụ huynh có quyền truy cập.
func (s *ParentScopeService) TogglePostLike(ctx context.Context, parentUserID, postID uuid.UUID) (bool, int, error) {
	if parentUserID == uuid.Nil {
		return false, 0, ErrInvalidUserID
	}
	if postID == uuid.Nil {
		return false, 0, ErrInvalidValue
	}

	allowed, err := s.postInteractRepo.ParentCanAccessPost(ctx, parentUserID, postID)
	if err != nil {
		return false, 0, fmt.Errorf("failed to verify post access: %w", err)
	}
	if !allowed {
		return false, 0, ErrForbidden
	}

	liked, likeCount, err := s.postInteractRepo.ToggleLike(ctx, parentUserID, postID)
	if err != nil {
		return false, 0, fmt.Errorf("failed to toggle like: %w", err)
	}

	return liked, likeCount, nil
}

// AddPostComment thêm bình luận cho bài đăng mà phụ huynh có quyền truy cập.
func (s *ParentScopeService) AddPostComment(ctx context.Context, parentUserID, postID uuid.UUID, content string) (model.PostComment, int, error) {
	if parentUserID == uuid.Nil {
		return model.PostComment{}, 0, ErrInvalidUserID
	}
	if postID == uuid.Nil {
		return model.PostComment{}, 0, ErrInvalidValue
	}
	trimmedContent := strings.TrimSpace(content)
	if trimmedContent == "" {
		return model.PostComment{}, 0, fmt.Errorf("%w: content cannot be empty", ErrInvalidValue)
	}

	allowed, err := s.postInteractRepo.ParentCanAccessPost(ctx, parentUserID, postID)
	if err != nil {
		return model.PostComment{}, 0, fmt.Errorf("failed to verify post access: %w", err)
	}
	if !allowed {
		return model.PostComment{}, 0, ErrForbidden
	}

	comment, err := s.postInteractRepo.AddComment(ctx, parentUserID, postID, trimmedContent)
	if err != nil {
		return model.PostComment{}, 0, fmt.Errorf("failed to add comment: %w", err)
	}

	commentCount, err := s.postInteractRepo.CountComments(ctx, postID)
	if err != nil {
		return model.PostComment{}, 0, fmt.Errorf("failed to count comments: %w", err)
	}

	return comment, commentCount, nil
}

// ListPostComments liệt kê bình luận của bài đăng mà phụ huynh có quyền truy cập.
func (s *ParentScopeService) ListPostComments(ctx context.Context, parentUserID, postID uuid.UUID, limit, offset int) ([]model.PostComment, int, error) {
	if parentUserID == uuid.Nil {
		return nil, 0, ErrInvalidUserID
	}
	if postID == uuid.Nil {
		return nil, 0, ErrInvalidValue
	}
	limit = normalizeParentScopeLimit(limit)

	allowed, err := s.postInteractRepo.ParentCanAccessPost(ctx, parentUserID, postID)
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

// SharePost ghi nhận một lượt chia sẻ cho bài đăng mà phụ huynh có quyền truy cập.
func (s *ParentScopeService) SharePost(ctx context.Context, parentUserID, postID uuid.UUID) (int, error) {
	if parentUserID == uuid.Nil {
		return 0, ErrInvalidUserID
	}
	if postID == uuid.Nil {
		return 0, ErrInvalidValue
	}

	allowed, err := s.postInteractRepo.ParentCanAccessPost(ctx, parentUserID, postID)
	if err != nil {
		return 0, fmt.Errorf("failed to verify post access: %w", err)
	}
	if !allowed {
		return 0, ErrForbidden
	}

	shareCount, err := s.postInteractRepo.AddShare(ctx, parentUserID, postID)
	if err != nil {
		return 0, fmt.Errorf("failed to add share: %w", err)
	}

	return shareCount, nil
}

func (s *ParentScopeService) GetMyAnalytics(ctx context.Context, parentUserID uuid.UUID) (*model.ParentAnalytics, error) {
	if parentUserID == uuid.Nil {
		return nil, ErrInvalidUserID
	}

	children, err := s.parentScopeRepo.CountMyChildren(ctx, parentUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to count children: %w", err)
	}

	upcoming, err := s.appointmentRepo.CountParentUpcomingAppointments(ctx, parentUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to count upcoming appointments: %w", err)
	}

	since := time.Now().AddDate(0, 0, -7)
	recentPosts, err := s.parentScopeRepo.CountMyRecentPosts(ctx, parentUserID, since)
	if err != nil {
		return nil, fmt.Errorf("failed to count recent posts: %w", err)
	}

	healthAlerts, err := s.parentScopeRepo.CountMyRecentHealthAlerts(ctx, parentUserID, since)
	if err != nil {
		return nil, fmt.Errorf("failed to count health alerts: %w", err)
	}

	return &model.ParentAnalytics{
		TotalChildren:        children,
		UpcomingAppointments: upcoming,
		RecentPosts7d:        recentPosts,
		RecentHealthAlerts7d: healthAlerts,
	}, nil
}
