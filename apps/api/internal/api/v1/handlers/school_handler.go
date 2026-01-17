package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/response"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/service"
)

type SchoolHandler struct {
	SchoolService *service.SchoolService
}

type CreateSchoolRequest struct {
	Name    string `json:"name" binding:"required,min=2"`
	Address string `json:"address"`
}

// Create tạo mới trường học
func (h *SchoolHandler) Create(c *gin.Context) {
	var req CreateSchoolRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid payload")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	school, err := h.SchoolService.Create(ctx, req.Name, req.Address)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "db error")
		return
	}

	response.Created(c, gin.H{"school_id": school.ID})
}

// List lấy danh sách trường học
func (h *SchoolHandler) List(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	schools, err := h.SchoolService.List(ctx)
	switch err {
	case nil:
		response.OK(c, schools)
		return
	default:
		response.Fail(c, http.StatusInternalServerError, "failed to fetch schools")
		return
	}
}
