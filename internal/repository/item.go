package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"

	"github.com/iokr/yichu/internal/model"
)

// ItemQuery 衣物列表查询条件（已由 service 归一化）。
type ItemQuery struct {
	CategoryID    uint64
	SubCategoryID uint64
	TagID         uint64
	Status        string
	Keyword       string
	Brand         string
	Colors        []string
	Seasons       []string
	Occasions     []string
	MinPrice      *float64
	MaxPrice      *float64
	Year          *int
	Sort          string
	Order         string
	Page          int
	Size          int
}

// CodeRow 季节/颜色统计所需的原始行。
type CodeRow struct {
	Season model.Codes
	Color  model.Codes
}

// ItemRepository 衣物数据访问。
type ItemRepository struct{ db *gorm.DB }

// NewItemRepository 构造仓储。
func NewItemRepository(db *gorm.DB) *ItemRepository { return &ItemRepository{db: db} }

// 允许的排序字段白名单，避免 SQL 注入。
var itemSortColumns = map[string]string{
	"created_at":     "created_at",
	"updated_at":     "updated_at",
	"purchase_date":  "purchase_date",
	"purchase_price": "purchase_price",
	"wear_count":     "wear_count",
	"last_worn_at":   "last_worn_at",
	"name":           "name",
}

func orderBy(sort, order string) string {
	col, ok := itemSortColumns[sort]
	if !ok {
		col = "created_at"
	}
	if strings.EqualFold(order, "asc") {
		return fmt.Sprintf("%s asc, id asc", col)
	}
	return fmt.Sprintf("%s desc, id desc", col)
}

func likePattern(keyword string) string {
	replacer := strings.NewReplacer("\\", "\\\\", "%", "\\%", "_", "\\_")
	return "%" + replacer.Replace(keyword) + "%"
}

func (r *ItemRepository) applyFilters(q *gorm.DB, f ItemQuery) *gorm.DB {
	if f.CategoryID > 0 {
		q = q.Where("category_id = ?", f.CategoryID)
	}
	if f.SubCategoryID > 0 {
		q = q.Where("sub_category_id = ?", f.SubCategoryID)
	}
	if f.TagID > 0 {
		q = q.Where("EXISTS (SELECT 1 FROM item_tags WHERE item_tags.item_id = items.id AND item_tags.tag_id = ?)", f.TagID)
	}
	if f.Status != "" {
		q = q.Where("status = ?", f.Status)
	}
	if f.Keyword != "" {
		pattern := likePattern(f.Keyword)
		q = q.Where("(name LIKE ? OR brand LIKE ?)", pattern, pattern)
	}
	if f.Brand != "" {
		q = q.Where("brand = ?", f.Brand)
	}
	for _, code := range f.Colors {
		q = q.Where("JSON_CONTAINS(color, JSON_QUOTE(?))", code)
	}
	for _, code := range f.Seasons {
		q = q.Where("JSON_CONTAINS(season, JSON_QUOTE(?))", code)
	}
	for _, code := range f.Occasions {
		q = q.Where("JSON_CONTAINS(occasion, JSON_QUOTE(?))", code)
	}
	if f.MinPrice != nil {
		q = q.Where("purchase_price >= ?", *f.MinPrice)
	}
	if f.MaxPrice != nil {
		q = q.Where("purchase_price <= ?", *f.MaxPrice)
	}
	if f.Year != nil {
		start := time.Date(*f.Year, 1, 1, 0, 0, 0, 0, time.Local)
		q = q.Where("purchase_date >= ? AND purchase_date < ?", start, start.AddDate(1, 0, 0))
	}
	return q
}

func (r *ItemRepository) preload(q *gorm.DB) *gorm.DB {
	return q.
		Preload("Images", func(db *gorm.DB) *gorm.DB { return db.Order("sort asc, id asc") }).
		Preload("Tags", func(db *gorm.DB) *gorm.DB { return db.Order("id asc") })
}

// List 分页查询。
func (r *ItemRepository) List(f ItemQuery) ([]model.Item, int64, error) {
	if f.Page <= 0 {
		f.Page = 1
	}
	if f.Size <= 0 {
		f.Size = 20
	}
	var total int64
	if err := r.applyFilters(r.db.Model(&model.Item{}), f).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	items := make([]model.Item, 0, f.Size)
	err := r.preload(r.applyFilters(r.db.Model(&model.Item{}), f)).
		Order(orderBy(f.Sort, f.Order)).
		Offset((f.Page - 1) * f.Size).
		Limit(f.Size).
		Find(&items).Error
	if err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

// All 不分页查询，用于导出/备份。
func (r *ItemRepository) All(f ItemQuery) ([]model.Item, error) {
	items := make([]model.Item, 0, 256)
	err := r.preload(r.applyFilters(r.db.Model(&model.Item{}), f)).
		Order("id asc").
		Find(&items).Error
	return items, err
}

// Get 查询详情，未找到返回 (nil, nil)。
func (r *ItemRepository) Get(id uint64) (*model.Item, error) {
	var item model.Item
	err := r.preload(r.db).First(&item, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &item, nil
}

// Create 新增衣物并写入标签。
func (r *ItemRepository) Create(item *model.Item, tagNames []string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Omit("Tags").Create(item).Error; err != nil {
			return err
		}
		return replaceTags(tx, item.ID, tagNames)
	})
}

// Update 全量更新衣物，替换图片与标签。
func (r *ItemRepository) Update(id uint64, fields map[string]any, images []model.ItemImage, tagNames []string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.Item{}).Where("id = ?", id).Updates(fields).Error; err != nil {
			return err
		}
		if err := tx.Where("item_id = ?", id).Delete(&model.ItemImage{}).Error; err != nil {
			return err
		}
		if len(images) > 0 {
			rows := make([]model.ItemImage, 0, len(images))
			now := time.Now()
			for i := range images {
				rows = append(rows, model.ItemImage{
					ItemID:    id,
					URL:       images[i].URL,
					ThumbURL:  images[i].ThumbURL,
					Sort:      i,
					CreatedAt: now,
				})
			}
			if err := tx.Create(&rows).Error; err != nil {
				return err
			}
		}
		return replaceTags(tx, id, tagNames)
	})
}

func replaceTags(tx *gorm.DB, itemID uint64, tagNames []string) error {
	tags := make([]model.Tag, 0, len(tagNames))
	seen := make(map[string]struct{}, len(tagNames))
	for _, raw := range tagNames {
		name := strings.TrimSpace(raw)
		if name == "" {
			continue
		}
		if _, dup := seen[name]; dup {
			continue
		}
		seen[name] = struct{}{}

		var tag model.Tag
		err := tx.Where("name = ?", name).First(&tag).Error
		switch {
		case err == nil:
		case errors.Is(err, gorm.ErrRecordNotFound):
			tag = model.Tag{Name: name}
			if cerr := tx.Create(&tag).Error; cerr != nil {
				if !errors.Is(cerr, gorm.ErrDuplicatedKey) {
					return cerr
				}
				if ferr := tx.Where("name = ?", name).First(&tag).Error; ferr != nil {
					return ferr
				}
			}
		default:
			return err
		}
		tags = append(tags, tag)
	}
	return tx.Model(&model.Item{ID: itemID}).Association("Tags").Replace(tags)
}

// Delete 软删除衣物。
func (r *ItemRepository) Delete(id uint64) error {
	return r.db.Delete(&model.Item{}, id).Error
}

// BatchDelete 批量软删除。
func (r *ItemRepository) BatchDelete(ids []uint64) error {
	return r.db.Where("id IN ?", ids).Delete(&model.Item{}).Error
}

// BatchStatus 批量修改状态。
func (r *ItemRepository) BatchStatus(ids []uint64, status string) error {
	return r.db.Model(&model.Item{}).Where("id IN ?", ids).Updates(map[string]any{"status": status}).Error
}

// Count 衣物总数。
func (r *ItemRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&model.Item{}).Count(&count).Error
	return count, err
}

// CountByStatus 按状态统计。
func (r *ItemRepository) CountByStatus() (map[string]int64, error) {
	var rows []struct {
		Status string
		Count  int64
	}
	if err := r.db.Model(&model.Item{}).
		Select("status, COUNT(*) AS count").
		Group("status").
		Scan(&rows).Error; err != nil {
		return nil, err
	}
	out := make(map[string]int64, len(rows))
	for _, row := range rows {
		out[row.Status] = row.Count
	}
	return out, nil
}

// SumPrice 汇总购买金额；year 为 nil 时统计全部。
func (r *ItemRepository) SumPrice(year *int) (*float64, error) {
	q := r.db.Model(&model.Item{})
	if year != nil {
		start := time.Date(*year, 1, 1, 0, 0, 0, 0, time.Local)
		q = q.Where("purchase_date >= ? AND purchase_date < ?", start, start.AddDate(1, 0, 0))
	}
	var sum sql.NullFloat64
	if err := q.Select("COALESCE(SUM(purchase_price), 0)").Scan(&sum).Error; err != nil {
		return nil, err
	}
	if !sum.Valid {
		return nil, nil
	}
	v := sum.Float64
	return &v, nil
}

// PricedCount 有价格的衣物数量（用于均价）。
func (r *ItemRepository) PricedCount() (int64, error) {
	var count int64
	err := r.db.Model(&model.Item{}).Where("purchase_price IS NOT NULL").Count(&count).Error
	return count, err
}

// SumWearCount 总穿着次数。
func (r *ItemRepository) SumWearCount() (int64, error) {
	var sum sql.NullInt64
	if err := r.db.Model(&model.Item{}).Select("COALESCE(SUM(wear_count), 0)").Scan(&sum).Error; err != nil {
		return 0, err
	}
	if !sum.Valid {
		return 0, nil
	}
	return sum.Int64, nil
}

// Recent 最近添加的衣物。
func (r *ItemRepository) Recent(limit int) ([]model.Item, error) {
	items := make([]model.Item, 0, limit)
	err := r.preload(r.db.Model(&model.Item{})).
		Order("created_at desc, id desc").
		Limit(limit).
		Find(&items).Error
	return items, err
}

// AllCodeRows 返回全部用于统计的 JSON 字段。
func (r *ItemRepository) AllCodeRows() ([]CodeRow, error) {
	rows := make([]CodeRow, 0, 256)
	err := r.db.Model(&model.Item{}).Select("season, color").Find(&rows).Error
	return rows, err
}

// CostPerWear 单次成本排行（价格非空且穿着次数大于 0）。
func (r *ItemRepository) CostPerWear(limit int) ([]model.Item, error) {
	items := make([]model.Item, 0, limit)
	err := r.preload(r.db.Model(&model.Item{})).
		Where("wear_count > 0 AND purchase_price IS NOT NULL").
		Order("purchase_price / wear_count asc, id asc").
		Limit(limit).
		Find(&items).Error
	return items, err
}
