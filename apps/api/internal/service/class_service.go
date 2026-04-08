package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/model"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/repo"
	"github.com/jackc/pgx/v5"
)

type ClassService struct {
	classRepo classRepo
}

type classRepo interface {
	Create(ctx context.Context, schoolID uuid.UUID, name, schoolYear string) (uuid.UUID, error)
	List(ctx context.Context, schoolID uuid.UUID, limit, offset int) ([]model.Class, int, error)
	GetByClassID(ctx context.Context, classID uuid.UUID) (*model.Class, error)
	Update(ctx context.Context, classID uuid.UUID, name, schoolYear string) error
	Delete(ctx context.Context, classID uuid.UUID) error
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

// Update cập nhật thông tin lớp học
func (s *ClassService) Update(ctx context.Context, adminSchoolID *uuid.UUID, classID uuid.UUID, name, schoolYear string) error {
	// Kiểm tra class thuộc school nào
	cls, err := s.classRepo.GetByClassID(ctx, classID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrClassNotFound
		}
		return err
	}
	if adminSchoolID != nil && *adminSchoolID != cls.SchoolID {
		return ErrSchoolAccessDenied
	}
	if err := s.classRepo.Update(ctx, classID, name, schoolYear); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrClassNotFound
		}
		return err
	}
	return nil
}

// Delete xóa lớp học
func (s *ClassService) Delete(ctx context.Context, adminSchoolID *uuid.UUID, classID uuid.UUID) error {
	cls, err := s.classRepo.GetByClassID(ctx, classID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrClassNotFound
		}
		return err
	}
	if adminSchoolID != nil && *adminSchoolID != cls.SchoolID {
		return ErrSchoolAccessDenied
	}
	if err := s.classRepo.Delete(ctx, classID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrClassNotFound
		}
		return err
	}
	return nil
}
