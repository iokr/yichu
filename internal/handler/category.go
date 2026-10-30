package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/iokr/yichu/internal/dto"
	"github.com/iokr/yichu/internal/pkg/response"
)

// ListCategories 分类树。
func (h *Handler) ListCategories(c *gin.Context) {
	res, err := h.deps.Categories.Tree(c.Request.Context())
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, res)
}

// CreateCategory 新增分类。
func (h *Handler) CreateCategory(c *gin.Context) {
	var req dto.CategoryRequest
	if !h.bindJSON(c, &req) {
		return
	}
	res, err := h.deps.Categories.Create(c.Request.Context(), req)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, res)
}

// UpdateCategory 更新分类。
func (h *Handler) UpdateCategory(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var req dto.CategoryRequest
	if !h.bindJSON(c, &req) {
		return
	}
	res, err := h.deps.Categories.Update(c.Request.Context(), id, req)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, res)
}

// DeleteCategory 删除分类。
func (h *Handler) DeleteCategory(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	if err := h.deps.Categories.Delete(c.Request.Context(), id); err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, gin.H{"id": id})
}

// SortCategories 调整分类排序。
func (h *Handler) SortCategories(c *gin.Context) {
	var req dto.CategorySortRequest
	if !h.bindJSON(c, &req) {
		return
	}
	if err := h.deps.Categories.Sort(c.Request.Context(), req.Items); err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, gin.H{"ok": true})
}

// ListTags 标签列表（含关联数量）。
func (h *Handler) ListTags(c *gin.Context) {
	res, err := h.deps.Tags.List(c.Request.Context())
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, res)
}
