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
)

type StudentHandler struct {
	studentService *service.StudentService
}

func NewStudentHandler(studentService *service.StudentService) *StudentHandler {
	return &StudentHandler{
		studentService: studentService,
	}
}

type CreateStudentReq struct {
	SchoolID       uuid.UUID `json:"school_id" binding:"required"` // TODO: gửi từ UI hoặc lấy từ class
	CurrentClassID uuid.UUID `json:"current_class_id" binding:"required"`
	FullName       string    `json:"full_name" binding:"required,min=1,max=120"`
	DOB            string    `json:"dob" binding:"required"`    // YYYY-MM-DD
	Gender         string    `json:"gender" binding:"required"` // male/female/other
}

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
	// TODO: Set Location header và trả về created resource
	// c.Header("Location", fmt.Sprintf("/api/v1/admin/students/%s", student.StudentID.String()))
	// response.Created(c, student)

}

func (h *StudentHandler) ListByClass(c *gin.Context) {
	adminSchoolID := extractAdminSchoolID(c)

	classID, err := uuid.Parse(c.Param("class_id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid class_id format")
		return
	}

	var params PaginationParams
	if err := c.ShouldBindQuery(&params); err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid pagination params")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	students, total, err := h.studentService.ListByClass(ctx, adminSchoolID, classID, params.Limit, params.Offset)
	if err != nil {
		if errors.Is(err, service.ErrSchoolAccessDenied) {
			response.Fail(c, http.StatusForbidden, "access denied")
			return
		}
		response.Fail(c, http.StatusInternalServerError, "failed to fetch students")
		return
	}

	response.OKPaginated(c, students, response.Pagination{
		Total:   total,
		Limit:   params.Limit,
		Offset:  params.Offset,
		HasMore: params.Offset+len(students) < total,
	})
}

// GetProfile lấy thông tin chi tiết của một học sinh
func (h *StudentHandler) GetProfile(c *gin.Context) {
	adminSchoolID := extractAdminSchoolID(c)

	studentID, err := uuid.Parse(c.Param("student_id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid student_id format")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	profile, err := h.studentService.GetProfile(ctx, adminSchoolID, studentID)
	if err != nil {
		if errors.Is(err, service.ErrSchoolAccessDenied) {
			response.Fail(c, http.StatusForbidden, "access denied")
			return
		}
		if errors.Is(err, service.ErrFailedToGetStudent) {
			response.Fail(c, http.StatusNotFound, "student not found")
			return
		}
		response.Fail(c, http.StatusInternalServerError, "failed to fetch student profile")
		return
	}

	response.OK(c, profile)
}
