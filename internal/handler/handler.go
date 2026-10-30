package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/iokr/yichu/internal/config"
	"github.com/iokr/yichu/internal/pkg/errcode"
	"github.com/iokr/yichu/internal/pkg/logx"
	"github.com/iokr/yichu/internal/pkg/response"
	"github.com/iokr/yichu/internal/service"
)

// Deps handler 依赖的服务与配置。
type Deps struct {
	Config     *config.Config
	Log        *logx.Logger
	Tokens     *service.TokenManager
	Auth       *service.AuthService
	Items      *service.ItemService
	Categories *service.CategoryService
	Tags       *service.TagService
	Upload     *service.UploadService
	Dashboard  *service.DashboardService
	Stats      *service.StatsService
	Wear       *service.WearService
	Settings   *service.SettingsService
	Backup     *service.BackupService
}

// Handler HTTP 层。
type Handler struct {
	deps Deps
}

// New 构造 Handler。
func New(deps Deps) *Handler { return &Handler{deps: deps} }

func (h *Handler) bindJSON(c *gin.Context, obj any) bool {
	if err := c.ShouldBindJSON(obj); err != nil {
		h.deps.Log.Debugf("请求体解析失败 %s: %v", c.FullPath(), err)
		response.Fail(c, errcode.Param("请求参数有误，请检查后重试"))
		return false
	}
	return true
}

func (h *Handler) bindQuery(c *gin.Context, obj any) bool {
	if err := c.ShouldBindQuery(obj); err != nil {
		h.deps.Log.Debugf("查询参数解析失败 %s: %v", c.FullPath(), err)
		response.Fail(c, errcode.Param("查询参数有误，请检查后重试"))
		return false
	}
	return true
}

func parseID(c *gin.Context) (uint64, bool) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		response.Fail(c, errcode.Param("无效的 ID"))
		return 0, false
	}
	return id, true
}
