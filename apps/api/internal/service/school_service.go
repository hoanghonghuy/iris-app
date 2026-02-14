package service

import (
	"context"

	"github.com/hoanghonghuy/iris-app/apps/api/internal/model"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/repo"
)

type SchoolService struct {
	schoolRepo *repo.SchoolRepo
}

func NewSchoolService(schoolRepo *repo.SchoolRepo) *SchoolService {
	return &SchoolService{
		schoolRepo: schoolRepo,
	}
}

// Create tạo mới trường học
func (s *SchoolService) Create(ctx context.Context, name, addr string) (*model.School, error) {
	id, err := s.schoolRepo.Create(ctx, name, addr)
	if err != nil {
		return nil, err
	}

	return &model.School{
		ID:      id,
		Name:    name,
		Address: addr,
	}, nil
}

// List lấy danh sách tất cả trường học
func (s *SchoolService) List(ctx context.Context, limit, offset int) ([]model.School, int, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}
	return s.schoolRepo.List(ctx, limit, offset)
}
