package web

import (
	"embed"
	"io/fs"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

//go:embed all:dist
var distFS embed.FS

// FS 返回前端构建产物文件系统。
func FS() (fs.FS, error) { return fs.Sub(distFS, "dist") }

const notBuiltHTML = `<!doctype html><html><head><meta charset="utf-8">
<meta name="viewport" content="width=device-width,initial-scale=1">
<title>衣橱</title></head>
<body style="font-family:system-ui;padding:64px;color:#6B6661;background:#FAFAF8">
<h2 style="color:#1F1D1B">前端尚未构建</h2>
<p>请执行 <code>make front</code> 或 <code>cd frontend &amp;&amp; npm run build</code> 后重新编译。</p>
<p>API 已正常运行，可访问 <code>/api/v1/version</code>。</p></body></html>`

// Register 挂载 SPA 静态资源与前端路由回退。
func Register(r *gin.Engine) error {
	sub, err := FS()
	if err != nil {
		return err
	}
	fileServer := http.FileServer(http.FS(sub))

	r.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path
		if strings.HasPrefix(path, "/api/") {
			c.JSON(http.StatusNotFound, gin.H{"code": 1004, "message": "resource not found", "data": nil})
			return
		}
		// 只允许 GET/HEAD 回退到前端。
		if c.Request.Method != http.MethodGet && c.Request.Method != http.MethodHead {
			c.JSON(http.StatusMethodNotAllowed, gin.H{"code": 1004, "message": "method not allowed", "data": nil})
			return
		}
		if f, ferr := sub.Open(strings.TrimPrefix(path, "/")); ferr == nil {
			_ = f.Close()
			if strings.HasPrefix(path, "/assets/") {
				c.Header("Cache-Control", "public, max-age=31536000, immutable")
			}
			fileServer.ServeHTTP(c.Writer, c.Request)
			return
		}
		index, ierr := fs.ReadFile(sub, "index.html")
		if ierr != nil {
			c.Header("Cache-Control", "no-store")
			c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(notBuiltHTML))
			return
		}
		c.Header("Cache-Control", "no-cache")
		c.Data(http.StatusOK, "text/html; charset=utf-8", index)
	})
	return nil
}
