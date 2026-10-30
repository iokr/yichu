package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	gormmysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"github.com/iokr/yichu/internal/config"
	"github.com/iokr/yichu/internal/model"
)

const (
	pingMaxRetries = 10
	pingInterval   = 2 * time.Second
)

// Init 连接 MySQL（带重试）、自动迁移、种子数据与首次建账号。
func Init(cfg *config.Config) *gorm.DB {
	db, err := connect(cfg)
	if err != nil {
		log.Fatalf("❌ 数据库初始化失败: %v", err)
	}
	if cfg.Database.AutoMigrate {
		if err := Migrate(db); err != nil {
			log.Fatalf("❌ 数据库迁移失败: %v", err)
		}
	}
	if err := Seed(db, cfg); err != nil {
		log.Fatalf("❌ 初始化数据失败: %v", err)
	}
	return db
}

func connect(cfg *config.Config) (*gorm.DB, error) {
	dsn := ensureDSNParams(cfg.Database.DSN)
	sqlDB, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("打开数据库失败: %w", err)
	}
	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.Database.ConnMaxLifetimeMinutes) * time.Minute)

	var pingErr error
	for attempt := 1; attempt <= pingMaxRetries; attempt++ {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		pingErr = sqlDB.PingContext(ctx)
		cancel()
		if pingErr == nil {
			break
		}
		if attempt < pingMaxRetries {
			log.Printf("⏳ 数据库未就绪（%d/%d），%s 后重试: %v", attempt, pingMaxRetries, pingInterval, pingErr)
			time.Sleep(pingInterval)
		}
	}
	if pingErr != nil {
		_ = sqlDB.Close()
		return nil, fmt.Errorf("连接 MySQL 失败: %w", pingErr)
	}

	gormLogger := gormlogger.New(log.New(os.Stdout, "[gorm] ", log.LstdFlags), gormlogger.Config{
		SlowThreshold:             300 * time.Millisecond,
		LogLevel:                  gormlogger.Warn,
		IgnoreRecordNotFoundError: true,
		Colorful:                  false,
	})

	db, err := gorm.Open(gormmysql.New(gormmysql.Config{Conn: sqlDB}), &gorm.Config{
		Logger:                                   gormLogger,
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		_ = sqlDB.Close()
		return nil, fmt.Errorf("初始化 GORM 失败: %w", err)
	}
	return db, nil
}

// ensureDSNParams 保证时间字段可正确解析。
func ensureDSNParams(dsn string) string {
	if !strings.Contains(dsn, "parseTime=") {
		if strings.Contains(dsn, "?") {
			dsn += "&parseTime=True"
		} else {
			dsn += "?parseTime=True"
		}
	}
	if !strings.Contains(dsn, "loc=") {
		dsn += "&loc=Local"
	}
	return dsn
}

// Migrate 自动迁移全部表结构。
func Migrate(db *gorm.DB) error {
	return db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4").
		AutoMigrate(
			&model.User{},
			&model.Category{},
			&model.Item{},
			&model.ItemImage{},
			&model.Tag{},
			&model.WearLog{},
			&model.Setting{},
		)
}
