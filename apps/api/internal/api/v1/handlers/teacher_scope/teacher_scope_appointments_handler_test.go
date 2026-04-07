package teacherscope

import (
	"context"
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

type fakeTeacherScopeAppointmentService struct {
	createSlotFn                     func(context.Context, uuid.UUID, uuid.UUID, time.Time, time.Time, string) (model.AppointmentSlot, error)
	listTeacherAppointmentsFn        func(context.Context, uuid.UUID, string, *time.Time, *time.Time, int, int) ([]model.Appointment, int, error)
	updateAppointmentStatusByTeacher func(context.Context, uuid.UUID, uuid.UUID, string, string) (model.Appointment, error)
}

func (f *fakeTeacherScopeAppointmentService) CreateSlot(ctx context.Context, teacherUserID, classID uuid.UUID, startTime, endTime time.Time, note string) (model.AppointmentSlot, error) {
	if f.createSlotFn == nil {
		return model.AppointmentSlot{}, errors.New("unexpected CreateSlot call")
	}
	return f.createSlotFn(ctx, teacherUserID, classID, startTime, endTime, note)
}

func (f *fakeTeacherScopeAppointmentService) ListTeacherAppointments(ctx context.Context, teacherUserID uuid.UUID, status string, from, to *time.Time, limit, offset int) ([]model.Appointment, int, error) {
	if f.listTeacherAppointmentsFn == nil {
		return nil, 0, errors.New("unexpected ListTeacherAppointments call")
	}
	return f.listTeacherAppointmentsFn(ctx, teacherUserID, status, from, to, limit, offset)
}

func (f *fakeTeacherScopeAppointmentService) UpdateAppointmentStatusByTeacher(ctx context.Context, teacherUserID, appointmentID uuid.UUID, status, cancelReason string) (model.Appointment, error) {
	if f.updateAppointmentStatusByTeacher == nil {
		return model.Appointment{}, errors.New("unexpected UpdateAppointmentStatusByTeacher call")
	}
	return f.updateAppointmentStatusByTeacher(ctx, teacherUserID, appointmentID, status, cancelReason)
}

func withTeacherClaims(userID uuid.UUID) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(middleware.CtxClaims, &auth.Claims{UserID: userID.String(), Roles: []string{"TEACHER"}})
		c.Next()
	}
}

func TestCreateAppointmentSlot_InvalidClassID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	h := &TeacherScopeHandler{}
	r := gin.New()
	r.POST("/appointments/slots", h.CreateAppointmentSlot)

	body := `{"class_id":"bad-uuid","start_time":"2026-04-10T10:00:00Z"}`
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/appointments/slots", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestCreateAppointmentSlot_InvalidStartTime(t *testing.T) {
	gin.SetMode(gin.TestMode)

	h := &TeacherScopeHandler{}
	r := gin.New()
	r.POST("/appointments/slots", h.CreateAppointmentSlot)

	body := `{"class_id":"00000000-0000-0000-0000-000000000001","start_time":"not-rfc3339"}`
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/appointments/slots", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestCreateAppointmentSlot_InvalidEndTime(t *testing.T) {
	gin.SetMode(gin.TestMode)

	h := &TeacherScopeHandler{}
	r := gin.New()
	r.POST("/appointments/slots", h.CreateAppointmentSlot)

	body := `{"class_id":"00000000-0000-0000-0000-000000000001","start_time":"2026-04-10T10:00:00Z","end_time":"invalid-end"}`
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/appointments/slots", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestUpdateAppointmentStatus_InvalidAppointmentID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	h := &TeacherScopeHandler{}
	r := gin.New()
	r.PATCH("/appointments/:appointment_id/status", h.UpdateAppointmentStatus)

	body := `{"status":"confirmed"}`
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPatch, "/appointments/not-a-uuid/status", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestUpdateAppointmentStatus_InvalidRequestBody(t *testing.T) {
	gin.SetMode(gin.TestMode)

	h := &TeacherScopeHandler{}
	r := gin.New()
	r.PATCH("/appointments/:appointment_id/status", h.UpdateAppointmentStatus)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPatch, "/appointments/00000000-0000-0000-0000-000000000001/status", strings.NewReader("{bad-json"))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestCreateAppointmentSlot_SuccessAndErrorMapping(t *testing.T) {
	gin.SetMode(gin.TestMode)
	teacherUserID := uuid.New()
	classID := uuid.New()

	t.Run("success", func(t *testing.T) {
		h := &TeacherScopeHandler{appointmentService: &fakeTeacherScopeAppointmentService{
			createSlotFn: func(_ context.Context, gotTeacher, gotClass uuid.UUID, _ time.Time, _ time.Time, _ string) (model.AppointmentSlot, error) {
				if gotTeacher != teacherUserID || gotClass != classID {
					t.Fatalf("unexpected ids forwarded")
				}
				return model.AppointmentSlot{SlotID: uuid.New(), TeacherID: gotTeacher, ClassID: gotClass}, nil
			},
		}}

		r := gin.New()
		r.Use(withTeacherClaims(teacherUserID))
		r.POST("/appointments/slots", h.CreateAppointmentSlot)

		body := `{"class_id":"` + classID.String() + `","start_time":"2030-01-01T10:00:00Z","duration_minutes":30}`
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/appointments/slots", strings.NewReader(body))
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
		{name: "forbidden", err: service.ErrForbidden, wantStatus: http.StatusForbidden},
		{name: "invalid value", err: service.ErrInvalidValue, wantStatus: http.StatusBadRequest},
		{name: "invalid class", err: service.ErrInvalidClassID, wantStatus: http.StatusBadRequest},
		{name: "internal", err: errors.New("boom"), wantStatus: http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &TeacherScopeHandler{appointmentService: &fakeTeacherScopeAppointmentService{
				createSlotFn: func(context.Context, uuid.UUID, uuid.UUID, time.Time, time.Time, string) (model.AppointmentSlot, error) {
					return model.AppointmentSlot{}, tt.err
				},
			}}

			r := gin.New()
			r.Use(withTeacherClaims(teacherUserID))
			r.POST("/appointments/slots", h.CreateAppointmentSlot)

			body := `{"class_id":"` + classID.String() + `","start_time":"2030-01-01T10:00:00Z"}`
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/appointments/slots", strings.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			r.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
		})
	}
}

func TestListMyAppointments_SuccessAndErrorMapping(t *testing.T) {
	gin.SetMode(gin.TestMode)
	teacherUserID := uuid.New()

	t.Run("success", func(t *testing.T) {
		h := &TeacherScopeHandler{appointmentService: &fakeTeacherScopeAppointmentService{
			listTeacherAppointmentsFn: func(context.Context, uuid.UUID, string, *time.Time, *time.Time, int, int) ([]model.Appointment, int, error) {
				return []model.Appointment{{AppointmentID: uuid.New(), TeacherID: uuid.New(), Status: "pending"}}, 1, nil
			},
		}}

		r := gin.New()
		r.Use(withTeacherClaims(teacherUserID))
		r.GET("/appointments", h.ListMyAppointments)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/appointments?limit=500&offset=-3", nil)
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
			h := &TeacherScopeHandler{appointmentService: &fakeTeacherScopeAppointmentService{
				listTeacherAppointmentsFn: func(context.Context, uuid.UUID, string, *time.Time, *time.Time, int, int) ([]model.Appointment, int, error) {
					return nil, 0, tt.err
				},
			}}

			r := gin.New()
			r.Use(withTeacherClaims(teacherUserID))
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
	teacherUserID := uuid.New()

	tests := []struct {
		name       string
		url        string
		wantLimit  int
		wantOffset int
	}{
		{name: "clamp and normalize", url: "/appointments?limit=500&offset=-3", wantLimit: 100, wantOffset: 0},
		{name: "fallback invalid query", url: "/appointments?limit=abc&offset=xyz", wantLimit: 20, wantOffset: 0},
		{name: "default when non-positive", url: "/appointments?limit=0&offset=-5", wantLimit: 20, wantOffset: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &TeacherScopeHandler{appointmentService: &fakeTeacherScopeAppointmentService{
				listTeacherAppointmentsFn: func(_ context.Context, gotTeacher uuid.UUID, _ string, _ *time.Time, _ *time.Time, gotLimit, gotOffset int) ([]model.Appointment, int, error) {
					if gotTeacher != teacherUserID {
						t.Fatalf("unexpected teacher forwarded")
					}
					if gotLimit != tt.wantLimit || gotOffset != tt.wantOffset {
						t.Fatalf("forwarded limit/offset = %d/%d, want %d/%d", gotLimit, gotOffset, tt.wantLimit, tt.wantOffset)
					}
					return []model.Appointment{}, 0, nil
				},
			}}

			r := gin.New()
			r.Use(withTeacherClaims(teacherUserID))
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
	teacherUserID := uuid.New()

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
			h := &TeacherScopeHandler{appointmentService: &fakeTeacherScopeAppointmentService{}}
			r := gin.New()
			r.Use(withTeacherClaims(teacherUserID))
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
	teacherUserID := uuid.New()
	fromRaw := "2026-04-07T12:00:00+07:00"
	toRaw := "2026-04-07T06:00:00+01:00"
	expectedFrom, _ := time.Parse(time.RFC3339, fromRaw)
	expectedTo, _ := time.Parse(time.RFC3339, toRaw)

	h := &TeacherScopeHandler{appointmentService: &fakeTeacherScopeAppointmentService{
		listTeacherAppointmentsFn: func(_ context.Context, gotTeacher uuid.UUID, _ string, gotFrom, gotTo *time.Time, _, _ int) ([]model.Appointment, int, error) {
			if gotTeacher != teacherUserID {
				t.Fatalf("unexpected teacher forwarded")
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
	r.Use(withTeacherClaims(teacherUserID))
	r.GET("/appointments", h.ListMyAppointments)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/appointments?from="+url.QueryEscape(fromRaw)+"&to="+url.QueryEscape(toRaw), nil)
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
}

func TestUpdateAppointmentStatus_SuccessAndErrorMapping(t *testing.T) {
	gin.SetMode(gin.TestMode)
	teacherUserID := uuid.New()
	appointmentID := uuid.New()

	t.Run("success", func(t *testing.T) {
		h := &TeacherScopeHandler{appointmentService: &fakeTeacherScopeAppointmentService{
			updateAppointmentStatusByTeacher: func(_ context.Context, gotTeacher, gotAppointment uuid.UUID, gotStatus, _ string) (model.Appointment, error) {
				if gotTeacher != teacherUserID || gotAppointment != appointmentID || gotStatus != "confirmed" {
					t.Fatalf("unexpected forwarded args")
				}
				return model.Appointment{AppointmentID: gotAppointment, Status: gotStatus}, nil
			},
		}}

		r := gin.New()
		r.Use(withTeacherClaims(teacherUserID))
		r.PATCH("/appointments/:appointment_id/status", h.UpdateAppointmentStatus)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPatch, "/appointments/"+appointmentID.String()+"/status", strings.NewReader(`{"status":"confirmed"}`))
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
			h := &TeacherScopeHandler{appointmentService: &fakeTeacherScopeAppointmentService{
				updateAppointmentStatusByTeacher: func(context.Context, uuid.UUID, uuid.UUID, string, string) (model.Appointment, error) {
					return model.Appointment{}, tt.err
				},
			}}

			r := gin.New()
			r.Use(withTeacherClaims(teacherUserID))
			r.PATCH("/appointments/:appointment_id/status", h.UpdateAppointmentStatus)

			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPatch, "/appointments/"+appointmentID.String()+"/status", strings.NewReader(`{"status":"confirmed"}`))
			req.Header.Set("Content-Type", "application/json")
			r.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
		})
	}
}
