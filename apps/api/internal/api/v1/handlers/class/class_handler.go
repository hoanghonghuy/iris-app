package classhandlers

import (
	"github.com/google/uuid"

	"github.com/hoanghonghuy/iris-app/apps/api/internal/service"
)

type ClassHandler struct {
	classService *service.ClassService
}

func NewClassHandler(classService *service.ClassService) *ClassHandler {
	return &ClassHandler{
		classService: classService,
	}
}

type CreateClassRequest struct {
	SchoolID   uuid.UUID `json:"school_id" binding:"required"`
	Name       string    `json:"name" binding:"required,min=1,max=100"`
	SchoolYear string    `json:"school_year" binding:"required,min=4,max=20"`
}

// UpdateClassRequest input để cập nhật lớp học.
type UpdateClassRequest struct {
	Name       string `json:"name" binding:"required,min=1,max=100"`
	SchoolYear string `json:"school_year" binding:"required,min=4,max=20"`
}
