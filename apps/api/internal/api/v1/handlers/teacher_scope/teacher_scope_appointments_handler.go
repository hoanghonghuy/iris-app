package teacherscope

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/hoanghonghuy/iris-app/apps/api/internal/api/v1/handlers/shared"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/response"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/service"
)

type CreateAppointmentSlotRequest struct {
	ClassID           string `json:"class_id" binding:"required"`
	StartTime         string `json:"start_time" binding:"required"`
	EndTime           string `json:"end_time"`
	DurationMinute    int    `json:"duration_minutes"`
	BufferMinutes     int    `json:"buffer_minutes"`
	MaxBookingsPerDay int    `json:"max_bookings_per_day"`
	Note              string `json:"note"`
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

	var endTime time.Time
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

	teacherUserID, ok := shared.RequireCurrentUserID(c)
	if !ok {
		return
	}

	slot, err := h.appointmentService.CreateSlot(
		c.Request.Context(),
		teacherUserID,
		classID,
		startTime,
		endTime,
		req.Note,
		req.BufferMinutes,
		req.MaxBookingsPerDay,
	)
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
	teacherUserID, ok := shared.RequireCurrentUserID(c)
	if !ok {
		return
	}

	status := c.Query("status")
	from, to, err := shared.ParseTimeRange(c.Query("from"), c.Query("to"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	limit, offset := shared.ParsePagination(c.Query("limit"), c.Query("offset"))
	limit, offset = shared.NormalizePagination(limit, offset)

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

	teacherUserID, ok := shared.RequireCurrentUserID(c)
	if !ok {
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
