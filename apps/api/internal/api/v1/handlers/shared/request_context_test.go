package shared

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/hoanghonghuy/iris-app/apps/api/internal/auth"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/middleware"
)

func TestNormalizePagination_DefaultAndBounds(t *testing.T) {
	tests := []struct {
		name       string
		limit      int
		offset     int
		wantLimit  int
		wantOffset int
	}{
		{
			name:       "default values when limit and offset are invalid",
			limit:      0,
			offset:     -1,
			wantLimit:  20,
			wantOffset: 0,
		},
		{
			name:       "clamp limit to max",
			limit:      101,
			offset:     5,
			wantLimit:  100,
			wantOffset: 5,
		},
		{
			name:       "keep valid values",
			limit:      15,
			offset:     10,
			wantLimit:  15,
			wantOffset: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotLimit, gotOffset := NormalizePagination(tt.limit, tt.offset)
			if gotLimit != tt.wantLimit {
				t.Fatalf("expected limit %d, got %d", tt.wantLimit, gotLimit)
			}
			if gotOffset != tt.wantOffset {
				t.Fatalf("expected offset %d, got %d", tt.wantOffset, gotOffset)
			}
		})
	}
}

func TestRequireCurrentClaims(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name       string
		setup      func(*gin.Context)
		wantStatus int
		wantOK     bool
	}{
		{
			name:       "missing claims",
			setup:      func(*gin.Context) {},
			wantStatus: http.StatusUnauthorized,
			wantOK:     false,
		},
		{
			name: "invalid claim type",
			setup: func(c *gin.Context) {
				c.Set(middleware.CtxClaims, "not-claims")
			},
			wantStatus: http.StatusUnauthorized,
			wantOK:     false,
		},
		{
			name: "valid claims",
			setup: func(c *gin.Context) {
				c.Set(middleware.CtxClaims, &auth.Claims{UserID: "11111111-1111-1111-1111-111111111111"})
			},
			wantStatus: http.StatusOK,
			wantOK:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.New()
			r.GET("/claims", func(c *gin.Context) {
				tt.setup(c)
				_, ok := RequireCurrentClaims(c)
				if ok {
					c.Status(http.StatusOK)
				}
			})

			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/claims", nil)
			r.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
			if !tt.wantOK {
				if got := decodeError(t, rec); got != "unauthorized" {
					t.Fatalf("error = %q, want %q", got, "unauthorized")
				}
			}
		})
	}
}

func TestRequireCurrentUserID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name       string
		claims     *auth.Claims
		wantStatus int
		wantOK     bool
		wantError  string
	}{
		{
			name:       "missing claims",
			claims:     nil,
			wantStatus: http.StatusUnauthorized,
			wantOK:     false,
			wantError:  "unauthorized",
		},
		{
			name:       "invalid user id format",
			claims:     &auth.Claims{UserID: "not-a-uuid"},
			wantStatus: http.StatusBadRequest,
			wantOK:     false,
			wantError:  "invalid user ID",
		},
		{
			name:       "valid user id",
			claims:     &auth.Claims{UserID: "11111111-1111-1111-1111-111111111111"},
			wantStatus: http.StatusOK,
			wantOK:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.New()
			r.GET("/user", func(c *gin.Context) {
				if tt.claims != nil {
					c.Set(middleware.CtxClaims, tt.claims)
				}
				_, ok := RequireCurrentUserID(c)
				if ok {
					c.Status(http.StatusOK)
				}
			})

			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/user", nil)
			r.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
			if !tt.wantOK {
				if got := decodeError(t, rec); got != tt.wantError {
					t.Fatalf("error = %q, want %q", got, tt.wantError)
				}
			}
		})
	}
}

func decodeError(t *testing.T, rec *httptest.ResponseRecorder) string {
	t.Helper()
	var body map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}
	return body["error"]
}

func TestParseTimeRange(t *testing.T) {
	from, to, err := ParseTimeRange("2026-04-07T10:00:00Z", "2026-04-07T11:00:00Z")
	if err != nil {
		t.Fatalf("ParseTimeRange() error = %v", err)
	}
	if from == nil || to == nil {
		t.Fatalf("expected both from/to to be non-nil")
	}
	if !to.After(*from) {
		t.Fatalf("expected to > from")
	}

	from, to, err = ParseTimeRange("", "")
	if err != nil {
		t.Fatalf("ParseTimeRange() error = %v", err)
	}
	if from != nil || to != nil {
		t.Fatalf("expected nil pointers when both query params are empty")
	}

	_, _, err = ParseTimeRange("bad-time", "")
	if err == nil {
		t.Fatalf("expected parse error for invalid from")
	}

	_, _, err = ParseTimeRange("", "bad-time")
	if err == nil {
		t.Fatalf("expected parse error for invalid to")
	}

	_, _, err = ParseTimeRange("2026-04-07T12:00:00Z", "2026-04-07T11:00:00Z")
	if err == nil {
		t.Fatalf("expected range validation error when from > to")
	}

	fromRaw := "2026-04-07T12:00:00+07:00"
	toRaw := "2026-04-07T06:00:00+01:00"
	expectedFrom, _ := time.Parse(time.RFC3339, fromRaw)
	expectedTo, _ := time.Parse(time.RFC3339, toRaw)

	from, to, err = ParseTimeRange(fromRaw, toRaw)
	if err != nil {
		t.Fatalf("ParseTimeRange() timezone offset error = %v", err)
	}
	if from == nil || to == nil {
		t.Fatalf("expected non-nil from/to for timezone inputs")
	}
	if !from.Equal(expectedFrom) || !to.Equal(expectedTo) {
		t.Fatalf("timezone parse mismatch: got from=%v to=%v", *from, *to)
	}
	if !from.Equal(*to) {
		t.Fatalf("expected same instant across different timezone offsets")
	}
}

func TestParsePagination(t *testing.T) {
	tests := []struct {
		name      string
		limitRaw  string
		offsetRaw string
		wantLimit int
		wantOff   int
	}{
		{name: "default values", limitRaw: "", offsetRaw: "", wantLimit: 20, wantOff: 0},
		{name: "valid values", limitRaw: "30", offsetRaw: "15", wantLimit: 30, wantOff: 15},
		{name: "invalid values fallback", limitRaw: "x", offsetRaw: "y", wantLimit: 20, wantOff: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotLimit, gotOff := ParsePagination(tt.limitRaw, tt.offsetRaw)
			if gotLimit != tt.wantLimit {
				t.Fatalf("limit = %d, want %d", gotLimit, tt.wantLimit)
			}
			if gotOff != tt.wantOff {
				t.Fatalf("offset = %d, want %d", gotOff, tt.wantOff)
			}
		})
	}
}
