package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/model"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/repo"
)

type ClassService struct {
	ClassRepo *repo.ClassRepo
}

// Create tạo mới lớp học
func (s *ClassService) Create(ctx context.Context, schoolID uuid.UUID, name, schoolYear string) (*model.Class, error) {
	id, err := s.ClassRepo.Create(ctx, schoolID, name, schoolYear)
	if err != nil {
		return nil, err
	}

	return &model.Class{
		ID:         id,
		SchoolID:   schoolID,
		Name:       name,
		SchoolYear: schoolYear,
	}, nil
}

// ListBySchool lấy danh sách lớp học theo trường
func (s *ClassService) ListBySchool(ctx context.Context, schoolID uuid.UUID) ([]model.Class, error) {
	return s.ClassRepo.List(ctx, schoolID)
}
