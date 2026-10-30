package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/iokr/yichu/internal/pkg/errcode"
	"github.com/iokr/yichu/internal/pkg/response"
	"github.com/iokr/yichu/internal/service"
)

// gin.Context 中保存登录用户信息的 key。
const (
	ContextUserID   = "yichu.user_id"
	ContextUsername = "yichu.username"
)

// Auth 校验 Authorization: Bearer <token>。
func Auth(tokens *service.TokenManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := strings.TrimSpace(c.GetHeader("Authorization"))
		if len(header) < 8 || !strings.EqualFold(header[:7], "bearer ") {
			unauthorized(c)
			return
		}
		claims, err := tokens.Parse(strings.TrimSpace(header[7:]), service.TokenTypeAccess)
		if err != nil {
			unauthorized(c)
			return
		}
		c.Set(ContextUserID, claims.UserID)
		c.Set(ContextUsername, claims.Username)
		c.Next()
	}
}

func unauthorized(c *gin.Context) {
	c.Header("WWW-Authenticate", `Bearer realm="yichu"`)
	response.FailCode(c, http.StatusUnauthorized, errcode.CodeUnauthorized, "未登录或登录已过期")
}
