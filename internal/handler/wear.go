package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/iokr/yichu/internal/dto"
	"github.com/iokr/yichu/internal/pkg/response"
)

// CreateWearLog 穿着打卡。
func (h *Handler) CreateWearLog(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var req dto.WearLogRequest
	if !h.bindJSON(c, &req) {
		return
	}
	res, err := h.deps.Wear.Create(c.Request.Context(), id, req)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, res)
}

// ListWearLogs 穿着记录列表。
func (h *Handler) ListWearLogs(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	res, err := h.deps.Wear.List(c.Request.Context(), id)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, res)
}
