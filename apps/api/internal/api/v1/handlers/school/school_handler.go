package schoolhandlers

import "github.com/hoanghonghuy/iris-app/apps/api/internal/service"

type SchoolHandler struct {
	schoolService *service.SchoolService
}

func NewSchoolHandler(schoolService *service.SchoolService) *SchoolHandler {
	return &SchoolHandler{
		schoolService: schoolService,
	}
}

type CreateSchoolRequest struct {
	Name    string `json:"name" binding:"required,min=2"`
	Address string `json:"address"`
}

// UpdateSchoolRequest input để cập nhật trường học.
type UpdateSchoolRequest struct {
	Name    string `json:"name" binding:"required,min=2"`
	Address string `json:"address"`
}
