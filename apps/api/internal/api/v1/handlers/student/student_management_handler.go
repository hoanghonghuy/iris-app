package studenthandlers

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

func (h *StudentHandler) Create(c *gin.Context) {
	adminSchoolID := extractAdminSchoolID(c)

	var req CreateStudentReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid request body")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	student, err := h.studentService.Create(ctx, adminSchoolID, req.SchoolID, req.CurrentClassID, req.FullName, req.DOB, req.Gender)
	if err != nil {
		if errors.Is(err, service.ErrSchoolAccessDenied) {
			response.Fail(c, http.StatusForbidden, "access denied")
			return
		}
		if errors.Is(err, service.ErrInvalidValue) {
			response.Fail(c, http.StatusBadRequest, "invalid dob format (expected YYYY-MM-DD)")
			return
		}
		response.Fail(c, http.StatusBadRequest, "failed to create student")
		return
	}

	response.Created(c, gin.H{"student_id": student.StudentID})
}

// Update cập nhật thông tin học sinh
func (h *StudentHandler) Update(c *gin.Context) {
	adminSchoolID := extractAdminSchoolID(c)

	studentID, err := uuid.Parse(c.Param("student_id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid student_id format")
		return
	}

	var req UpdateStudentReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid request body")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	if err := h.studentService.Update(ctx, adminSchoolID, studentID, req.FullName, req.DOB, req.Gender); err != nil {
		if errors.Is(err, service.ErrSchoolAccessDenied) {
			response.Fail(c, http.StatusForbidden, "access denied")
			return
		}
		if errors.Is(err, service.ErrStudentNotFound) {
			response.Fail(c, http.StatusNotFound, "student not found")
			return
		}
		if errors.Is(err, service.ErrInvalidValue) {
			response.Fail(c, http.StatusBadRequest, "invalid dob format (expected YYYY-MM-DD)")
			return
		}
		response.Fail(c, http.StatusInternalServerError, "failed to update student")
		return
	}

	response.OK(c, gin.H{"message": "student updated successfully", "student_id": studentID.String()})
}

// Delete xóa học sinh
func (h *StudentHandler) Delete(c *gin.Context) {
	adminSchoolID := extractAdminSchoolID(c)

	studentID, err := uuid.Parse(c.Param("student_id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid student_id format")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	if err := h.studentService.Delete(ctx, adminSchoolID, studentID); err != nil {
		if errors.Is(err, service.ErrSchoolAccessDenied) {
			response.Fail(c, http.StatusForbidden, "access denied")
			return
		}
		if errors.Is(err, service.ErrStudentNotFound) {
			response.Fail(c, http.StatusNotFound, "student not found")
			return
		}
		response.Fail(c, http.StatusInternalServerError, "failed to delete student")
		return
	}

	response.OK(c, gin.H{"message": "student deleted successfully", "student_id": studentID.String()})
}
