package repository

import (
	"database/sql"

	"gorm.io/gorm"
)

// GroupRow 分组统计行。
type GroupRow struct {
	Key    string          `gorm:"column:key"`
	Label  string          `gorm:"column:label"`
	Count  int64           `gorm:"column:count"`
	Amount sql.NullFloat64 `gorm:"column:amount"`
}

// StatsRepository 统计报表数据访问。
type StatsRepository struct{ db *gorm.DB }

// NewStatsRepository 构造仓储。
func NewStatsRepository(db *gorm.DB) *StatsRepository { return &StatsRepository{db: db} }

func (r *StatsRepository) group(selectClause, joinClause, whereClause, groupClause, orderClause string, args ...any) ([]GroupRow, error) {
	q := r.db.Table("items").
		Select(selectClause).
		Where("items.deleted_at IS NULL")
	if joinClause != "" {
		q = q.Joins(joinClause)
	}
	if whereClause != "" {
		q = q.Where(whereClause, args...)
	}
	q = q.Group(groupClause)
	if orderClause != "" {
		q = q.Order(orderClause)
	}
	rows := make([]GroupRow, 0, 16)
	if err := q.Scan(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

// ByCategory 按一级分类统计。
func (r *StatsRepository) ByCategory() ([]GroupRow, error) {
	return r.group(
		"CAST(items.category_id AS CHAR) AS `key`, COALESCE(categories.name, '未分类') AS label, COUNT(*) AS count, SUM(items.purchase_price) AS amount",
		"LEFT JOIN categories ON categories.id = items.category_id",
		"", "items.category_id, categories.name", "`count` DESC",
	)
}

// ByBrand 按品牌统计。
func (r *StatsRepository) ByBrand() ([]GroupRow, error) {
	return r.group(
		"items.brand AS `key`, items.brand AS label, COUNT(*) AS count, SUM(items.purchase_price) AS amount",
		"", "items.brand IS NOT NULL AND items.brand <> ''", "items.brand", "`count` DESC",
	)
}

// ByYear 按购买年份统计。
// 年份先落到派生表再分组：MySQL 8 的 ONLY_FULL_GROUP_BY 不认 GROUP BY 表达式上套 CAST/CONCAT 的写法。
func (r *StatsRepository) ByYear() ([]GroupRow, error) {
	const query = "SELECT CAST(t.purchase_year AS CHAR) AS `key`, CONCAT(t.purchase_year, ' 年') AS label, COUNT(*) AS count, SUM(t.purchase_price) AS amount" +
		" FROM (SELECT YEAR(items.purchase_date) AS purchase_year, items.purchase_price AS purchase_price FROM items" +
		" WHERE items.deleted_at IS NULL AND items.purchase_date IS NOT NULL) AS t" +
		" GROUP BY t.purchase_year ORDER BY t.purchase_year DESC"
	rows := make([]GroupRow, 0, 8)
	if err := r.db.Raw(query).Scan(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

// ByStatus 按状态统计。
func (r *StatsRepository) ByStatus() ([]GroupRow, error) {
	return r.group(
		"items.status AS `key`, items.status AS label, COUNT(*) AS count, SUM(items.purchase_price) AS amount",
		"", "", "items.status", "`count` DESC",
	)
}

// TopBrands 品牌数量排行（用于仪表盘/统计页）。
func (r *StatsRepository) TopBrands(limit int) ([]GroupRow, error) {
	rows, err := r.ByBrand()
	if err != nil {
		return nil, err
	}
	if limit > 0 && len(rows) > limit {
		return rows[:limit], nil
	}
	return rows, nil
}
