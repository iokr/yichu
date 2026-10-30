package repository

import (
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/iokr/yichu/internal/model"
)

// SettingRepository settings 表数据访问。
type SettingRepository struct{ db *gorm.DB }

// NewSettingRepository 构造仓储。
func NewSettingRepository(db *gorm.DB) *SettingRepository { return &SettingRepository{db: db} }

// Get 读取设置项，第二个返回值表示是否存在。
func (r *SettingRepository) Get(key string) (string, bool, error) {
	var s model.Setting
	err := r.db.Where("`key` = ?", key).First(&s).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}
	return s.Value, true, nil
}

// Set 写入或更新设置项。
func (r *SettingRepository) Set(key, value string) error {
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "key"}},
		DoUpdates: clause.Assignments(map[string]any{"value": value}),
	}).Create(&model.Setting{Key: key, Value: value}).Error
}

// All 返回全部设置项。
func (r *SettingRepository) All() ([]model.Setting, error) {
	list := make([]model.Setting, 0, 8)
	err := r.db.Order("`key` asc").Find(&list).Error
	return list, err
}
