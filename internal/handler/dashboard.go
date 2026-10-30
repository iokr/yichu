package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/iokr/yichu/internal/pkg/response"
)

// Dashboard 仪表盘数据。
func (h *Handler) Dashboard(c *gin.Context) {
	res, err := h.deps.Dashboard.Get(c.Request.Context())
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, res)
}
