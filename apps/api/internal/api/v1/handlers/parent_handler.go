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
	"github.com/jackc/pgx/v5"
)

type ParentHandler struct {
	parentService *service.ParentService
}

func NewParentHandler(parentService *service.ParentService) *ParentHandler {
	return &ParentHandler{
		parentService: parentService,
	}
}

type AssignStudentRequest struct {
	Relationship string `json:"relationship"` // father, mother, guardian, etc. (optional)
}

// AssignStudent gán phụ huynh cho học sinh (admin only)
func (h *ParentHandler) AssignStudent(c *gin.Context) {
	adminSchoolID := extractAdminSchoolID(c)

	parentID, err := uuid.Parse(c.Param("parent_id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid parent_id format")
		return
	}

	studentID, err := uuid.Parse(c.Param("student_id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid student_id format")
		return
	}

	var req AssignStudentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// Relationship is optional, so binding error is ok
		req.Relationship = ""
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	err = h.parentService.AssignStudent(ctx, adminSchoolID, parentID, studentID, req.Relationship)
	if err != nil {
		if errors.Is(err, service.ErrSchoolAccessDenied) {
			response.Fail(c, http.StatusForbidden, "access denied")
			return
		}
		response.Fail(c, http.StatusInternalServerError, "failed to assign parent to student")
		return
	}

	response.OK(c, gin.H{
		"message":      "parent assigned to student successfully",
		"parent_id":    parentID.String(),
		"student_id":   studentID.String(),
		"relationship": req.Relationship,
	})
}

// UnassignStudent hủy gán phụ huynh khỏi học sinh (admin only)
func (h *ParentHandler) UnassignStudent(c *gin.Context) {
	adminSchoolID := extractAdminSchoolID(c)

	parentID, err := uuid.Parse(c.Param("parent_id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid parent_id format")
		return
	}

	studentID, err := uuid.Parse(c.Param("student_id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid student_id format")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	err = h.parentService.UnassignStudent(ctx, adminSchoolID, parentID, studentID)
	if err != nil {
		if errors.Is(err, service.ErrSchoolAccessDenied) {
			response.Fail(c, http.StatusForbidden, "access denied")
			return
		}
		response.Fail(c, http.StatusInternalServerError, "failed to unassign parent from student")
		return
	}

	response.OK(c, gin.H{
		"message":    "parent unassigned from student successfully",
		"parent_id":  parentID.String(),
		"student_id": studentID.String(),
	})
}

// List lấy danh sách phụ huynh (admin only)
func (h *ParentHandler) List(c *gin.Context) {
	adminSchoolID := extractAdminSchoolID(c)

	var params PaginationParams
	if err := c.ShouldBindQuery(&params); err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid pagination params")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	parents, total, err := h.parentService.List(ctx, adminSchoolID, params.Limit, params.Offset)
	if err != nil {
		if errors.Is(err, service.ErrSchoolAccessDenied) {
			response.Fail(c, http.StatusForbidden, "access denied")
			return
		}
		response.Fail(c, http.StatusInternalServerError, "failed to fetch parents")
		return
	}

	response.OKPaginated(c, parents, response.Pagination{
		Total:   total,
		Limit:   params.Limit,
		Offset:  params.Offset,
		HasMore: params.Offset+len(parents) < total,
	})
}

// GetByID lấy thông tin phụ huynh theo parent_id (admin only)
func (h *ParentHandler) GetByID(c *gin.Context) {
	adminSchoolID := extractAdminSchoolID(c)

	parentID, err := uuid.Parse(c.Param("parent_id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid parent_id format")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	parent, err := h.parentService.GetByParentID(ctx, adminSchoolID, parentID)
	if err != nil {
		if errors.Is(err, service.ErrSchoolAccessDenied) {
			response.Fail(c, http.StatusForbidden, "access denied")
			return
		}
		if errors.Is(err, pgx.ErrNoRows) {
			response.Fail(c, http.StatusNotFound, "parent not found")
			return
		}
		response.Fail(c, http.StatusInternalServerError, "failed to fetch parent")
		return
	}

	response.OK(c, parent)
}
