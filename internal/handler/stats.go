package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/iokr/yichu/internal/pkg/response"
)

// StatsOverview 统计概览。
func (h *Handler) StatsOverview(c *gin.Context) {
	res, err := h.deps.Stats.Overview(c.Request.Context())
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, res)
}

// StatsByCategory 分类统计。
func (h *Handler) StatsByCategory(c *gin.Context) {
	res, err := h.deps.Stats.ByCategory(c.Request.Context())
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, res)
}

// StatsByBrand 品牌统计。
func (h *Handler) StatsByBrand(c *gin.Context) {
	res, err := h.deps.Stats.ByBrand(c.Request.Context())
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, res)
}

// StatsByYear 年度统计。
func (h *Handler) StatsByYear(c *gin.Context) {
	res, err := h.deps.Stats.ByYear(c.Request.Context())
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, res)
}

// StatsByMonth 月度购买统计。
func (h *Handler) StatsByMonth(c *gin.Context) {
	year, _ := strconv.Atoi(c.DefaultQuery("year", "0"))
	res, err := h.deps.Stats.ByMonth(c.Request.Context(), year)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, res)
}

// StatsBySeason 季节统计。
func (h *Handler) StatsBySeason(c *gin.Context) {
	res, err := h.deps.Stats.BySeason(c.Request.Context())
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, res)
}

// StatsByColor 颜色统计。
func (h *Handler) StatsByColor(c *gin.Context) {
	res, err := h.deps.Stats.ByColor(c.Request.Context())
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, res)
}

// StatsCostPerWear 单次成本排行。
func (h *Handler) StatsCostPerWear(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	res, err := h.deps.Stats.CostPerWear(c.Request.Context(), limit)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, res)
}
