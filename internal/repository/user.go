package repository

import (
	"errors"

	"gorm.io/gorm"

	"github.com/iokr/yichu/internal/model"
)

// UserRepository 用户数据访问。
type UserRepository struct{ db *gorm.DB }

// NewUserRepository 构造仓储。
func NewUserRepository(db *gorm.DB) *UserRepository { return &UserRepository{db: db} }

// ByUsername 按用户名查询，未找到返回 (nil, nil)。
func (r *UserRepository) ByUsername(username string) (*model.User, error) {
	var u model.User
	err := r.db.Where("username = ?", username).First(&u).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// ByID 按 ID 查询，未找到返回 (nil, nil)。
func (r *UserRepository) ByID(id uint64) (*model.User, error) {
	var u model.User
	err := r.db.First(&u, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// All 返回全部用户（备份用，调用方需剔除密码哈希）。
func (r *UserRepository) All() ([]model.User, error) {
	list := make([]model.User, 0, 2)
	err := r.db.Select("id, username, nickname, created_at").Order("id asc").Find(&list).Error
	return list, err
}
