package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/viper"
)

// Config 对应 config.yaml，全部字段均可被 YICHU_ 前缀环境变量覆盖。
type Config struct {
	App      AppConfig      `mapstructure:"app"`
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Auth     AuthConfig     `mapstructure:"auth"`
	Storage  StorageConfig  `mapstructure:"storage"`
	Log      LogConfig      `mapstructure:"log"`
}

// AppConfig 应用配置。
type AppConfig struct {
	Name                       string `mapstructure:"name"`
	HidePrice                  bool   `mapstructure:"hide_price"`
	PriceToggleRequirePassword bool   `mapstructure:"price_toggle_require_password"`
}

// ServerConfig HTTP 服务配置。
type ServerConfig struct {
	Addr string    `mapstructure:"addr"`
	Mode string    `mapstructure:"mode"`
	TLS  TLSConfig `mapstructure:"tls"`
}

// TLSConfig 内置 TLS 配置，与 Caddy 方案互斥。
type TLSConfig struct {
	Enabled  bool   `mapstructure:"enabled"`
	CertFile string `mapstructure:"cert_file"`
	KeyFile  string `mapstructure:"key_file"`
}

// DatabaseConfig MySQL 配置。
type DatabaseConfig struct {
	Driver                 string `mapstructure:"driver"`
	DSN                    string `mapstructure:"dsn"`
	MaxOpenConns           int    `mapstructure:"max_open_conns"`
	MaxIdleConns           int    `mapstructure:"max_idle_conns"`
	ConnMaxLifetimeMinutes int    `mapstructure:"conn_max_lifetime_minutes"`
	AutoMigrate            bool   `mapstructure:"auto_migrate"`
}

// AuthConfig 鉴权配置。
type AuthConfig struct {
	JWTSecret     string `mapstructure:"jwt_secret"`
	TokenTTLHours int    `mapstructure:"token_ttl_hours"`
	InitUsername  string `mapstructure:"init_username"`
	InitPassword  string `mapstructure:"init_password"`
}

// StorageConfig 图片存储配置。
type StorageConfig struct {
	UploadDir string `mapstructure:"upload_dir"`
	MaxSizeMB int    `mapstructure:"max_size_mb"`
	MaxCount  int    `mapstructure:"max_count"`
}

// LogConfig 日志配置。
type LogConfig struct {
	Level string `mapstructure:"level"`
	File  string `mapstructure:"file"`
}

// envBindings 显式绑定环境变量，避免 viper AutomaticEnv 在 Unmarshal 时的已知缺陷。
var envBindings = map[string]string{
	"app.name":                           "YICHU_APP_NAME",
	"app.hide_price":                     "YICHU_APP_HIDE_PRICE",
	"app.price_toggle_require_password":  "YICHU_APP_PRICE_TOGGLE_REQUIRE_PASSWORD",
	"server.addr":                        "YICHU_SERVER_ADDR",
	"server.mode":                        "YICHU_SERVER_MODE",
	"server.tls.enabled":                 "YICHU_SERVER_TLS_ENABLED",
	"server.tls.cert_file":               "YICHU_SERVER_TLS_CERT_FILE",
	"server.tls.key_file":                "YICHU_SERVER_TLS_KEY_FILE",
	"database.driver":                    "YICHU_DATABASE_DRIVER",
	"database.dsn":                       "YICHU_DATABASE_DSN",
	"database.max_open_conns":            "YICHU_DATABASE_MAX_OPEN_CONNS",
	"database.max_idle_conns":            "YICHU_DATABASE_MAX_IDLE_CONNS",
	"database.conn_max_lifetime_minutes": "YICHU_DATABASE_CONN_MAX_LIFETIME_MINUTES",
	"database.auto_migrate":              "YICHU_DATABASE_AUTO_MIGRATE",
	"auth.jwt_secret":                    "YICHU_AUTH_JWT_SECRET",
	"auth.token_ttl_hours":               "YICHU_AUTH_TOKEN_TTL_HOURS",
	"auth.init_username":                 "YICHU_AUTH_INIT_USERNAME",
	"auth.init_password":                 "YICHU_AUTH_INIT_PASSWORD",
	"storage.upload_dir":                 "YICHU_STORAGE_UPLOAD_DIR",
	"storage.max_size_mb":                "YICHU_STORAGE_MAX_SIZE_MB",
	"storage.max_count":                  "YICHU_STORAGE_MAX_COUNT",
	"log.level":                          "YICHU_LOG_LEVEL",
	"log.file":                           "YICHU_LOG_FILE",
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("app.name", "衣橱")
	v.SetDefault("app.hide_price", false)
	v.SetDefault("app.price_toggle_require_password", false)

	v.SetDefault("server.addr", ":8080")
	v.SetDefault("server.mode", "release")
	v.SetDefault("server.tls.enabled", false)
	v.SetDefault("server.tls.cert_file", "")
	v.SetDefault("server.tls.key_file", "")

	v.SetDefault("database.driver", "mysql")
	v.SetDefault("database.dsn", "yichu:yichu@tcp(127.0.0.1:3306)/yichu?charset=utf8mb4&parseTime=True&loc=Local&timeout=5s&readTimeout=30s&writeTimeout=30s")
	v.SetDefault("database.max_open_conns", 25)
	v.SetDefault("database.max_idle_conns", 10)
	v.SetDefault("database.conn_max_lifetime_minutes", 30)
	v.SetDefault("database.auto_migrate", true)

	v.SetDefault("auth.jwt_secret", "change-me-in-production")
	v.SetDefault("auth.token_ttl_hours", 168)
	v.SetDefault("auth.init_username", "admin")
	v.SetDefault("auth.init_password", "yichu123")

	v.SetDefault("storage.upload_dir", "./uploads")
	v.SetDefault("storage.max_size_mb", 10)
	v.SetDefault("storage.max_count", 5)

	v.SetDefault("log.level", "info")
	v.SetDefault("log.file", "./logs/yichu.log")
}

// Load 读取默认路径的配置。
func Load() *Config { return LoadFrom("") }

// LoadFrom 读取指定配置文件；path 为空时依次尝试 ./config.yaml、~/.yichu、/opt/yichu。
func LoadFrom(path string) *Config {
	v := viper.New()
	v.SetConfigType("yaml")
	if path != "" {
		v.SetConfigFile(path)
	} else {
		v.SetConfigName("config")
		v.AddConfigPath(".")
		v.AddConfigPath("$HOME/.yichu")
		v.AddConfigPath("/opt/yichu")
	}
	setDefaults(v)

	v.SetEnvPrefix("YICHU")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()
	for key, env := range envBindings {
		if err := v.BindEnv(key, env); err != nil {
			log.Printf("[warn] 绑定环境变量 %s 失败: %v", env, err)
		}
	}

	if err := v.ReadInConfig(); err != nil {
		var notFound viper.ConfigFileNotFoundError
		if !errors.As(err, &notFound) && !os.IsNotExist(err) {
			log.Fatalf("读取配置文件失败: %v", err)
		}
	}

	cfg := &Config{}
	if err := v.Unmarshal(cfg); err != nil {
		log.Fatalf("解析配置失败: %v", err)
	}
	cfg.normalize()
	if err := cfg.validate(); err != nil {
		log.Fatalf("配置校验失败: %v", err)
	}
	return cfg
}

func (c *Config) normalize() {
	c.App.Name = strings.TrimSpace(c.App.Name)
	if c.App.Name == "" {
		c.App.Name = "衣橱"
	}
	c.Server.Addr = strings.TrimSpace(c.Server.Addr)
	if c.Server.Addr == "" {
		c.Server.Addr = ":8080"
	}
	if strings.ToLower(c.Server.Mode) == "debug" {
		c.Server.Mode = "debug"
	} else {
		c.Server.Mode = "release"
	}

	if strings.TrimSpace(c.Auth.JWTSecret) == "" {
		c.Auth.JWTSecret = "change-me-in-production"
	}
	if c.Auth.TokenTTLHours <= 0 {
		c.Auth.TokenTTLHours = 168
	}
	if strings.TrimSpace(c.Auth.InitUsername) == "" {
		c.Auth.InitUsername = "admin"
	}
	if strings.TrimSpace(c.Auth.InitPassword) == "" {
		c.Auth.InitPassword = "yichu123"
	}

	if c.Database.MaxOpenConns <= 0 {
		c.Database.MaxOpenConns = 25
	}
	if c.Database.MaxIdleConns <= 0 {
		c.Database.MaxIdleConns = 10
	}
	if c.Database.ConnMaxLifetimeMinutes <= 0 {
		c.Database.ConnMaxLifetimeMinutes = 30
	}
	// 支持 DSN 中写 ${DB_PASSWORD}，由环境变量填充。
	if strings.Contains(c.Database.DSN, "${") {
		c.Database.DSN = os.Expand(c.Database.DSN, func(key string) string { return os.Getenv(key) })
	}

	if strings.TrimSpace(c.Storage.UploadDir) == "" {
		c.Storage.UploadDir = "./uploads"
	}
	if c.Storage.MaxSizeMB <= 0 {
		c.Storage.MaxSizeMB = 10
	}
	if c.Storage.MaxCount <= 0 {
		c.Storage.MaxCount = 5
	}

	if strings.TrimSpace(c.Log.Level) == "" {
		c.Log.Level = "info"
	}
	if strings.TrimSpace(c.Log.File) == "" {
		c.Log.File = "./logs/yichu.log"
	}
	c.Server.TLS.CertFile = strings.TrimSpace(c.Server.TLS.CertFile)
	c.Server.TLS.KeyFile = strings.TrimSpace(c.Server.TLS.KeyFile)
}

func (c *Config) validate() error {
	if !c.Server.TLS.Enabled {
		return nil
	}
	if c.Server.TLS.CertFile == "" || c.Server.TLS.KeyFile == "" {
		return errors.New("server.tls.enabled=true 时必须配置 cert_file 与 key_file")
	}
	for _, f := range []string{c.Server.TLS.CertFile, c.Server.TLS.KeyFile} {
		if _, err := os.Stat(f); err != nil {
			return fmt.Errorf("TLS 文件不可用 %s: %w", f, err)
		}
	}
	return nil
}

// IsRelease 是否生产模式。
func (c *Config) IsRelease() bool { return c.Server.Mode == "release" }

// MaxUploadBytes 单张图片最大字节数。
func (c *Config) MaxUploadBytes() int64 { return int64(c.Storage.MaxSizeMB) * 1024 * 1024 }
