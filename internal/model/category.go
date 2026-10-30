package model

import "time"

// Category 衣物分类，ParentID=0 表示一级分类。
type Category struct {
	ID        uint64    `gorm:"primaryKey"`
	ParentID  uint64    `gorm:"not null;default:0;index:idx_parent"`
	Name      string    `gorm:"type:varchar(50);not null;index:idx_name"`
	Sort      int       `gorm:"not null;default:0"`
	IsSystem  bool      `gorm:"not null;default:false"`
	CreatedAt time.Time `gorm:"type:datetime(3);not null"`
}

// TableName 表名。
func (Category) TableName() string { return "categories" }

// IsTop 是否一级分类。
func (c Category) IsTop() bool { return c.ParentID == 0 }
