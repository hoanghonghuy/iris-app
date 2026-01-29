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

// Create tạo mới trường học
func (h *SchoolHandler) Create(c *gin.Context) {
	var req CreateSchoolRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid request body")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	school, err := h.schoolService.Create(ctx, req.Name, req.Address)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "failed to create school")
		return
	}

	response.Created(c, gin.H{"school_id": school.ID})
}

// List lấy danh sách trường học
func (h *SchoolHandler) List(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	schools, err := h.schoolService.List(ctx)
	switch err {
	case nil:
		response.OK(c, schools)
		return
	default:
		response.Fail(c, http.StatusInternalServerError, "failed to fetch schools")
		return
	}
}
