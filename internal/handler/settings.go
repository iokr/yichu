package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/iokr/yichu/internal/buildinfo"
	"github.com/iokr/yichu/internal/dto"
	"github.com/iokr/yichu/internal/pkg/errcode"
	"github.com/iokr/yichu/internal/pkg/response"
)

// GetSettings 读取运行时设置。
func (h *Handler) GetSettings(c *gin.Context) {
	response.OK(c, h.settingsResponse())
}

// UpdateSettings 更新运行时设置。
func (h *Handler) UpdateSettings(c *gin.Context) {
	var req dto.SettingsUpdateRequest
	if !h.bindJSON(c, &req) {
		return
	}
	if req.HidePrice != nil {
		if err := h.deps.Settings.SetHidePrice(*req.HidePrice); err != nil {
			response.Fail(c, errcode.Internal(err))
			return
		}
		h.deps.Log.Infof("价格可见性切换为 hide_price=%v", *req.HidePrice)
	}
	response.OK(c, h.settingsResponse())
}

func (h *Handler) settingsResponse() dto.SettingsResponse {
	return dto.SettingsResponse{
		HidePrice:                  h.deps.Settings.HidePrice(),
		PriceToggleRequirePassword: h.deps.Settings.RequirePassword(),
		AppName:                    h.deps.Config.App.Name,
		Version:                    buildinfo.Version,
		BuildTime:                  buildinfo.BuildTime,
		PWAEnabled:                 true,
	}
}
