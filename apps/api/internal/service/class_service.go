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

// Create tạo mới lớp học.
func (s *ClassService) Create(ctx context.Context, adminSchoolID *uuid.UUID, schoolID uuid.UUID, name, schoolYear string) (*model.Class, error) {
	// SCHOOL_ADMIN: validate school_id trong request phải khớp school của admin
	// adminSchoolID == nil => SUPER_ADMIN: không cần validate
	if adminSchoolID != nil && *adminSchoolID != schoolID {
		return nil, ErrSchoolAccessDenied
	}

	id, err := s.classRepo.Create(ctx, schoolID, name, schoolYear)
	if err != nil {
		return nil, err
	}

	return &model.Class{
		ClassID:    id,
		SchoolID:   schoolID,
		Name:       name,
		SchoolYear: schoolYear,
	}, nil
}

// ListBySchool lấy danh sách lớp học theo trường
func (s *ClassService) ListBySchool(ctx context.Context, adminSchoolID *uuid.UUID, schoolID uuid.UUID, limit, offset int) ([]model.Class, int, error) {
	// SCHOOL_ADMIN: validate school_id param phải khớp school của admin
	// adminSchoolID == nil => SUPER_ADMIN: không cần validate
	if adminSchoolID != nil && *adminSchoolID != schoolID {
		return nil, 0, ErrSchoolAccessDenied
	}

	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}
	return s.classRepo.List(ctx, schoolID, limit, offset)
}
