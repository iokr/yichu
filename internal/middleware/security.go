package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// SecurityHeaders 基础安全响应头与静态图片缓存策略。
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Header("X-Frame-Options", "SAMEORIGIN")
		// 图片文件名由服务端 UUID 生成，可长期缓存。
		if strings.HasPrefix(c.Request.URL.Path, "/uploads/") {
			c.Header("Cache-Control", "public, max-age=2592000, immutable")
		}
		c.Next()
	}
}

// NoStore 用于鉴权与设置接口，禁止任何缓存。
func NoStore() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Cache-Control", "no-store")
		c.Next()
	}
}
