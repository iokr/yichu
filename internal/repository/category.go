package repository

import (
	"errors"

	"gorm.io/gorm"

	"github.com/iokr/yichu/internal/model"
)

// SortItem 分类排序项。
type SortItem struct {
	ID   uint64
	Sort int
}

// CategoryRepository 分类数据访问。
type CategoryRepository struct{ db *gorm.DB }

// NewCategoryRepository 构造仓储。
func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

// All 返回全部分类，按父级、排序号排列。
func (r *CategoryRepository) All() ([]model.Category, error) {
	list := make([]model.Category, 0, 32)
	err := r.db.Order("parent_id asc, sort asc, id asc").Find(&list).Error
	return list, err
}

// Get 按 ID 查询，未找到返回 (nil, nil)。
func (r *CategoryRepository) Get(id uint64) (*model.Category, error) {
	var c model.Category
	err := r.db.First(&c, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &c, nil
}

// Create 新增分类。
func (r *CategoryRepository) Create(c *model.Category) error { return r.db.Create(c).Error }

// Update 更新名称、父级与排序。
func (r *CategoryRepository) Update(c *model.Category) error {
	return r.db.Model(&model.Category{}).Where("id = ?", c.ID).Updates(map[string]any{
		"name":      c.Name,
		"parent_id": c.ParentID,
		"sort":      c.Sort,
	}).Error
}

// Delete 物理删除分类（分类无软删除）。
func (r *CategoryRepository) Delete(id uint64) error {
	return r.db.Delete(&model.Category{}, id).Error
}

// NameExists 同层重名校验。
func (r *CategoryRepository) NameExists(parentID uint64, name string, excludeID uint64) (bool, error) {
	q := r.db.Model(&model.Category{}).Where("parent_id = ? AND name = ?", parentID, name)
	if excludeID > 0 {
		q = q.Where("id <> ?", excludeID)
	}
	var count int64
	if err := q.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// ItemCount 统计该分类（一级或二级）关联的衣物数量。
func (r *CategoryRepository) ItemCount(id uint64) (int64, error) {
	var count int64
	err := r.db.Model(&model.Item{}).
		Where("category_id = ? OR sub_category_id = ?", id, id).
		Count(&count).Error
	return count, err
}

// ItemCounts 一级分类衣物数。
func (r *CategoryRepository) ItemCounts() (map[uint64]int64, error) {
	var rows []struct {
		CategoryID uint64
		Count      int64
	}
	if err := r.db.Model(&model.Item{}).
		Select("category_id, COUNT(*) AS count").
		Group("category_id").
		Scan(&rows).Error; err != nil {
		return nil, err
	}
	out := make(map[uint64]int64, len(rows))
	for _, row := range rows {
		out[row.CategoryID] = row.Count
	}
	return out, nil
}

// SubItemCounts 二级分类衣物数。
func (r *CategoryRepository) SubItemCounts() (map[uint64]int64, error) {
	var rows []struct {
		SubCategoryID uint64
		Count         int64
	}
	if err := r.db.Model(&model.Item{}).
		Select("sub_category_id, COUNT(*) AS count").
		Where("sub_category_id IS NOT NULL").
		Group("sub_category_id").
		Scan(&rows).Error; err != nil {
		return nil, err
	}
	out := make(map[uint64]int64, len(rows))
	for _, row := range rows {
		out[row.SubCategoryID] = row.Count
	}
	return out, nil
}

// UpdateSorts 批量更新排序号。
func (r *CategoryRepository) UpdateSorts(items []SortItem) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, item := range items {
			if err := tx.Model(&model.Category{}).Where("id = ?", item.ID).
				Update("sort", item.Sort).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
