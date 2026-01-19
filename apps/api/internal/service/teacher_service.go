package service

import (
	"context"

	"github.com/hoanghonghuy/iris-app/apps/api/internal/model"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/repo"
)

type TeacherService struct {
	TeacherRepo *repo.TeacherRepo
}

func (s *TeacherService) List(ctx context.Context) ([]model.Teacher, error) {
	return s.TeacherRepo.List(ctx)
}