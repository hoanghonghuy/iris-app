package shared

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/hoanghonghuy/iris-app/apps/api/internal/auth"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/middleware"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/response"
)

// RequireCurrentClaims extracts JWT claims from context and writes auth error response on failure.
func RequireCurrentClaims(c *gin.Context) (*auth.Claims, bool) {
	claimsAny, exists := c.Get(middleware.CtxClaims)
	if !exists {
		response.Fail(c, http.StatusUnauthorized, "unauthorized")
		return nil, false
	}

	claims, ok := claimsAny.(*auth.Claims)
	if !ok {
		response.Fail(c, http.StatusUnauthorized, "unauthorized")
		return nil, false
	}

	return claims, true
}

// RequireCurrentUser extracts current user UUID from claims and writes error response on failure.
func RequireCurrentUser(c *gin.Context) (uuid.UUID, *auth.Claims, bool) {
	claims, ok := RequireCurrentClaims(c)
	if !ok {
		return uuid.Nil, nil, false
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid user ID")
		return uuid.Nil, nil, false
	}

	return userID, claims, true
}

// RequireCurrentUserID extracts only current user UUID and writes error response on failure.
func RequireCurrentUserID(c *gin.Context) (uuid.UUID, bool) {
	userID, _, ok := RequireCurrentUser(c)
	if !ok {
		return uuid.Nil, false
	}
	return userID, true
}

// ParseTimeRange parses optional RFC3339 range query params.
func ParseTimeRange(fromRaw, toRaw string) (*time.Time, *time.Time, error) {
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

// ParsePagination parses limit/offset with default values.
func ParsePagination(limitRaw, offsetRaw string) (int, int) {
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
