package errcode

import (
	"errors"
	"fmt"
	"net/http"
)

// 统一业务错误码，见设计文档 8.4。
const (
	CodeOK           = 0
	CodeParam        = 1001
	CodeUnauthorized = 1002
	CodeForbidden    = 1003
	CodeNotFound     = 1004
	CodeServer       = 5000
)

// Error 是携带业务码与 HTTP 状态码的错误类型。
// cause 仅用于日志，绝不返回给客户端。
type Error struct {
	Code    int
	Message string
	Status  int
	cause   error
}

func (e *Error) Error() string {
	if e.cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.cause)
	}
	return e.Message
}

func (e *Error) Unwrap() error { return e.cause }

func defaultStatus(code int) int {
	switch code {
	case CodeParam:
		return http.StatusBadRequest
	case CodeUnauthorized:
		return http.StatusUnauthorized
	case CodeForbidden:
		return http.StatusForbidden
	case CodeNotFound:
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}

// New 构造业务错误，HTTP 状态码由业务码推导。
func New(code int, message string) *Error {
	return &Error{Code: code, Message: message, Status: defaultStatus(code)}
}

// Param 参数错误。
func Param(format string, args ...any) *Error {
	return &Error{Code: CodeParam, Message: fmt.Sprintf(format, args...), Status: http.StatusBadRequest}
}

// Unauthorized 未登录 / 凭证失效。
func Unauthorized(message string) *Error {
	return &Error{Code: CodeUnauthorized, Message: message, Status: http.StatusUnauthorized}
}

// Forbidden 无权限或业务规则禁止。
func Forbidden(format string, args ...any) *Error {
	return &Error{Code: CodeForbidden, Message: fmt.Sprintf(format, args...), Status: http.StatusForbidden}
}

// NotFound 资源不存在。
func NotFound(message string) *Error {
	return &Error{Code: CodeNotFound, Message: message, Status: http.StatusNotFound}
}

// Internal 服务器错误，原始错误仅进入日志。
func Internal(err error) *Error {
	return &Error{Code: CodeServer, Message: "服务器开小差了，请稍后重试", Status: http.StatusInternalServerError, cause: err}
}

// Wrap 使用指定业务码包装底层错误。
func Wrap(code int, message string, err error) *Error {
	return &Error{Code: code, Message: message, Status: defaultStatus(code), cause: err}
}

// As 提取业务错误。
func As(err error) (*Error, bool) {
	var e *Error
	if errors.As(err, &e) {
		return e, true
	}
	return nil, false
}
