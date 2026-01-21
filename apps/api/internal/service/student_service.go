package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/model"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/repo"
)

type StudentService struct {
	StudentRepo *repo.StudentRepo
}

func (s *StudentService) Create(ctx context.Context, schoolID, classID uuid.UUID,
	fullName string, dob time.Time, gender string) (*model.Student, error) {
	id, err := s.StudentRepo.Create(ctx, schoolID, classID, fullName, dob, gender)
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

func (s *StudentService) ListByClass(ctx context.Context, classID uuid.UUID) ([]model.Student, error) {
	return s.StudentRepo.ListByClass(ctx, classID)
}
