package parenthandlers

import (
	"context"

	"github.com/google/uuid"

	"github.com/hoanghonghuy/iris-app/apps/api/internal/model"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/service"
)

type parentService interface {
	List(ctx context.Context, adminSchoolID *uuid.UUID, limit, offset int) ([]model.Parent, int, error)
	GetByParentID(ctx context.Context, adminSchoolID *uuid.UUID, parentID uuid.UUID) (*model.Parent, error)
	AssignStudent(ctx context.Context, adminSchoolID *uuid.UUID, parentID, studentID uuid.UUID, relationship string) error
	UnassignStudent(ctx context.Context, adminSchoolID *uuid.UUID, parentID, studentID uuid.UUID) error
}

type ParentHandler struct {
	parentService parentService
}

func NewParentHandler(parentService *service.ParentService) *ParentHandler {
	return &ParentHandler{
		parentService: parentService,
	}
}

type AssignStudentRequest struct {
	Relationship string `json:"relationship"`
}
