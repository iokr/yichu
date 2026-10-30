package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/iokr/yichu/internal/dto"
	"github.com/iokr/yichu/internal/pkg/response"
)

// ListItems 衣物列表。
func (h *Handler) ListItems(c *gin.Context) {
	var q dto.ItemQuery
	if !h.bindQuery(c, &q) {
		return
	}
	res, err := h.deps.Items.List(c.Request.Context(), q)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, res)
}

// GetItem 衣物详情。
func (h *Handler) GetItem(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	res, err := h.deps.Items.Get(c.Request.Context(), id)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, res)
}

// CreateItem 新增衣物。
func (h *Handler) CreateItem(c *gin.Context) {
	var req dto.ItemRequest
	if !h.bindJSON(c, &req) {
		return
	}
	res, err := h.deps.Items.Create(c.Request.Context(), req)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, res)
}

// UpdateItem 更新衣物。
func (h *Handler) UpdateItem(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var req dto.ItemRequest
	if !h.bindJSON(c, &req) {
		return
	}
	res, err := h.deps.Items.Update(c.Request.Context(), id, req)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, res)
}

// DeleteItem 删除衣物。
func (h *Handler) DeleteItem(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	if err := h.deps.Items.Delete(c.Request.Context(), id); err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, gin.H{"id": id})
}

// BatchItems 批量删除 / 批量改状态。
func (h *Handler) BatchItems(c *gin.Context) {
	var req dto.BatchRequest
	if !h.bindJSON(c, &req) {
		return
	}
	if err := h.deps.Items.Batch(c.Request.Context(), req); err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, gin.H{"affected": len(req.IDs)})
}

// ExportItems 导出 CSV / Excel。
func (h *Handler) ExportItems(c *gin.Context) {
	var q dto.ItemQuery
	if !h.bindQuery(c, &q) {
		return
	}
	result, err := h.deps.Items.Export(c.Request.Context(), q, c.DefaultQuery("format", "csv"))
	if err != nil {
		response.Fail(c, err)
		return
	}
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%q", result.Filename))
	c.Header("Cache-Control", "no-store")
	c.Header("X-Price-Masked", strconv.FormatBool(result.Masked))
	c.Data(http.StatusOK, result.ContentType, result.Data)
}
