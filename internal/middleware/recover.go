package middleware

import (
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"

	"github.com/iokr/yichu/internal/pkg/errcode"
	"github.com/iokr/yichu/internal/pkg/logx"
	"github.com/iokr/yichu/internal/pkg/response"
)

// Recover 捕获 panic，记录堆栈并返回统一错误。
func Recover(logger *logx.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if rec := recover(); rec != nil {
				logger.Errorf("panic: %v\n%s", rec, string(debug.Stack()))
				if c.Writer.Written() {
					c.Abort()
					return
				}
				response.FailCode(c, http.StatusInternalServerError, errcode.CodeServer, "服务器开小差了，请稍后重试")
			}
		}()
		c.Next()
	}
}
