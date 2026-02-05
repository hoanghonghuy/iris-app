package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/hoanghonghuy/iris-app/apps/api/internal/model"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/repo"
)

type ParentScopeService struct {
	parentScopeRepo *repo.ParentScopeRepo
}

func NewParentScopeService(parentScopeRepo *repo.ParentScopeRepo) *ParentScopeService {
	return &ParentScopeService{
		parentScopeRepo: parentScopeRepo,
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

// ListMyChildClassPosts liệt kê bài đăng của lớp con mình đang học
func (s *ParentScopeService) ListMyChildClassPosts(ctx context.Context, parentUserID, studentID uuid.UUID,
	limit, offset int) ([]model.Post, error) {
	// Validate parentUserID
	if parentUserID == uuid.Nil {
		return nil, ErrInvalidUserID
	}

	// Validate studentID
	if studentID == uuid.Nil {
		return nil, ErrInvalidUserID
	}

	// Kiểm tra xem user có phải là parent của student không
	isParent, err := s.parentScopeRepo.IsParentOfStudent(ctx, parentUserID, studentID)
	if err != nil {
		return nil, fmt.Errorf("failed to verify parent-child relationship: %w", err)
	}
	if !isParent {
		return nil, ErrForbidden
	}

	// Default limit
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	posts, err := s.parentScopeRepo.ListMyChildClassPosts(ctx, parentUserID, studentID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list class posts: %w", err)
	}

	return posts, nil
}

// ListMyChildStudentPosts liệt kê bài đăng riêng của con mình (student scope)
func (s *ParentScopeService) ListMyChildStudentPosts(ctx context.Context, parentUserID, studentID uuid.UUID,
	limit, offset int) ([]model.Post, error) {
	// Validate parentUserID
	if parentUserID == uuid.Nil {
		return nil, ErrInvalidUserID
	}

	// Validate studentID
	if studentID == uuid.Nil {
		return nil, ErrInvalidUserID
	}

	// Kiểm tra xem user có phải là parent của student không
	isParent, err := s.parentScopeRepo.IsParentOfStudent(ctx, parentUserID, studentID)
	if err != nil {
		return nil, fmt.Errorf("failed to verify parent-child relationship: %w", err)
	}
	if !isParent {
		return nil, ErrForbidden
	}

	// Default limit
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	posts, err := s.parentScopeRepo.ListMyChildStudentPosts(ctx, parentUserID, studentID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list student posts: %w", err)
	}

	return posts, nil
}

// ListAllMyChildPosts liệt kê tất cả bài đăng liên quan đến con mình (cả class và student scope)
func (s *ParentScopeService) ListAllMyChildPosts(ctx context.Context, parentUserID, studentID uuid.UUID,
	limit, offset int) ([]model.Post, error) {
	// Validate parentUserID
	if parentUserID == uuid.Nil {
		return nil, ErrInvalidUserID
	}

	// Validate studentID
	if studentID == uuid.Nil {
		return nil, ErrInvalidUserID
	}

	// Kiểm tra xem user có phải là parent của student không
	isParent, err := s.parentScopeRepo.IsParentOfStudent(ctx, parentUserID, studentID)
	if err != nil {
		return nil, fmt.Errorf("failed to verify parent-child relationship: %w", err)
	}
	if !isParent {
		return nil, ErrForbidden
	}

	// Default limit
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	posts, err := s.parentScopeRepo.ListAllMyChildPosts(ctx, parentUserID, studentID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list all child posts: %w", err)
	}

	return posts, nil
}

// GetMyFeed lấy tất cả bài đăng liên quan đến tất cả con của phụ huynh (aggregated feed)
func (s *ParentScopeService) GetMyFeed(ctx context.Context, parentUserID uuid.UUID,
	limit, offset int) ([]model.Post, error) {

	if parentUserID == uuid.Nil {
		return nil, ErrInvalidUserID
	}

	// default limit
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	posts, err := s.parentScopeRepo.GetMyFeed(ctx, parentUserID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get feed: %w", err)
	}

	return posts, nil
}
