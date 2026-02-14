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
}

func NewStudentService(studentRepo *repo.StudentRepo) *StudentService {
	return &StudentService{
		studentRepo: studentRepo,
	}
}

func (s *StudentService) Create(ctx context.Context, schoolID, classID uuid.UUID,
	fullName string, dobStr string, gender string) (*model.Student, error) {
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
		ID:             id,
		SchoolID:       schoolID,
		CurrentClassID: classID,
		FullName:       fullName,
		DOB:            dob,
		Gender:         gender,
	}, nil
}

func (s *StudentService) ListByClass(ctx context.Context, classID uuid.UUID, limit, offset int) ([]model.Student, int, error) {
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
