package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/iokr/yichu/internal/pkg/errcode"
	"github.com/iokr/yichu/internal/pkg/response"
)

// Upload 上传图片。
func (h *Handler) Upload(c *gin.Context) {
	maxBytes := h.deps.Config.MaxUploadBytes()
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxBytes+1024*1024)

	file, err := c.FormFile("file")
	if err != nil {
		response.Fail(c, errcode.Param("请选择要上传的图片"))
		return
	}
	res, err := h.deps.Upload.Save(c.Request.Context(), file)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, res)
}
