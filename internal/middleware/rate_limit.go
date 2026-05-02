package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimiter 基于内存的滑动窗口限流器
type RateLimiter struct {
	requests map[string][]time.Time // key (IP:path) -> timestamps
	limit    int                    // 窗口内最大请求数
	window   time.Duration          // 窗口大小
	mu       sync.RWMutex
}

// NewRateLimiter 创建限流器
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}
}

// Allow 检查是否允许请求
func (rl *RateLimiter) Allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-rl.window)

	// 获取该 key 的历史请求
	timestamps := rl.requests[key]

	// 清理过期的请求记录（滑动窗口）
	valid := timestamps[:0]
	for _, ts := range timestamps {
		if ts.After(cutoff) {
			valid = append(valid, ts)
		}
	}

	// 检查是否超过限制
	if len(valid) >= rl.limit {
		rl.requests[key] = valid
		return false
	}

	// 记录本次请求
	valid = append(valid, now)
	rl.requests[key] = valid
	return true
}

// RateLimitMiddleware Gin 中间件
func RateLimitMiddleware() gin.HandlerFunc {
	// 默认：每 1 秒最多 5 次请求
	rl := NewRateLimiter(5, time.Second)

	return func(c *gin.Context) {
		key := c.ClientIP() + ":" + c.Request.URL.Path

		if !rl.Allow(key) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"code":    "1005",
				"message": "请求过于频繁，请稍后重试",
				"data":    nil,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
