package parentscope

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/hoanghonghuy/iris-app/apps/api/internal/auth"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/middleware"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/model"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/service"
)

type fakeParentScopeAppointmentService struct {
	listAvailableSlotsForParentFn func(context.Context, uuid.UUID, uuid.UUID, *time.Time, *time.Time, int, int) ([]model.AppointmentSlot, int, error)
	createAppointmentFn           func(context.Context, uuid.UUID, uuid.UUID, uuid.UUID, string) (model.Appointment, error)
	listParentAppointmentsFn      func(context.Context, uuid.UUID, string, *time.Time, *time.Time, int, int) ([]model.Appointment, int, error)
	cancelAppointmentByParentFn   func(context.Context, uuid.UUID, uuid.UUID, string) (model.Appointment, error)
}

func (f *fakeParentScopeAppointmentService) ListAvailableSlotsForParent(ctx context.Context, parentUserID, studentID uuid.UUID, from, to *time.Time, limit, offset int) ([]model.AppointmentSlot, int, error) {
	if f.listAvailableSlotsForParentFn == nil {
		return nil, 0, errors.New("unexpected ListAvailableSlotsForParent call")
	}
	return f.listAvailableSlotsForParentFn(ctx, parentUserID, studentID, from, to, limit, offset)
}

func (f *fakeParentScopeAppointmentService) CreateAppointment(ctx context.Context, parentUserID, studentID, slotID uuid.UUID, note string) (model.Appointment, error) {
	if f.createAppointmentFn == nil {
		return model.Appointment{}, errors.New("unexpected CreateAppointment call")
	}
	return f.createAppointmentFn(ctx, parentUserID, studentID, slotID, note)
}

func (f *fakeParentScopeAppointmentService) ListParentAppointments(ctx context.Context, parentUserID uuid.UUID, status string, from, to *time.Time, limit, offset int) ([]model.Appointment, int, error) {
	if f.listParentAppointmentsFn == nil {
		return nil, 0, errors.New("unexpected ListParentAppointments call")
	}
	return f.listParentAppointmentsFn(ctx, parentUserID, status, from, to, limit, offset)
}

func (f *fakeParentScopeAppointmentService) CancelAppointmentByParent(ctx context.Context, parentUserID, appointmentID uuid.UUID, cancelReason string) (model.Appointment, error) {
	if f.cancelAppointmentByParentFn == nil {
		return model.Appointment{}, errors.New("unexpected CancelAppointmentByParent call")
	}
	return f.cancelAppointmentByParentFn(ctx, parentUserID, appointmentID, cancelReason)
}

func withParentClaims(userID uuid.UUID) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(middleware.CtxClaims, &auth.Claims{UserID: userID.String(), Roles: []string{"PARENT"}})
		c.Next()
	}
}

func decodeParentScopeError(t *testing.T, rec *httptest.ResponseRecorder) string {
	t.Helper()
	var body map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}
	return body["error"]
}

func TestCancelAppointment_ReturnsBadRequestOnMalformedJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	h := &ParentScopeHandler{}
	r := gin.New()
	r.PATCH("/appointments/:appointment_id/cancel", h.CancelAppointment)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPatch, "/appointments/00000000-0000-0000-0000-000000000001/cancel", strings.NewReader("{invalid-json"))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestCancelAppointment_AllowsEmptyJSONBodyAndContinuesFlow(t *testing.T) {
	gin.SetMode(gin.TestMode)

	h := &ParentScopeHandler{}
	r := gin.New()
	r.PATCH("/appointments/:appointment_id/cancel", h.CancelAppointment)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPatch, "/appointments/00000000-0000-0000-0000-000000000001/cancel", strings.NewReader(""))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, rec.Code)
	}
}

func TestListAvailableAppointmentSlots_InvalidStudentID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	h := &ParentScopeHandler{}
	r := gin.New()
	r.GET("/appointments/slots", h.ListAvailableAppointmentSlots)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/appointments/slots?student_id=bad-uuid", nil)
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestListAvailableAppointmentSlots_InvalidTimeRangeQuery(t *testing.T) {
	gin.SetMode(gin.TestMode)
	parentUserID := uuid.New()
	studentID := uuid.New()

	tests := []struct {
		name string
		url  string
	}{
		{name: "invalid from", url: "/appointments/slots?student_id=" + studentID.String() + "&from=bad-time"},
		{name: "invalid to", url: "/appointments/slots?student_id=" + studentID.String() + "&to=bad-time"},
		{name: "from greater than to", url: "/appointments/slots?student_id=" + studentID.String() + "&from=2026-04-07T12:00:00Z&to=2026-04-07T11:00:00Z"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &ParentScopeHandler{appointmentService: &fakeParentScopeAppointmentService{}}
			r := gin.New()
			r.Use(withParentClaims(parentUserID))
			r.GET("/appointments/slots", h.ListAvailableAppointmentSlots)

			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, tt.url, nil)
			r.ServeHTTP(rec, req)

			if rec.Code != http.StatusBadRequest {
				t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
			}
		})
	}
}

func TestListAvailableAppointmentSlots_TimezoneOffsetRangeForwarding(t *testing.T) {
	gin.SetMode(gin.TestMode)
	parentUserID := uuid.New()
	studentID := uuid.New()
	fromRaw := "2026-04-07T12:00:00+07:00"
	toRaw := "2026-04-07T06:00:00+01:00"
	expectedFrom, _ := time.Parse(time.RFC3339, fromRaw)
	expectedTo, _ := time.Parse(time.RFC3339, toRaw)

	h := &ParentScopeHandler{appointmentService: &fakeParentScopeAppointmentService{
		listAvailableSlotsForParentFn: func(_ context.Context, gotParent, gotStudent uuid.UUID, gotFrom, gotTo *time.Time, _, _ int) ([]model.AppointmentSlot, int, error) {
			if gotParent != parentUserID || gotStudent != studentID {
				t.Fatalf("unexpected ids forwarded")
			}
			if gotFrom == nil || gotTo == nil {
				t.Fatalf("expected non-nil from/to")
			}
			if !gotFrom.Equal(expectedFrom) || !gotTo.Equal(expectedTo) {
				t.Fatalf("timezone forwarding mismatch: from=%v to=%v", *gotFrom, *gotTo)
			}
			return []model.AppointmentSlot{}, 0, nil
		},
	}}

	r := gin.New()
	r.Use(withParentClaims(parentUserID))
	r.GET("/appointments/slots", h.ListAvailableAppointmentSlots)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/appointments/slots?student_id="+studentID.String()+"&from="+url.QueryEscape(fromRaw)+"&to="+url.QueryEscape(toRaw), nil)
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
}

func TestCreateAppointment_InvalidRequestBody(t *testing.T) {
	gin.SetMode(gin.TestMode)

	h := &ParentScopeHandler{}
	r := gin.New()
	r.POST("/appointments", h.CreateAppointment)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/appointments", strings.NewReader("{bad-json"))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestCreateAppointment_InvalidStudentID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	h := &ParentScopeHandler{}
	r := gin.New()
	r.POST("/appointments", h.CreateAppointment)

	body := `{"student_id":"bad-uuid","slot_id":"00000000-0000-0000-0000-000000000001"}`
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/appointments", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestCreateAppointment_InvalidSlotID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	h := &ParentScopeHandler{}
	r := gin.New()
	r.POST("/appointments", h.CreateAppointment)

	body := `{"student_id":"00000000-0000-0000-0000-000000000001","slot_id":"bad-uuid"}`
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/appointments", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestCreateAppointment_SuccessAndErrorMapping(t *testing.T) {
	gin.SetMode(gin.TestMode)
	parentUserID := uuid.New()
	studentID := uuid.New()
	slotID := uuid.New()

	t.Run("success", func(t *testing.T) {
		h := &ParentScopeHandler{appointmentService: &fakeParentScopeAppointmentService{
			createAppointmentFn: func(_ context.Context, gotParent, gotStudent, gotSlot uuid.UUID, _ string) (model.Appointment, error) {
				if gotParent != parentUserID || gotStudent != studentID || gotSlot != slotID {
					t.Fatalf("unexpected ids forwarded")
				}
				return model.Appointment{AppointmentID: uuid.New(), ParentID: gotParent, StudentID: gotStudent, SlotID: gotSlot, Status: "pending"}, nil
			},
		}}

		r := gin.New()
		r.Use(withParentClaims(parentUserID))
		r.POST("/appointments", h.CreateAppointment)

		body := `{"student_id":"` + studentID.String() + `","slot_id":"` + slotID.String() + `","note":"hello"}`
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/appointments", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusCreated {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusCreated)
		}
	})

	tests := []struct {
		name       string
		err        error
		wantStatus int
	}{
		{name: "conflict", err: service.ErrAppointmentSlotUnavailable, wantStatus: http.StatusConflict},
		{name: "forbidden", err: service.ErrForbidden, wantStatus: http.StatusForbidden},
		{name: "bad request", err: service.ErrInvalidValue, wantStatus: http.StatusBadRequest},
		{name: "internal", err: errors.New("boom"), wantStatus: http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &ParentScopeHandler{appointmentService: &fakeParentScopeAppointmentService{
				createAppointmentFn: func(context.Context, uuid.UUID, uuid.UUID, uuid.UUID, string) (model.Appointment, error) {
					return model.Appointment{}, tt.err
				},
			}}

			r := gin.New()
			r.Use(withParentClaims(parentUserID))
			r.POST("/appointments", h.CreateAppointment)

			body := `{"student_id":"` + studentID.String() + `","slot_id":"` + slotID.String() + `"}`
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/appointments", strings.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			r.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
		})
	}
}

func TestListAvailableAppointmentSlots_SuccessAndErrorMapping(t *testing.T) {
	gin.SetMode(gin.TestMode)
	parentUserID := uuid.New()
	studentID := uuid.New()

	t.Run("success", func(t *testing.T) {
		h := &ParentScopeHandler{appointmentService: &fakeParentScopeAppointmentService{
			listAvailableSlotsForParentFn: func(context.Context, uuid.UUID, uuid.UUID, *time.Time, *time.Time, int, int) ([]model.AppointmentSlot, int, error) {
				return []model.AppointmentSlot{{SlotID: uuid.New(), ClassID: uuid.New(), TeacherID: uuid.New()}}, 1, nil
			},
		}}

		r := gin.New()
		r.Use(withParentClaims(parentUserID))
		r.GET("/appointments/slots", h.ListAvailableAppointmentSlots)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/appointments/slots?student_id="+studentID.String()+"&limit=999&offset=-3", nil)
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
		}
	})

	tests := []struct {
		name       string
		err        error
		wantStatus int
	}{
		{name: "forbidden", err: service.ErrForbidden, wantStatus: http.StatusForbidden},
		{name: "bad request", err: service.ErrInvalidValue, wantStatus: http.StatusBadRequest},
		{name: "internal", err: errors.New("boom"), wantStatus: http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &ParentScopeHandler{appointmentService: &fakeParentScopeAppointmentService{
				listAvailableSlotsForParentFn: func(context.Context, uuid.UUID, uuid.UUID, *time.Time, *time.Time, int, int) ([]model.AppointmentSlot, int, error) {
					return nil, 0, tt.err
				},
			}}

			r := gin.New()
			r.Use(withParentClaims(parentUserID))
			r.GET("/appointments/slots", h.ListAvailableAppointmentSlots)

			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/appointments/slots?student_id="+studentID.String(), nil)
			r.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
		})
	}
}

func TestListAvailableAppointmentSlots_PaginationNormalizationForwarding(t *testing.T) {
	gin.SetMode(gin.TestMode)
	parentUserID := uuid.New()
	studentID := uuid.New()

	tests := []struct {
		name       string
		url        string
		wantLimit  int
		wantOffset int
	}{
		{name: "clamp and normalize", url: "/appointments/slots?student_id=" + studentID.String() + "&limit=999&offset=-3", wantLimit: 100, wantOffset: 0},
		{name: "fallback invalid query", url: "/appointments/slots?student_id=" + studentID.String() + "&limit=abc&offset=xyz", wantLimit: 20, wantOffset: 0},
		{name: "default when non-positive", url: "/appointments/slots?student_id=" + studentID.String() + "&limit=0&offset=-1", wantLimit: 20, wantOffset: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &ParentScopeHandler{appointmentService: &fakeParentScopeAppointmentService{
				listAvailableSlotsForParentFn: func(_ context.Context, gotParent, gotStudent uuid.UUID, _ *time.Time, _ *time.Time, gotLimit, gotOffset int) ([]model.AppointmentSlot, int, error) {
					if gotParent != parentUserID || gotStudent != studentID {
						t.Fatalf("unexpected ids forwarded")
					}
					if gotLimit != tt.wantLimit || gotOffset != tt.wantOffset {
						t.Fatalf("forwarded limit/offset = %d/%d, want %d/%d", gotLimit, gotOffset, tt.wantLimit, tt.wantOffset)
					}
					return []model.AppointmentSlot{}, 0, nil
				},
			}}

			r := gin.New()
			r.Use(withParentClaims(parentUserID))
			r.GET("/appointments/slots", h.ListAvailableAppointmentSlots)

			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, tt.url, nil)
			r.ServeHTTP(rec, req)

			if rec.Code != http.StatusOK {
				t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
			}
		})
	}
}

func TestListMyAppointments_SuccessAndErrorMapping(t *testing.T) {
	gin.SetMode(gin.TestMode)
	parentUserID := uuid.New()

	t.Run("success", func(t *testing.T) {
		h := &ParentScopeHandler{appointmentService: &fakeParentScopeAppointmentService{
			listParentAppointmentsFn: func(context.Context, uuid.UUID, string, *time.Time, *time.Time, int, int) ([]model.Appointment, int, error) {
				return []model.Appointment{{AppointmentID: uuid.New(), ParentID: parentUserID, Status: "pending"}}, 1, nil
			},
		}}

		r := gin.New()
		r.Use(withParentClaims(parentUserID))
		r.GET("/appointments", h.ListMyAppointments)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/appointments?limit=500&offset=-1", nil)
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
		}
	})

	tests := []struct {
		name       string
		err        error
		wantStatus int
	}{
		{name: "bad request", err: service.ErrInvalidValue, wantStatus: http.StatusBadRequest},
		{name: "internal", err: errors.New("boom"), wantStatus: http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &ParentScopeHandler{appointmentService: &fakeParentScopeAppointmentService{
				listParentAppointmentsFn: func(context.Context, uuid.UUID, string, *time.Time, *time.Time, int, int) ([]model.Appointment, int, error) {
					return nil, 0, tt.err
				},
			}}

			r := gin.New()
			r.Use(withParentClaims(parentUserID))
			r.GET("/appointments", h.ListMyAppointments)

			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/appointments", nil)
			r.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
		})
	}
}

func TestListMyAppointments_PaginationNormalizationForwarding(t *testing.T) {
	gin.SetMode(gin.TestMode)
	parentUserID := uuid.New()

	tests := []struct {
		name       string
		url        string
		wantLimit  int
		wantOffset int
	}{
		{name: "clamp and normalize", url: "/appointments?limit=500&offset=-1", wantLimit: 100, wantOffset: 0},
		{name: "fallback invalid query", url: "/appointments?limit=abc&offset=xyz", wantLimit: 20, wantOffset: 0},
		{name: "default when non-positive", url: "/appointments?limit=0&offset=-5", wantLimit: 20, wantOffset: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &ParentScopeHandler{appointmentService: &fakeParentScopeAppointmentService{
				listParentAppointmentsFn: func(_ context.Context, gotParent uuid.UUID, _ string, _ *time.Time, _ *time.Time, gotLimit, gotOffset int) ([]model.Appointment, int, error) {
					if gotParent != parentUserID {
						t.Fatalf("unexpected parent forwarded")
					}
					if gotLimit != tt.wantLimit || gotOffset != tt.wantOffset {
						t.Fatalf("forwarded limit/offset = %d/%d, want %d/%d", gotLimit, gotOffset, tt.wantLimit, tt.wantOffset)
					}
					return []model.Appointment{}, 0, nil
				},
			}}

			r := gin.New()
			r.Use(withParentClaims(parentUserID))
			r.GET("/appointments", h.ListMyAppointments)

			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, tt.url, nil)
			r.ServeHTTP(rec, req)

			if rec.Code != http.StatusOK {
				t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
			}
		})
	}
}

func TestListMyAppointments_InvalidTimeRangeQuery(t *testing.T) {
	gin.SetMode(gin.TestMode)
	parentUserID := uuid.New()

	tests := []struct {
		name string
		url  string
	}{
		{name: "invalid from", url: "/appointments?from=bad-time"},
		{name: "invalid to", url: "/appointments?to=bad-time"},
		{name: "from greater than to", url: "/appointments?from=2026-04-07T12:00:00Z&to=2026-04-07T11:00:00Z"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &ParentScopeHandler{appointmentService: &fakeParentScopeAppointmentService{}}
			r := gin.New()
			r.Use(withParentClaims(parentUserID))
			r.GET("/appointments", h.ListMyAppointments)

			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, tt.url, nil)
			r.ServeHTTP(rec, req)

			if rec.Code != http.StatusBadRequest {
				t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
			}
		})
	}
}

func TestListMyAppointments_TimezoneOffsetRangeForwarding(t *testing.T) {
	gin.SetMode(gin.TestMode)
	parentUserID := uuid.New()
	fromRaw := "2026-04-07T12:00:00+07:00"
	toRaw := "2026-04-07T06:00:00+01:00"
	expectedFrom, _ := time.Parse(time.RFC3339, fromRaw)
	expectedTo, _ := time.Parse(time.RFC3339, toRaw)

	h := &ParentScopeHandler{appointmentService: &fakeParentScopeAppointmentService{
		listParentAppointmentsFn: func(_ context.Context, gotParent uuid.UUID, _ string, gotFrom, gotTo *time.Time, _, _ int) ([]model.Appointment, int, error) {
			if gotParent != parentUserID {
				t.Fatalf("unexpected parent forwarded")
			}
			if gotFrom == nil || gotTo == nil {
				t.Fatalf("expected non-nil from/to")
			}
			if !gotFrom.Equal(expectedFrom) || !gotTo.Equal(expectedTo) {
				t.Fatalf("timezone forwarding mismatch: from=%v to=%v", *gotFrom, *gotTo)
			}
			return []model.Appointment{}, 0, nil
		},
	}}

	r := gin.New()
	r.Use(withParentClaims(parentUserID))
	r.GET("/appointments", h.ListMyAppointments)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/appointments?from="+url.QueryEscape(fromRaw)+"&to="+url.QueryEscape(toRaw), nil)
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
}

func TestCancelAppointment_SuccessAndErrorMapping(t *testing.T) {
	gin.SetMode(gin.TestMode)
	parentUserID := uuid.New()
	appointmentID := uuid.New()

	t.Run("success", func(t *testing.T) {
		h := &ParentScopeHandler{appointmentService: &fakeParentScopeAppointmentService{
			cancelAppointmentByParentFn: func(_ context.Context, gotParent, gotAppointment uuid.UUID, _ string) (model.Appointment, error) {
				if gotParent != parentUserID || gotAppointment != appointmentID {
					t.Fatalf("unexpected ids forwarded")
				}
				return model.Appointment{AppointmentID: gotAppointment, ParentID: gotParent, Status: "cancelled"}, nil
			},
		}}

		r := gin.New()
		r.Use(withParentClaims(parentUserID))
		r.PATCH("/appointments/:appointment_id/cancel", h.CancelAppointment)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPatch, "/appointments/"+appointmentID.String()+"/cancel", strings.NewReader(`{"cancel_reason":"parent request"}`))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
		}
	})

	tests := []struct {
		name       string
		err        error
		wantStatus int
	}{
		{name: "forbidden", err: service.ErrForbidden, wantStatus: http.StatusForbidden},
		{name: "bad request", err: service.ErrInvalidValue, wantStatus: http.StatusBadRequest},
		{name: "internal", err: errors.New("boom"), wantStatus: http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &ParentScopeHandler{appointmentService: &fakeParentScopeAppointmentService{
				cancelAppointmentByParentFn: func(context.Context, uuid.UUID, uuid.UUID, string) (model.Appointment, error) {
					return model.Appointment{}, tt.err
				},
			}}

			r := gin.New()
			r.Use(withParentClaims(parentUserID))
			r.PATCH("/appointments/:appointment_id/cancel", h.CancelAppointment)

			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPatch, "/appointments/"+appointmentID.String()+"/cancel", strings.NewReader(`{"cancel_reason":"x"}`))
			req.Header.Set("Content-Type", "application/json")
			r.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
		})
	}
}

func TestCancelAppointment_RejectsInvalidAppointmentID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := &ParentScopeHandler{}
	r := gin.New()
	r.PATCH("/appointments/:appointment_id/cancel", h.CancelAppointment)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPatch, "/appointments/not-uuid/cancel", strings.NewReader(`{"cancel_reason":"x"}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
	if got := decodeParentScopeError(t, rec); got != "invalid appointment_id" {
		t.Fatalf("error = %q, want %q", got, "invalid appointment_id")
	}
}
