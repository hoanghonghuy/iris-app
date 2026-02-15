package handlers

import (
	"context"
	"errors"
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

	response.Created(c, gin.H{"school_id": school.SchoolID})
}

// List lấy danh sách trường học
func (h *SchoolHandler) List(c *gin.Context) {
	adminSchoolID := extractAdminSchoolID(c)

	var params PaginationParams
	if err := c.ShouldBindQuery(&params); err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid pagination params")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	schools, total, err := h.schoolService.List(ctx, adminSchoolID, params.Limit, params.Offset)
	if err != nil {
		if errors.Is(err, service.ErrSchoolAccessDenied) {
			response.Fail(c, http.StatusForbidden, "access denied")
			return
		}
		response.Fail(c, http.StatusInternalServerError, "failed to fetch schools")
		return
	}

	response.OKPaginated(c, schools, response.Pagination{
		Total:   total,
		Limit:   params.Limit,
		Offset:  params.Offset,
		HasMore: params.Offset+len(schools) < total,
	})
}
