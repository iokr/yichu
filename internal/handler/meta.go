package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/iokr/yichu/internal/buildinfo"
	"github.com/iokr/yichu/internal/dto"
	"github.com/iokr/yichu/internal/pkg/dict"
	"github.com/iokr/yichu/internal/pkg/response"
)

// Meta 产品元信息（免登录）。
func (h *Handler) Meta(c *gin.Context) {
	response.OK(c, dto.MetaResponse{
		Name:      h.deps.Config.App.Name,
		Version:   buildinfo.Version,
		BuildTime: buildinfo.BuildTime,
		HidePrice: h.deps.Settings.HidePrice(),
	})
}

// Dicts 枚举字典（免登录）。
func (h *Handler) Dicts(c *gin.Context) {
	response.OK(c, dict.All())
}

// Version 版本信息（免登录）。
func (h *Handler) Version(c *gin.Context) {
	response.OK(c, dto.VersionResponse{
		Version:   buildinfo.Version,
		BuildTime: buildinfo.BuildTime,
	})
}
