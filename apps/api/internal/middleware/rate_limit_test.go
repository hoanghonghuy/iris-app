package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestIPFixedWindowRateLimitBlocksAfterMaxRequests(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(NewIPFixedWindowRateLimitWithConfig(FixedWindowRateLimitConfig{
		MaxRequests:  2,
		Window:       time.Minute,
		CleanupEvery: 8,
		StaleTTL:     5 * time.Minute,
	}))
	r.POST("/api/v1/auth/reset-password", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	makeReq := func() *httptest.ResponseRecorder {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/reset-password", nil)
		req.RemoteAddr = "127.0.0.1:54321"
		r.ServeHTTP(rec, req)
		return rec
	}

	first := makeReq()
	if first.Code != http.StatusOK {
		t.Fatalf("first request code = %d, want %d", first.Code, http.StatusOK)
	}

	second := makeReq()
	if second.Code != http.StatusOK {
		t.Fatalf("second request code = %d, want %d", second.Code, http.StatusOK)
	}

	third := makeReq()
	if third.Code != http.StatusTooManyRequests {
		t.Fatalf("third request code = %d, want %d", third.Code, http.StatusTooManyRequests)
	}
	if third.Header().Get("Retry-After") == "" {
		t.Fatalf("expected Retry-After header on 429 response")
	}
}
