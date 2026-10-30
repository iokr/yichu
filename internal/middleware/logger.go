package middleware

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/iokr/yichu/internal/pkg/logx"
)

// AccessLogger 记录请求日志，不记录请求体、响应体与价格。
func AccessLogger(logger *logx.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		c.Next()
		if strings.HasPrefix(path, "/assets/") || strings.HasPrefix(path, "/uploads/") {
			return
		}
		logger.Infof("%s %s %d %s ip=%s ua=%q",
			c.Request.Method,
			path,
			c.Writer.Status(),
			time.Since(start).Round(time.Millisecond),
			c.ClientIP(),
			c.Request.UserAgent(),
		)
	}
}
