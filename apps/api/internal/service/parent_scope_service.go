package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"

	"github.com/hoanghonghuy/iris-app/apps/api/internal/model"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/repo"
)

type ParentScopeService struct {
	parentScopeRepo  *repo.ParentScopeRepo
	postInteractRepo *repo.PostInteractionRepo
}

func NewParentScopeService(parentScopeRepo *repo.ParentScopeRepo, postInteractRepo *repo.PostInteractionRepo) *ParentScopeService {
	return &ParentScopeService{
		parentScopeRepo:  parentScopeRepo,
		postInteractRepo: postInteractRepo,
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
	// Validate parentUserID
	if parentUserID == uuid.Nil {
		return nil, 0, ErrInvalidUserID
	}

	// Validate studentID
	if studentID == uuid.Nil {
		return nil, 0, ErrInvalidUserID
	}

	// Kiểm tra xem user có phải là parent của student không
	isParent, err := s.parentScopeRepo.IsParentOfStudent(ctx, parentUserID, studentID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to verify parent-child relationship: %w", err)
	}
	if !isParent {
		return nil, 0, ErrForbidden
	}

	// Default limit
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	posts, total, err := s.parentScopeRepo.ListMyChildClassPosts(ctx, parentUserID, studentID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list class posts: %w", err)
	}

	return posts, total, nil
}

// ListMyChildStudentPosts liệt kê bài đăng riêng của con phụ huynh (student scope)
func (s *ParentScopeService) ListMyChildStudentPosts(ctx context.Context, parentUserID, studentID uuid.UUID,
	limit, offset int) ([]model.Post, int, error) {
	// Validate parentUserID
	if parentUserID == uuid.Nil {
		return nil, 0, ErrInvalidUserID
	}

	// Validate studentID
	if studentID == uuid.Nil {
		return nil, 0, ErrInvalidUserID
	}

	// Kiểm tra xem user có phải là parent của student không
	isParent, err := s.parentScopeRepo.IsParentOfStudent(ctx, parentUserID, studentID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to verify parent-child relationship: %w", err)
	}
	if !isParent {
		return nil, 0, ErrForbidden
	}

	// Default limit
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	posts, total, err := s.parentScopeRepo.ListMyChildStudentPosts(ctx, parentUserID, studentID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list student posts: %w", err)
	}

	return posts, total, nil
}

// ListAllMyChildPosts liệt kê tất cả bài đăng liên quan đến con của phụ huynh (cả class và student scope)
func (s *ParentScopeService) ListAllMyChildPosts(ctx context.Context, parentUserID, studentID uuid.UUID,
	limit, offset int) ([]model.Post, int, error) {
	// Validate parentUserID
	if parentUserID == uuid.Nil {
		return nil, 0, ErrInvalidUserID
	}

	// Validate studentID
	if studentID == uuid.Nil {
		return nil, 0, ErrInvalidUserID
	}

	// Kiểm tra xem user có phải là parent của student không
	isParent, err := s.parentScopeRepo.IsParentOfStudent(ctx, parentUserID, studentID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to verify parent-child relationship: %w", err)
	}
	if !isParent {
		return nil, 0, ErrForbidden
	}

	// Default limit
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

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

	// default limit
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

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
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

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
