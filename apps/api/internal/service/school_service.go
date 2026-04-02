package service

import (
	"context"

	"github.com/google/uuid"
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
func (s *SchoolService) Create(ctx context.Context, name, address string) (*model.School, error) {
	id, err := s.schoolRepo.Create(ctx, name, address)
	if err != nil {
		return nil, err
	}

	return &model.School{
		SchoolID: id,
		Name:     name,
		Address:  address,
	}, nil
}

// List lấy danh sách trường học.
//
// adminSchoolID == nil → tất cả trường (SUPER_ADMIN)
//
// adminSchoolID != nil → chỉ trường của admin đó (SCHOOL_ADMIN)
func (s *SchoolService) List(ctx context.Context, adminSchoolID *uuid.UUID, limit, offset int) ([]model.School, int, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	// SCHOOL_ADMIN: chỉ trả về trường của admin hiện tại đang truy vấn.
	if adminSchoolID != nil {
		school, err := s.schoolRepo.GetByID(ctx, *adminSchoolID)
		if err != nil {
			return nil, 0, err
		}
		return []model.School{*school}, 1, nil
	}

	return s.schoolRepo.List(ctx, limit, offset)
}

// Update cập nhật thông tin trường học (SUPER_ADMIN only)
func (s *SchoolService) Update(ctx context.Context, schoolID uuid.UUID, name, address string) error {
	return s.schoolRepo.Update(ctx, schoolID, name, address)
}

// Delete xóa trường học (SUPER_ADMIN only)
func (s *SchoolService) Delete(ctx context.Context, schoolID uuid.UUID) error {
	return s.schoolRepo.Delete(ctx, schoolID)
}
