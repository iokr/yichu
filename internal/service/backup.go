package service

import (
	"context"
	"encoding/json"
	"time"

	"github.com/iokr/yichu/internal/dto"
	"github.com/iokr/yichu/internal/pkg/errcode"
	"github.com/iokr/yichu/internal/repository"
)

// BackupService 备份导出。
type BackupService struct {
	appName     string
	version     string
	users       *repository.UserRepository
	categories  *CategoryService
	items       *ItemService
	tags        *TagService
	wears       *WearService
	settings    *SettingsService
	settingRepo *repository.SettingRepository
}

// NewBackupService 构造服务。
func NewBackupService(
	appName, version string,
	users *repository.UserRepository,
	categories *CategoryService,
	items *ItemService,
	tags *TagService,
	wears *WearService,
	settings *SettingsService,
	settingRepo *repository.SettingRepository,
) *BackupService {
	return &BackupService{
		appName:     appName,
		version:     version,
		users:       users,
		categories:  categories,
		items:       items,
		tags:        tags,
		wears:       wears,
		settings:    settings,
		settingRepo: settingRepo,
	}
}

type backupUser struct {
	ID        uint64 `json:"id"`
	Username  string `json:"username"`
	Nickname  string `json:"nickname"`
	CreatedAt string `json:"created_at"`
}

type backupFile struct {
	AppName     string             `json:"app_name"`
	Version     string             `json:"version"`
	CreatedAt   string             `json:"created_at"`
	PriceMasked bool               `json:"price_masked"`
	Users       []backupUser       `json:"users"`
	Categories  []dto.CategoryNode `json:"categories"`
	Items       []dto.ItemDetail   `json:"items"`
	Tags        []dto.TagDTO       `json:"tags"`
	WearLogs    []dto.WearLogItem  `json:"wear_logs"`
	Settings    map[string]string  `json:"settings"`
}

// Export 生成 JSON 备份，受脱敏约束。
func (s *BackupService) Export(ctx context.Context) ([]byte, bool, error) {
	categories, err := s.categories.Tree(ctx)
	if err != nil {
		return nil, false, err
	}
	items, masked, err := s.items.AllDetails(ctx)
	if err != nil {
		return nil, false, err
	}
	tags, err := s.tags.List(ctx)
	if err != nil {
		return nil, false, err
	}
	logs, err := s.wears.All(ctx)
	if err != nil {
		return nil, false, err
	}
	settingsRows, err := s.settingRepo.All()
	if err != nil {
		return nil, false, errcode.Internal(err)
	}
	users, err := s.users.All()
	if err != nil {
		return nil, false, errcode.Internal(err)
	}

	settings := make(map[string]string, len(settingsRows))
	for _, row := range settingsRows {
		settings[row.Key] = row.Value
	}

	backupUsers := make([]backupUser, 0, len(users))
	for _, u := range users {
		backupUsers = append(backupUsers, backupUser{
			ID:        u.ID,
			Username:  u.Username,
			Nickname:  u.Nickname,
			CreatedAt: u.CreatedAt.Format(dateTimeLayout),
		})
	}

	file := backupFile{
		AppName:     s.appName,
		Version:     s.version,
		CreatedAt:   time.Now().Format(dateTimeLayout),
		PriceMasked: masked,
		Users:       backupUsers,
		Categories:  categories,
		Items:       items,
		Tags:        tags,
		WearLogs:    s.wears.MapLogs(logs),
		Settings:    settings,
	}
	data, merr := json.MarshalIndent(file, "", "  ")
	if merr != nil {
		return nil, false, errcode.Internal(merr)
	}
	return data, masked, nil
}
