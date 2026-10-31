package router

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/iokr/yichu/internal/buildinfo"
	"github.com/iokr/yichu/internal/config"
	"github.com/iokr/yichu/internal/handler"
	"github.com/iokr/yichu/internal/middleware"
	"github.com/iokr/yichu/internal/pkg/logx"
	"github.com/iokr/yichu/internal/repository"
	"github.com/iokr/yichu/internal/service"
)

// New 组装仓储、服务、handler 与路由。
// 分层依赖：router → handler → service → repository → model。
func New(cfg *config.Config, db *gorm.DB, logger *logx.Logger) *gin.Engine {
	if cfg.IsRelease() {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	userRepo := repository.NewUserRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	itemRepo := repository.NewItemRepository(db)
	tagRepo := repository.NewTagRepository(db)
	wearRepo := repository.NewWearLogRepository(db)
	settingRepo := repository.NewSettingRepository(db)
	statsRepo := repository.NewStatsRepository(db)

	settingsSvc := service.NewSettingsService(settingRepo, cfg)
	if err := settingsSvc.Init(); err != nil {
		log.Fatalf("❌ 加载设置失败: %v", err)
	}
	storage, err := service.NewLocalStorage(cfg.Storage.UploadDir)
	if err != nil {
		log.Fatalf("❌ 初始化图片存储失败: %v", err)
	}

	tokens := service.NewTokenManager(cfg.Auth.JWTSecret, cfg.Auth.TokenTTLHours)
	authSvc := service.NewAuthService(userRepo, tokens)
	itemSvc := service.NewItemService(itemRepo, categoryRepo, settingsSvc, cfg)
	categorySvc := service.NewCategoryService(categoryRepo)
	tagSvc := service.NewTagService(tagRepo)
	wearSvc := service.NewWearService(itemRepo, wearRepo)
	uploadSvc := service.NewUploadService(storage, cfg.MaxUploadBytes(), logger)
	dashboardSvc := service.NewDashboardService(itemRepo, statsRepo, categoryRepo, settingsSvc)
	statsSvc := service.NewStatsService(itemRepo, statsRepo, settingsSvc)
	backupSvc := service.NewBackupService(
		cfg.App.Name, buildinfo.Version,
		userRepo, categorySvc, itemSvc, tagSvc, wearSvc, settingsSvc, settingRepo,
	)

	h := handler.New(handler.Deps{
		Config:     cfg,
		Log:        logger,
		Tokens:     tokens,
		Auth:       authSvc,
		Items:      itemSvc,
		Categories: categorySvc,
		Tags:       tagSvc,
		Upload:     uploadSvc,
		Dashboard:  dashboardSvc,
		Stats:      statsSvc,
		Wear:       wearSvc,
		Settings:   settingsSvc,
		Backup:     backupSvc,
	})

	r := gin.New()
	r.Use(
		middleware.Recover(logger),
		middleware.SecurityHeaders(),
		middleware.CORS(),
		middleware.AccessLogger(logger),
	)

	// 图片走显式路由，不会被 SPA 回退拦截。
	r.StaticFS("/uploads", http.Dir(storage.Root()))

	api := r.Group("/api/v1")
	{
		api.GET("/version", h.Version)
		api.GET("/meta", h.Meta)
		api.GET("/dicts", h.Dicts)
		api.POST("/auth/login", middleware.NoStore(), h.Login)
		api.POST("/auth/refresh", middleware.NoStore(), h.Refresh)
	}

	auth := api.Group("", middleware.Auth(tokens))
	{
		auth.POST("/auth/logout", middleware.NoStore(), h.Logout)

		auth.GET("/dashboard", h.Dashboard)
		auth.GET("/tags", h.ListTags)

		items := auth.Group("/items")
		{
			items.GET("", h.ListItems)
			items.POST("", h.CreateItem)
			items.POST("/batch", h.BatchItems)
			items.GET("/export", h.ExportItems)
			items.GET("/:id", h.GetItem)
			items.PUT("/:id", h.UpdateItem)
			items.DELETE("/:id", h.DeleteItem)
			items.POST("/:id/wear-logs", h.CreateWearLog)
			items.GET("/:id/wear-logs", h.ListWearLogs)
		}

		categories := auth.Group("/categories")
		{
			categories.GET("", h.ListCategories)
			categories.POST("", h.CreateCategory)
			categories.PUT("/sort", h.SortCategories)
			categories.PUT("/:id", h.UpdateCategory)
			categories.DELETE("/:id", h.DeleteCategory)
		}

		stats := auth.Group("/stats")
		{
			stats.GET("/overview", h.StatsOverview)
			stats.GET("/by-category", h.StatsByCategory)
			stats.GET("/by-brand", h.StatsByBrand)
			stats.GET("/by-year", h.StatsByYear)
			stats.GET("/by-season", h.StatsBySeason)
			stats.GET("/by-color", h.StatsByColor)
			stats.GET("/cost-per-wear", h.StatsCostPerWear)
		}

		auth.POST("/upload", h.Upload)
		auth.GET("/settings", middleware.NoStore(), h.GetSettings)
		auth.PUT("/settings", middleware.NoStore(), h.UpdateSettings)
		auth.GET("/backup", h.DownloadBackup)
	}

	return r
}
