package handlers

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/response"
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

// Create tạo mới lớp học
func (h *ClassHandler) Create(c *gin.Context) {
	adminSchoolID := extractAdminSchoolID(c)

	var req CreateClassRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid request body")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	class, err := h.classService.Create(ctx, adminSchoolID, req.SchoolID, req.Name, req.SchoolYear)
	if err != nil {
		if errors.Is(err, service.ErrSchoolAccessDenied) {
			response.Fail(c, http.StatusForbidden, "access denied")
			return
		}
		response.Fail(c, http.StatusInternalServerError, "failed to create class")
		return
	}

	response.Created(c, gin.H{"class_id": class.ClassID})
}

// ListBySchool lấy danh sách lớp theo trường
func (h *ClassHandler) ListBySchool(c *gin.Context) {
	adminSchoolID := extractAdminSchoolID(c)

	schoolID, err := uuid.Parse(c.Param("school_id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid school_id format")
		return
	}

	var params PaginationParams
	if err := c.ShouldBindQuery(&params); err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid pagination params")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	classes, total, err := h.classService.ListBySchool(ctx, adminSchoolID, schoolID, params.Limit, params.Offset)
	if err != nil {
		if errors.Is(err, service.ErrSchoolAccessDenied) {
			response.Fail(c, http.StatusForbidden, "access denied")
			return
		}
		response.Fail(c, http.StatusInternalServerError, "failed to fetch classes")
		return
	}

	response.OKPaginated(c, classes, response.Pagination{
		Total:   total,
		Limit:   params.Limit,
		Offset:  params.Offset,
		HasMore: params.Offset+len(classes) < total,
	})
}
