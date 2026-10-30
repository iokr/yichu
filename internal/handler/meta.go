package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/iokr/yichu/internal/buildinfo"
	"github.com/iokr/yichu/internal/dto"
	"github.com/iokr/yichu/internal/pkg/dict"
	"github.com/iokr/yichu/internal/pkg/response"
)

// Meta 产品元信息（免登录）。
func (h *Handler) Meta(c *gin.Context) {
	response.OK(c, dto.MetaResponse{
		Name:       h.deps.Config.App.Name,
		Version:    buildinfo.Version,
		BuildTime:  buildinfo.BuildTime,
		HidePrice:  h.deps.Settings.HidePrice(),
		PWAEnabled: true,
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

// Manifest 动态生成 PWA manifest，产品名随配置变化。
func (h *Handler) Manifest(c *gin.Context) {
	name := h.deps.Config.App.Name
	c.Header("Content-Type", "application/manifest+json; charset=utf-8")
	c.Header("Cache-Control", "public, max-age=3600")
	c.JSON(http.StatusOK, gin.H{
		"name":             name,
		"short_name":       name,
		"description":      name + " · 家庭衣物数字化管理",
		"start_url":        "/",
		"scope":            "/",
		"display":          "standalone",
		"orientation":      "portrait",
		"background_color": "#FAFAF8",
		"theme_color":      "#A67C89",
		"lang":             "zh-CN",
		"icons": []gin.H{
			{"src": "/icons/icon-192.png", "sizes": "192x192", "type": "image/png"},
			{"src": "/icons/icon-512.png", "sizes": "512x512", "type": "image/png"},
			{"src": "/icons/icon-maskable-512.png", "sizes": "512x512", "type": "image/png", "purpose": "maskable"},
		},
	})
}
