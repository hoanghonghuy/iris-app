package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/response"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/service"
)

type ClassHandler struct {
	ClassService *service.ClassService
}

type CreateClassRequest struct {
	Name       string `json:"name" binding:"required,min=1,max=100"`
	SchoolYear string `json:"school_year" binding:"required,min=4,max=20"`
}

// Create tạo mới lớp học
func (h *ClassHandler) Create(c *gin.Context) {
	schoolID, err := uuid.Parse(c.Param("school_id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid school_id format")
		return
	}

	var req CreateClassRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid request body")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	class, err := h.ClassService.Create(ctx, schoolID, req.Name, req.SchoolYear)
	if err != nil {
		// FK fail -> school_id not found
		response.Fail(c, http.StatusInternalServerError, "failed to create class")
		return
	}

	response.Created(c, gin.H{"class_id": class.ID})
}

// ListBySchool lấy danh sách lớp theo trường
func (h *ClassHandler) ListBySchool(c *gin.Context) {
	schoolID, err := uuid.Parse(c.Param("school_id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid school_id format")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	classes, err := h.ClassService.ListBySchool(ctx, schoolID)
	switch err {
	case nil:
		response.OK(c, classes)
		return
	default:
		response.Fail(c, http.StatusInternalServerError, "failed to fetch classes")
		return
	}
}
