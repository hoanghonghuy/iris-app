package chathandlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/hoanghonghuy/iris-app/apps/api/internal/auth"
)

func TestHandleWS_PreUpgradeValidation(t *testing.T) {
	gin.SetMode(gin.TestMode)
	secret := "test-secret"

	t.Run("missing token in protocol", func(t *testing.T) {
		h := &ChatHandler{jwtSecret: secret, allowedOrigins: map[string]struct{}{}}
		r := gin.New()
		r.GET("/chat/ws", h.HandleWS)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/chat/ws", nil)
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
		}
		if got := decodeChatError(t, rec); got != "missing token in Sec-WebSocket-Protocol" {
			t.Fatalf("error = %q, want %q", got, "missing token in Sec-WebSocket-Protocol")
		}
	})

	t.Run("invalid token", func(t *testing.T) {
		h := &ChatHandler{jwtSecret: secret, allowedOrigins: map[string]struct{}{}}
		r := gin.New()
		r.GET("/chat/ws", h.HandleWS)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/chat/ws", nil)
		req.Header.Set("Sec-WebSocket-Protocol", "Bearer, not-a-jwt")
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
		}
		if got := decodeChatError(t, rec); got != "invalid token" {
			t.Fatalf("error = %q, want %q", got, "invalid token")
		}
	})

	t.Run("invalid user id in token", func(t *testing.T) {
		token, err := auth.Sign(secret, time.Minute, "not-a-uuid", "u@example.com", []string{"PARENT"}, "")
		if err != nil {
			t.Fatalf("auth.Sign() error = %v", err)
		}

		h := &ChatHandler{jwtSecret: secret, allowedOrigins: map[string]struct{}{}}
		r := gin.New()
		r.GET("/chat/ws", h.HandleWS)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/chat/ws", nil)
		req.Header.Set("Sec-WebSocket-Protocol", "Bearer, "+token)
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
		}
		if got := decodeChatError(t, rec); got != "invalid user ID in token" {
			t.Fatalf("error = %q, want %q", got, "invalid user ID in token")
		}
	})
}
