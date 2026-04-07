package parentcodehandlers

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

// GenerateCodeForStudent tạo parent code cho student (admin only)
func (h *ParentCodeHandler) GenerateCodeForStudent(c *gin.Context) {
	adminSchoolID := shared.ExtractAdminSchoolID(c)

	studentID, err := uuid.Parse(c.Param("student_id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid student ID format")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	code, err := h.parentCodeService.GenerateCodeForStudent(ctx, adminSchoolID, studentID)
	if err != nil {
		if errors.Is(err, service.ErrSchoolAccessDenied) {
			response.Fail(c, http.StatusForbidden, "access denied")
			return
		}
		response.Fail(c, http.StatusInternalServerError, "failed to generate parent code")
		return
	}

	response.OK(c, gin.H{
		"student_id":  studentID,
		"parent_code": code,
		"message":     "share this code with parent to allow registration",
		"max_usage":   4,
		"expires_at":  time.Now().AddDate(0, 0, 7),
	})
}

// RevokeParentCode thu hồi parent code cua hs
func (h *ParentCodeHandler) RevokeParentCode(c *gin.Context) {
	adminSchoolID := shared.ExtractAdminSchoolID(c)

	studentID, err := uuid.Parse(c.Param("student_id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid student ID format")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	if err := h.parentCodeService.RevokeCode(ctx, adminSchoolID, studentID); err != nil {
		if errors.Is(err, service.ErrSchoolAccessDenied) {
			response.Fail(c, http.StatusForbidden, "access denied")
			return
		}
		response.Fail(c, http.StatusInternalServerError, "failed to revoke parent code")
		return
	}

	response.OK(c, gin.H{"message": "parent code revoked successfully"})
}
