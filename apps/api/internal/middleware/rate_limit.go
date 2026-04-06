package middleware

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/hoanghonghuy/iris-app/apps/api/internal/response"
)

type fixedWindowCounter struct {
	// windowStart là mốc bắt đầu cửa sổ giới hạn hiện tại của key.
	windowStart time.Time
	// lastSeen dùng để dọn các key cũ không còn traffic.
	lastSeen time.Time
	count    int
}

// FixedWindowRateLimitConfig cấu hình cho limiter fixed-window theo IP + route.
type FixedWindowRateLimitConfig struct {
	MaxRequests  int
	Window       time.Duration
	CleanupEvery int
	StaleTTL     time.Duration
}

// NewIPFixedWindowRateLimitWithConfig tạo middleware với cấu hình chi tiết.
func NewIPFixedWindowRateLimitWithConfig(cfg FixedWindowRateLimitConfig) gin.HandlerFunc {
	if cfg.MaxRequests <= 0 {
		cfg.MaxRequests = 1
	}
	if cfg.Window <= 0 {
		cfg.Window = time.Minute
	}
	if cfg.CleanupEvery <= 0 {
		cfg.CleanupEvery = 256
	}
	if cfg.StaleTTL <= 0 {
		cfg.StaleTTL = 5 * cfg.Window
	}

	var mu sync.Mutex
	counters := make(map[string]fixedWindowCounter)
	requestCount := 0

	return func(c *gin.Context) {
		now := time.Now()
		key := fmt.Sprintf("%s|%s", c.Request.URL.Path, c.ClientIP())

		mu.Lock()
		counter := counters[key]
		if counter.windowStart.IsZero() || now.Sub(counter.windowStart) >= cfg.Window {
			counter.windowStart = now
			counter.count = 0
		}
		counter.lastSeen = now

		if counter.count >= cfg.MaxRequests {
			retryAfter := int(cfg.Window.Seconds())
			if elapsed := now.Sub(counter.windowStart); elapsed < cfg.Window {
				retryAfter = int((cfg.Window - elapsed).Seconds())
				if retryAfter < 1 {
					retryAfter = 1
				}
			}

			c.Header("Retry-After", fmt.Sprintf("%d", retryAfter))
			mu.Unlock()
			response.Fail(c, http.StatusTooManyRequests, "too many requests")
			c.Abort()
			return
		}

		counter.count++
		counters[key] = counter
		requestCount++

		if requestCount%cfg.CleanupEvery == 0 {
			for staleKey, staleCounter := range counters {
				if now.Sub(staleCounter.lastSeen) > cfg.StaleTTL {
					delete(counters, staleKey)
				}
			}
		}
		mu.Unlock()

		c.Next()
	}
}
