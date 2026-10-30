package response

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/iokr/yichu/internal/pkg/errcode"
)

// Body 统一响应结构：{ "code": 0, "message": "ok", "data": {} }。
type Body struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

// OK 返回成功响应。
func OK(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Body{Code: errcode.CodeOK, Message: "ok", Data: data})
}

// Fail 按业务错误返回；未知错误统一为 5000，避免泄漏内部细节。
func Fail(c *gin.Context, err error) {
	e, ok := errcode.As(err)
	if !ok {
		log.Printf("[error] 未包装的错误: %v", err)
		e = errcode.Internal(err)
	} else if e.Code == errcode.CodeServer {
		log.Printf("[error] %v", err)
	}
	c.AbortWithStatusJSON(e.Status, Body{Code: e.Code, Message: e.Message, Data: nil})
}

// FailCode 直接以业务码返回。
func FailCode(c *gin.Context, status, code int, message string) {
	c.AbortWithStatusJSON(status, Body{Code: code, Message: message, Data: nil})
}
