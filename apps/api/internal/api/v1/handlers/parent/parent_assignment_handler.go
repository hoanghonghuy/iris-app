package parenthandlers

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/hoanghonghuy/iris-app/apps/api/internal/api/v1/handlers/shared"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/response"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/service"
)

// AssignStudent gán phụ huynh cho học sinh (admin only)
func (h *ParentHandler) AssignStudent(c *gin.Context) {
	adminSchoolID := shared.ExtractAdminSchoolID(c)

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
	adminSchoolID := shared.ExtractAdminSchoolID(c)

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
