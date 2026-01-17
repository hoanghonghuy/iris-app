package service

import (
	"context"

	"github.com/hoanghonghuy/iris-app/apps/api/internal/model"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/repo"
)

type SchoolService struct {
	Repo *repo.SchoolRepo
}

// Create tạo mới trường học
func (s *SchoolService) Create(ctx context.Context, name, addr string) (*model.School, error) {
	id, err := s.Repo.Create(ctx, name, addr)
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
func (s *SchoolService) List(ctx context.Context) ([]model.School, error) {
	return s.Repo.List(ctx)
}
