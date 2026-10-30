package service

import (
	"strconv"
	"sync"

	"github.com/iokr/yichu/internal/config"
	"github.com/iokr/yichu/internal/model"
	"github.com/iokr/yichu/internal/repository"
)

// SettingsService 运行时设置，优先级：settings 表 > config.yaml > 默认值。
type SettingsService struct {
	repo            *repository.SettingRepository
	defaultHidden   bool
	requirePassword bool

	mu     sync.RWMutex
	hidden bool
	loaded bool
}

// NewSettingsService 构造服务。
func NewSettingsService(repo *repository.SettingRepository, cfg *config.Config) *SettingsService {
	return &SettingsService{
		repo:            repo,
		defaultHidden:   cfg.App.HidePrice,
		requirePassword: cfg.App.PriceToggleRequirePassword,
	}
}

// Init 启动时加载设置，表内无记录则沿用配置文件默认值。
func (s *SettingsService) Init() error {
	raw, ok, err := s.repo.Get(model.SettingHidePrice)
	if err != nil {
		return err
	}
	hidden := s.defaultHidden
	if ok {
		if v, perr := strconv.ParseBool(raw); perr == nil {
			hidden = v
		}
	}
	s.mu.Lock()
	s.hidden = hidden
	s.loaded = true
	s.mu.Unlock()
	return nil
}

// HidePrice 当前是否隐藏价格。
func (s *SettingsService) HidePrice() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if !s.loaded {
		return s.defaultHidden
	}
	return s.hidden
}

// SetHidePrice 切换价格可见性并持久化。
func (s *SettingsService) SetHidePrice(hidden bool) error {
	if err := s.repo.Set(model.SettingHidePrice, strconv.FormatBool(hidden)); err != nil {
		return err
	}
	s.mu.Lock()
	s.hidden = hidden
	s.loaded = true
	s.mu.Unlock()
	return nil
}

// RequirePassword 切换价格可见性是否需要密码（V1 默认关闭）。
func (s *SettingsService) RequirePassword() bool { return s.requirePassword }

// Masker 生成当前脱敏器。
func (s *SettingsService) Masker() PriceMasker { return NewPriceMasker(s.HidePrice()) }
