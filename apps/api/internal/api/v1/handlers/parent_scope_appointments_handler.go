package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/hoanghonghuy/iris-app/apps/api/internal/response"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/service"
)

type CreateAppointmentRequest struct {
	StudentID string `json:"student_id" binding:"required"`
	SlotID    string `json:"slot_id" binding:"required"`
	Note      string `json:"note"`
}

type CancelAppointmentRequest struct {
	CancelReason string `json:"cancel_reason"`
}

func (h *ParentScopeHandler) ListAvailableAppointmentSlots(c *gin.Context) {
	studentID, err := uuid.Parse(c.Query("student_id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid student_id")
		return
	}

	parentUserID, ok := requireCurrentUserID(c)
	if !ok {
		return
	}

	from, to, err := parseTimeRange(c.Query("from"), c.Query("to"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	limit, offset := parsePagination(c.Query("limit"), c.Query("offset"))

	items, total, err := h.appointmentService.ListAvailableSlotsForParent(c.Request.Context(), parentUserID, studentID, from, to, limit, offset)
	if err != nil {
		if errors.Is(err, service.ErrForbidden) || errors.Is(err, service.ErrInvalidValue) {
			response.Fail(c, http.StatusBadRequest, err.Error())
			return
		}
		response.Fail(c, http.StatusInternalServerError, "failed to list available slots")
		return
	}

	response.OKPaginated(c, items, response.Pagination{
		Total:   total,
		Limit:   limit,
		Offset:  offset,
		HasMore: offset+len(items) < total,
	})
}

func (h *ParentScopeHandler) CreateAppointment(c *gin.Context) {
	var req CreateAppointmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid request body")
		return
	}

	studentID, err := uuid.Parse(req.StudentID)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid student_id")
		return
	}
	slotID, err := uuid.Parse(req.SlotID)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid slot_id")
		return
	}

	parentUserID, ok := requireCurrentUserID(c)
	if !ok {
		return
	}

	appointment, err := h.appointmentService.CreateAppointment(c.Request.Context(), parentUserID, studentID, slotID, req.Note)
	if err != nil {
		if errors.Is(err, service.ErrForbidden) || errors.Is(err, service.ErrInvalidValue) {
			response.Fail(c, http.StatusBadRequest, err.Error())
			return
		}
		response.Fail(c, http.StatusInternalServerError, "failed to create appointment")
		return
	}

	response.Created(c, appointment)
}

func (h *ParentScopeHandler) ListMyAppointments(c *gin.Context) {
	parentUserID, ok := requireCurrentUserID(c)
	if !ok {
		return
	}

	status := c.Query("status")
	from, to, err := parseTimeRange(c.Query("from"), c.Query("to"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	limit, offset := parsePagination(c.Query("limit"), c.Query("offset"))

	items, total, err := h.appointmentService.ListParentAppointments(c.Request.Context(), parentUserID, status, from, to, limit, offset)
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

func (h *ParentScopeHandler) CancelAppointment(c *gin.Context) {
	appointmentID, err := uuid.Parse(c.Param("appointment_id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid appointment_id")
		return
	}

	var req CancelAppointmentRequest
	_ = c.ShouldBindJSON(&req)

	parentUserID, ok := requireCurrentUserID(c)
	if !ok {
		return
	}

	updated, err := h.appointmentService.CancelAppointmentByParent(c.Request.Context(), parentUserID, appointmentID, req.CancelReason)
	if err != nil {
		if errors.Is(err, service.ErrForbidden) || errors.Is(err, service.ErrInvalidValue) {
			response.Fail(c, http.StatusBadRequest, err.Error())
			return
		}
		response.Fail(c, http.StatusInternalServerError, "failed to cancel appointment")
		return
	}

	response.OK(c, updated)
}

func (h *ParentScopeHandler) GetMyAnalytics(c *gin.Context) {
	parentUserID, ok := requireCurrentUserID(c)
	if !ok {
		return
	}

	stats, err := h.parentScopeService.GetMyAnalytics(c.Request.Context(), parentUserID)
	if err != nil {
		if errors.Is(err, service.ErrInvalidUserID) {
			response.Fail(c, http.StatusBadRequest, err.Error())
			return
		}
		response.Fail(c, http.StatusInternalServerError, "failed to fetch parent analytics")
		return
	}
	response.OK(c, stats)
}
