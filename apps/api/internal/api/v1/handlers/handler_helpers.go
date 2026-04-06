package handlers

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

// requireCurrentUserID đọc user_id từ JWT claims và tự trả lỗi HTTP nếu token/claims không hợp lệ.
func requireCurrentUserID(c *gin.Context) (uuid.UUID, bool) {
	claimsAny, exists := c.Get(middleware.CtxClaims)
	if !exists {
		response.Fail(c, http.StatusUnauthorized, "unauthorized")
		return uuid.Nil, false
	}

	claims, ok := claimsAny.(*auth.Claims)
	if !ok {
		response.Fail(c, http.StatusUnauthorized, "unauthorized")
		return uuid.Nil, false
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid user ID")
		return uuid.Nil, false
	}

	return userID, true
}

func parseTimeRange(fromRaw, toRaw string) (*time.Time, *time.Time, error) {
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
