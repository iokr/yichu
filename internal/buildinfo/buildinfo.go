package buildinfo

// Version 与 BuildTime 由构建时 -ldflags 注入：
// -X 'main.Version=...' -X 'main.BuildTime=...'，main 启动后同步到本包。
var (
	Version   = "dev"
	BuildTime = "unknown"
)
