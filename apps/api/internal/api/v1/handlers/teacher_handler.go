package handlers

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/hoanghonghuy/iris-app/apps/api/internal/response"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/service"
)

type TeacherHandler struct {
	teacherService *service.TeacherService
}

// UpdateTeacherRequest input để admin cập nhật thông tin giáo viên
type UpdateTeacherRequest struct {
	FullName string    `json:"full_name"`
	Phone    string    `json:"phone"`
	SchoolID uuid.UUID `json:"school_id"`
}

func NewTeacherHandler(teacherService *service.TeacherService) *TeacherHandler {
	return &TeacherHandler{
		teacherService: teacherService,
	}
}

func (h *TeacherHandler) List(c *gin.Context) {
	var params PaginationParams
	if err := c.ShouldBindQuery(&params); err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid pagination params")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	teachers, total, err := h.teacherService.List(ctx, params.Limit, params.Offset)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "failed to fetch teachers")
		return
	}

	response.OKPaginated(c, teachers, response.Pagination{
		Total:   total,
		Limit:   params.Limit,
		Offset:  params.Offset,
		HasMore: params.Offset+len(teachers) < total,
	})
}

func (h *TeacherHandler) ListTeacherOfClass(c *gin.Context) {
	classID, err := uuid.Parse(c.Param("class_id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid class_id format")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	teachers, err := h.teacherService.ListTeachersOfClass(ctx, classID)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "failed to fetch teachers of class")
		return
	}
	response.OK(c, teachers)
}

func (h *TeacherHandler) GetByTeacherID(c *gin.Context) {
	teacherID, err := uuid.Parse(c.Param("teacher_id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid teacher_id format")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	teacher, err := h.teacherService.GetByTeacherID(ctx, teacherID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			response.Fail(c, http.StatusNotFound, "teacher not found")
			return
		}
		response.Fail(c, http.StatusInternalServerError, "failed to fetch teacher")
		return
	}
	response.OK(c, teacher)
}

func (h *TeacherHandler) Assign(c *gin.Context) {
	teacherID, err := uuid.Parse(c.Param("teacher_id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid teacher_id format")
		return
	}
	classID, err := uuid.Parse(c.Param("class_id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid class_id format")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	err = h.teacherService.Assign(ctx, teacherID, classID)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "failed to assign teacher to class")
		return
	}
	response.OK(c, gin.H{
		"message":    "teacher assigned to class successfully",
		"teacher_id": teacherID,
		"class_id":   classID,
	})
}

func (h *TeacherHandler) Unassign(c *gin.Context) {
	teacherID, err := uuid.Parse(c.Param("teacher_id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid teacher_id format")
		return
	}
	classID, err := uuid.Parse(c.Param("class_id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid class_id format")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	err = h.teacherService.Unassign(ctx, teacherID, classID)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "failed to unassign teacher from class")
		return
	}
	response.OK(c, gin.H{
		"message":    "teacher unassigned from class successfully",
		"teacher_id": teacherID,
		"class_id":   classID,
	})
}

// Update updates a teacher's information (admin only - can update all fields)
func (h *TeacherHandler) Update(c *gin.Context) {
	teacherID, err := uuid.Parse(c.Param("teacher_id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid teacher_id format")
		return
	}

	var req UpdateTeacherRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid request body")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	err = h.teacherService.Update(ctx, teacherID, req.FullName, req.Phone, req.SchoolID)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "failed to update teacher")
		return
	}

	response.OK(c, gin.H{
		"message":    "teacher updated successfully",
		"teacher_id": teacherID.String(),
	})
}
