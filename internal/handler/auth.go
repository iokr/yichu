package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/iokr/yichu/internal/dto"
	"github.com/iokr/yichu/internal/pkg/response"
)

// Login 账密登录。
func (h *Handler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if !h.bindJSON(c, &req) {
		return
	}
	res, err := h.deps.Auth.Login(c.Request.Context(), req)
	if err != nil {
		response.Fail(c, err)
		return
	}
	h.deps.Log.Infof("登录成功 username=%s ip=%s", res.User.Username, c.ClientIP())
	response.OK(c, res)
}

// Refresh 刷新 Token。
func (h *Handler) Refresh(c *gin.Context) {
	var req dto.RefreshRequest
	if !h.bindJSON(c, &req) {
		return
	}
	res, err := h.deps.Auth.Refresh(c.Request.Context(), req.RefreshToken)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, res)
}

// Logout 退出登录（JWT 无状态，由前端清除本地凭证）。
func (h *Handler) Logout(c *gin.Context) {
	response.OK(c, gin.H{"ok": true})
}
