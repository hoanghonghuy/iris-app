package httpapi

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/repo"
)

type AdminClassesHandler struct {
	Classes *repo.ClassRepo
}

type CreateClassReq struct {
	Name       string `json:"name" binding:"required,min=1,max=100"`
	SchoolYear string `json:"school_year" binding:"required,min=4,max=20"`
}

func (h *AdminClassesHandler) Create(c *gin.Context) {
	schoolID, err := uuid.Parse(c.Param("school_id"))
	if err != nil {
		fail(c, http.StatusBadRequest, "invalid school_id")
		return
	}

	var req CreateClassReq
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, http.StatusBadRequest, "invalid payload")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	id, err := h.Classes.Create(ctx, schoolID, req.Name, req.SchoolYear)
	if err != nil {
		// FK fail -> school_id not found
		fail(c, http.StatusInternalServerError, "db error")
		return
	}

	created(c, gin.H{
		"class_id": id,
	})
}

func (h *AdminClassesHandler) ListBySchool(c *gin.Context) {
	schoolId, err := uuid.Parse(c.Param("school_id"))
	if err != nil {
		fail(c, http.StatusBadRequest, "invalid school_id")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	classes, err := h.Classes.List(ctx, schoolId)
	if err != nil {
		fail(c, http.StatusInternalServerError, "db error")
		return
	}

	ok(c, classes)
}
