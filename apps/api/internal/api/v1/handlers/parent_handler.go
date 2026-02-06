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

	err = h.parentService.AssignStudent(ctx, parentID, studentID, req.Relationship)
	if err != nil {
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

	err = h.parentService.UnassignStudent(ctx, parentID, studentID)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "failed to unassign parent from student")
		return
	}

	response.OK(c, gin.H{
		"message":    "parent unassigned from student successfully",
		"parent_id":  parentID.String(),
		"student_id": studentID.String(),
	})
}

// List lấy danh sách tất cả phụ huynh (admin only)
func (h *ParentHandler) List(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	parents, err := h.parentService.List(ctx)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "failed to fetch parents")
		return
	}

	response.OK(c, parents)
}

// GetByID lấy thông tin phụ huynh theo parent_id (admin only)
func (h *ParentHandler) GetByID(c *gin.Context) {
	parentID, err := uuid.Parse(c.Param("parent_id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid parent_id format")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	parent, err := h.parentService.GetByParentID(ctx, parentID)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "failed to fetch parent")
		return
	}

	response.OK(c, parent)
}
