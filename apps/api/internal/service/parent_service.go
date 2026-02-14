package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/model"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/repo"
)

type ParentService struct {
	parentRepo        *repo.ParentRepo
	studentParentRepo *repo.StudentParentRepo
}

func NewParentService(parentRepo *repo.ParentRepo, studentParentRepo *repo.StudentParentRepo) *ParentService {
	return &ParentService{
		parentRepo:        parentRepo,
		studentParentRepo: studentParentRepo,
	}
}

// List lấy danh sách tất cả phụ huynh
func (s *ParentService) List(ctx context.Context, limit, offset int) ([]model.Parent, int, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}
	return s.parentRepo.List(ctx, limit, offset)
}

// Create tạo mới phụ huynh
func (s *ParentService) Create(ctx context.Context, userID, schoolID uuid.UUID, fullName, phone string) (*model.Parent, error) {
	id, err := s.parentRepo.Create(ctx, userID, schoolID, fullName, phone)
	if err != nil {
		return nil, err
	}

	return &model.Parent{
		ParentID: id,
		UserID:   userID,
		SchoolID: schoolID,
		FullName: fullName,
		Phone:    phone,
	}, nil
}

// GetByParentID lấy thông tin phụ huynh theo parent_id
func (s *ParentService) GetByParentID(ctx context.Context, parentID uuid.UUID) (*model.Parent, error) {
	return s.parentRepo.GetByParentID(ctx, parentID)
}

// AssignStudent gán phụ huynh cho học sinh
func (s *ParentService) AssignStudent(ctx context.Context, parentID, studentID uuid.UUID, relationship string) error {
	if parentID == uuid.Nil {
		return ErrInvalidUserID
	}

	if studentID == uuid.Nil {
		return ErrInvalidUserID
	}

	// Validate relationship
	// Relationship examples: "father", "mother", "guardian", etc.

	return s.studentParentRepo.Assign(ctx, studentID, parentID, relationship)
}

// UnassignStudent hủy gán phụ huynh khỏi học sinh
func (s *ParentService) UnassignStudent(ctx context.Context, parentID, studentID uuid.UUID) error {
	if parentID == uuid.Nil {
		return ErrInvalidUserID
	}

	if studentID == uuid.Nil {
		return ErrInvalidUserID
	}

	return s.studentParentRepo.Unassign(ctx, studentID, parentID)
}
