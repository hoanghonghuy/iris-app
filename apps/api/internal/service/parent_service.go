package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/model"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/repo"
)

type ParentService struct {
	parentRepo *repo.ParentRepo
}

func NewParentService(parentRepo *repo.ParentRepo) *ParentService {
	return &ParentService{
		parentRepo: parentRepo,
	}
}

// List lấy danh sách tất cả phụ huynh
func (s *ParentService) List(ctx context.Context) ([]model.Parent, error) {
	return s.parentRepo.List(ctx)
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
