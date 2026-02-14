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

func (s *StudentHandler) Create(c *gin.Context) {
	var req CreateStudentReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid request body")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	student, err := s.studentService.Create(ctx, req.SchoolID, req.CurrentClassID, req.FullName, req.DOB, req.Gender)
	if err != nil {
		if errors.Is(err, service.ErrInvalidValue) {
			response.Fail(c, http.StatusBadRequest, "invalid dob format (expected YYYY-MM-DD)")
			return
		}
		response.Fail(c, http.StatusBadRequest, "failed to create student")
		return
	}

	response.Created(c, gin.H{"student_id": student.ID})
	// TODO: Set Location header và trả về created resource
	// c.Header("Location", fmt.Sprintf("/api/v1/admin/students/%s", student.ID.String()))
	// response.Created(c, student)

}

func (s *StudentHandler) ListByClass(c *gin.Context) {
	classID, err := uuid.Parse(c.Param("class_id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid class_id format")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	students, err := s.studentService.ListByClass(ctx, classID)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "failed to fetch students")
		return
	}

	response.OK(c, students)
}
