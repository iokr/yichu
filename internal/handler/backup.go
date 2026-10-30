package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/iokr/yichu/internal/pkg/response"
)

// DownloadBackup 下载 JSON 备份。
func (h *Handler) DownloadBackup(c *gin.Context) {
	data, masked, err := h.deps.Backup.Export(c.Request.Context())
	if err != nil {
		response.Fail(c, err)
		return
	}
	filename := fmt.Sprintf("yichu-backup-%s.json", time.Now().Format("20060102-150405"))
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%q", filename))
	c.Header("Cache-Control", "no-store")
	c.Header("X-Price-Masked", strconv.FormatBool(masked))
	c.Data(http.StatusOK, "application/json; charset=utf-8", data)
}
