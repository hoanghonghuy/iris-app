package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/hoanghonghuy/iris-app/apps/api/internal/auth"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/middleware"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/response"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/service"
)

type CreateAppointmentSlotRequest struct {
	ClassID        string `json:"class_id" binding:"required"`
	StartTime      string `json:"start_time" binding:"required"`
	EndTime        string `json:"end_time"`
	DurationMinute int    `json:"duration_minutes"`
	Note           string `json:"note"`
}

type UpdateAppointmentStatusRequest struct {
	Status       string `json:"status" binding:"required"`
	CancelReason string `json:"cancel_reason"`
}

func (h *TeacherScopeHandler) CreateAppointmentSlot(c *gin.Context) {
	var req CreateAppointmentSlotRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid request body")
		return
	}

	classID, err := uuid.Parse(req.ClassID)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid class_id")
		return
	}

	startTime, err := time.Parse(time.RFC3339, req.StartTime)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid start_time (RFC3339)")
		return
	}

	endTime := time.Time{}
	if req.EndTime != "" {
		endTime, err = time.Parse(time.RFC3339, req.EndTime)
		if err != nil {
			response.Fail(c, http.StatusBadRequest, "invalid end_time (RFC3339)")
			return
		}
	} else {
		d := req.DurationMinute
		if d <= 0 {
			d = 30
		}
		endTime = startTime.Add(time.Duration(d) * time.Minute)
	}

	claimsAny, exists := c.Get(middleware.CtxClaims)
	if !exists {
		response.Fail(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	claims := claimsAny.(*auth.Claims)
	teacherUserID, err := uuid.Parse(claims.UserID)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid user ID")
		return
	}

	slot, err := h.appointmentService.CreateSlot(c.Request.Context(), teacherUserID, classID, startTime, endTime, req.Note)
	if err != nil {
		if errors.Is(err, service.ErrForbidden) {
			response.Fail(c, http.StatusForbidden, "forbidden")
			return
		}
		if errors.Is(err, service.ErrInvalidValue) || errors.Is(err, service.ErrInvalidClassID) {
			response.Fail(c, http.StatusBadRequest, err.Error())
			return
		}
		response.Fail(c, http.StatusInternalServerError, "failed to create appointment slot")
		return
	}

	response.Created(c, slot)
}

func (h *TeacherScopeHandler) ListMyAppointments(c *gin.Context) {
	claimsAny, exists := c.Get(middleware.CtxClaims)
	if !exists {
		response.Fail(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	claims := claimsAny.(*auth.Claims)
	teacherUserID, err := uuid.Parse(claims.UserID)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid user ID")
		return
	}

	status := c.Query("status")
	from, to, err := h.parseTimeRange(c.Query("from"), c.Query("to"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	limit, offset := parsePagination(c.Query("limit"), c.Query("offset"))

	items, total, err := h.appointmentService.ListTeacherAppointments(c.Request.Context(), teacherUserID, status, from, to, limit, offset)
	if err != nil {
		if errors.Is(err, service.ErrInvalidValue) {
			response.Fail(c, http.StatusBadRequest, err.Error())
			return
		}
		response.Fail(c, http.StatusInternalServerError, "failed to list appointments")
		return
	}

	response.OKPaginated(c, items, response.Pagination{
		Total:   total,
		Limit:   limit,
		Offset:  offset,
		HasMore: offset+len(items) < total,
	})
}

func (h *TeacherScopeHandler) UpdateAppointmentStatus(c *gin.Context) {
	appointmentID, err := uuid.Parse(c.Param("appointment_id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid appointment_id")
		return
	}

	var req UpdateAppointmentStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid request body")
		return
	}

	claimsAny, exists := c.Get(middleware.CtxClaims)
	if !exists {
		response.Fail(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	claims := claimsAny.(*auth.Claims)
	teacherUserID, err := uuid.Parse(claims.UserID)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid user ID")
		return
	}

	updated, err := h.appointmentService.UpdateAppointmentStatusByTeacher(c.Request.Context(), teacherUserID, appointmentID, req.Status, req.CancelReason)
	if err != nil {
		if errors.Is(err, service.ErrForbidden) {
			response.Fail(c, http.StatusForbidden, "forbidden")
			return
		}
		if errors.Is(err, service.ErrInvalidValue) {
			response.Fail(c, http.StatusBadRequest, err.Error())
			return
		}
		response.Fail(c, http.StatusInternalServerError, "failed to update appointment status")
		return
	}

	response.OK(c, updated)
}

func (h *TeacherScopeHandler) parseTimeRange(fromRaw, toRaw string) (*time.Time, *time.Time, error) {
	var from *time.Time
	var to *time.Time
	if fromRaw != "" {
		v, err := time.Parse(time.RFC3339, fromRaw)
		if err != nil {
			return nil, nil, err
		}
		from = &v
	}
	if toRaw != "" {
		v, err := time.Parse(time.RFC3339, toRaw)
		if err != nil {
			return nil, nil, err
		}
		to = &v
	}
	return from, to, nil
}

func parsePagination(limitRaw, offsetRaw string) (int, int) {
	limit := 20
	if limitRaw != "" {
		if n, err := strconv.Atoi(limitRaw); err == nil {
			limit = n
		}
	}
	offset := 0
	if offsetRaw != "" {
		if n, err := strconv.Atoi(offsetRaw); err == nil {
			offset = n
		}
	}
	return limit, offset
}
