package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/model"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/repo"
)

type ClassService struct {
	classRepo *repo.ClassRepo
}

func NewClassService(classRepo *repo.ClassRepo) *ClassService {
	return &ClassService{
		classRepo: classRepo,
	}
}

// Create tạo mới lớp học
func (s *ClassService) Create(ctx context.Context, schoolID uuid.UUID, name, schoolYear string) (*model.Class, error) {
	id, err := s.classRepo.Create(ctx, schoolID, name, schoolYear)
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
	return s.classRepo.List(ctx, schoolID)
}
