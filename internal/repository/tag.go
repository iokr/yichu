package repository

import (
	"errors"

	"gorm.io/gorm"

	"github.com/iokr/yichu/internal/model"
)

// TagRepository 标签数据访问。
type TagRepository struct{ db *gorm.DB }

// NewTagRepository 构造仓储。
func NewTagRepository(db *gorm.DB) *TagRepository { return &TagRepository{db: db} }

// All 返回全部标签。
func (r *TagRepository) All() ([]model.Tag, error) {
	list := make([]model.Tag, 0, 32)
	err := r.db.Order("id asc").Find(&list).Error
	return list, err
}

// Counts 统计每个标签关联的有效衣物数量。
func (r *TagRepository) Counts() (map[uint64]int64, error) {
	var rows []struct {
		TagID uint64
		Count int64
	}
	err := r.db.Table("item_tags").
		Select("item_tags.tag_id AS tag_id, COUNT(*) AS count").
		Joins("JOIN items ON items.id = item_tags.item_id AND items.deleted_at IS NULL").
		Group("item_tags.tag_id").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	out := make(map[uint64]int64, len(rows))
	for _, row := range rows {
		out[row.TagID] = row.Count
	}
	return out, nil
}

// GetByName 按名称查询标签，未找到返回 (nil, nil)。
func (r *TagRepository) GetByName(name string) (*model.Tag, error) {
	var tag model.Tag
	err := r.db.Where("name = ?", name).First(&tag).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &tag, nil
}
