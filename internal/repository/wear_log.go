package repository

import (
	"gorm.io/gorm"

	"github.com/iokr/yichu/internal/model"
)

// WearLogRepository 穿着记录数据访问。
type WearLogRepository struct{ db *gorm.DB }

// NewWearLogRepository 构造仓储。
func NewWearLogRepository(db *gorm.DB) *WearLogRepository {
	return &WearLogRepository{db: db}
}

// ListByItem 查询单件衣物的穿着记录。
func (r *WearLogRepository) ListByItem(itemID uint64, limit int) ([]model.WearLog, error) {
	list := make([]model.WearLog, 0, 16)
	q := r.db.Where("item_id = ?", itemID).Order("worn_date desc, id desc")
	if limit > 0 {
		q = q.Limit(limit)
	}
	err := q.Find(&list).Error
	return list, err
}

// All 返回全部穿着记录（备份用）。
func (r *WearLogRepository) All() ([]model.WearLog, error) {
	list := make([]model.WearLog, 0, 128)
	err := r.db.Order("id asc").Find(&list).Error
	return list, err
}

// Add 新增穿着记录并同步衣物穿着次数与最近穿着时间。
func (r *WearLogRepository) Add(log *model.WearLog) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(log).Error; err != nil {
			return err
		}
		if err := tx.Model(&model.Item{}).Where("id = ?", log.ItemID).
			UpdateColumn("wear_count", gorm.Expr("wear_count + 1")).Error; err != nil {
			return err
		}
		return tx.Model(&model.Item{}).
			Where("id = ? AND (last_worn_at IS NULL OR last_worn_at < ?)", log.ItemID, log.WornDate).
			UpdateColumn("last_worn_at", log.WornDate).Error
	})
}
