package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/model"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/repo"
)

type StudentService struct {
	studentRepo *repo.StudentRepo
	classRepo   *repo.ClassRepo
}

func NewStudentService(studentRepo *repo.StudentRepo, classRepo *repo.ClassRepo) *StudentService {
	return &StudentService{
		studentRepo: studentRepo,
		classRepo:   classRepo,
	}
}

// Create tạo mới học sinh.
func (s *StudentService) Create(ctx context.Context, adminSchoolID *uuid.UUID, schoolID, classID uuid.UUID,
	fullName string, dobStr string, gender string) (*model.Student, error) {
	// SCHOOL_ADMIN: validate school_id trong request phải khớp school của admin
	// adminSchoolID == nil => SUPER_ADMIN: không cần validate
	if adminSchoolID != nil && *adminSchoolID != schoolID {
		return nil, ErrSchoolAccessDenied
	}

	// SCHOOL_ADMIN: validate class thuộc cùng school với admin
	// adminSchoolID == nil => SUPER_ADMIN: không cần validate
	if adminSchoolID != nil {
		class, err := s.classRepo.GetByClassID(ctx, classID)
		if err != nil {
			return nil, ErrInvalidClassID
		}
		if class.SchoolID != *adminSchoolID {
			return nil, ErrSchoolAccessDenied
		}
	}

	// Parse DOB from YYYY-MM-DD string
	dob, err := time.Parse("2006-01-02", dobStr)
	if err != nil {
		return nil, ErrInvalidValue
	}

	id, err := s.studentRepo.Create(ctx, schoolID, classID, fullName, dob, gender)
	if err != nil {
		return nil, err
	}

	return &model.Student{
		StudentID:      id,
		SchoolID:       schoolID,
		CurrentClassID: classID,
		FullName:       fullName,
		DOB:            dob,
		Gender:         gender,
	}, nil
}

// ListByClass lấy danh sách học sinh theo lớp.
func (s *StudentService) ListByClass(ctx context.Context, adminSchoolID *uuid.UUID, classID uuid.UUID, limit, offset int) ([]model.Student, int, error) {
	// SCHOOL_ADMIN: validate class thuộc cùng school với admin
	// adminSchoolID == nil => SUPER_ADMIN: không cần validate
	if adminSchoolID != nil {
		class, err := s.classRepo.GetByClassID(ctx, classID)
		if err != nil {
			return nil, 0, ErrInvalidClassID
		}
		if class.SchoolID != *adminSchoolID {
			return nil, 0, ErrSchoolAccessDenied
		}
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
	return s.studentRepo.ListByClass(ctx, classID, limit, offset)
}
