package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/iokr/yichu/internal/buildinfo"
	"github.com/iokr/yichu/internal/config"
	"github.com/iokr/yichu/internal/database"
	"github.com/iokr/yichu/internal/pkg/logx"
	"github.com/iokr/yichu/internal/pkg/netutil"
	"github.com/iokr/yichu/internal/router"
	"github.com/iokr/yichu/web"
)

// Version 与 BuildTime 由 -ldflags 注入。
var (
	Version   = "dev"
	BuildTime = "unknown"
)

func main() {
	configPath := flag.String("c", "", "配置文件路径（默认依次查找 ./config.yaml、~/.yichu、/opt/yichu）")
	flag.Parse()

	buildinfo.Version = Version
	buildinfo.BuildTime = BuildTime

	cfg := config.LoadFrom(*configPath)

	logger, err := logx.New(cfg.Log.File, cfg.Log.Level)
	if err != nil {
		log.Fatalf("❌ 初始化日志失败: %v", err)
	}
	defer func() { _ = logger.Close() }()
	log.SetOutput(logger.Writer())
	log.SetFlags(0)

	logger.Infof("🌿 %s %s (build %s) 启动中…", cfg.App.Name, Version, BuildTime)

	db := database.Init(cfg)
	r := router.New(cfg, db, logger)
	if err := web.Register(r); err != nil {
		logger.Errorf("❌ 挂载前端失败: %v", err)
		os.Exit(1)
	}

	srv := &http.Server{
		Addr:              cfg.Server.Addr,
		Handler:           r,
		ReadHeaderTimeout: 10 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	go func() {
		logger.Infof("🌿 %s %s (build %s) 已启动: http://%s", cfg.App.Name, Version, BuildTime, cfg.Server.Addr)
		if ip := netutil.LocalIP(); ip != "" {
			logger.Infof("📱 局域网访问: http://%s%s", ip, cfg.Server.Addr)
		}
		var serveErr error
		if cfg.Server.TLS.Enabled {
			serveErr = srv.ListenAndServeTLS(cfg.Server.TLS.CertFile, cfg.Server.TLS.KeyFile)
		} else {
			serveErr = srv.ListenAndServe()
		}
		if serveErr != nil && !errors.Is(serveErr, http.ErrServerClosed) {
			logger.Errorf("❌ 服务异常: %v", serveErr)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Warnf("优雅退出超时: %v", err)
	}
	logger.Infof("👋 %s 已退出", cfg.App.Name)
}
