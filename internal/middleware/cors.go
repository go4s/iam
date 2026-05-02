package middleware

import (
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

// CORSMiddleware 跨域中间件，允许 *.zhangli2946.club
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		allowed := isAllowedOrigin(origin)

		if allowed {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
			c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		}

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// isAllowedOrigin 检查 Origin 是否属于 *.zhangli2946.club
func isAllowedOrigin(origin string) bool {
	if origin == "" {
		return true
	}

	u, err := url.Parse(origin)
	if err != nil {
		return false
	}

	host := u.Host
	// 去掉端口
	if idx := strings.Index(host, ":"); idx != -1 {
		host = host[:idx]
	}

	// 匹配 *.zhangli2946.club
	if strings.HasSuffix(host, ".zhangli2946.club") {
		return true
	}
	// 匹配 *.zhangli2946.club
	if strings.HasSuffix(host, "localhost") {
		return true
	}

	return false
}
