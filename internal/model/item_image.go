package model

import "time"

// ItemImage 衣物图片，Sort 升序，第一张为封面。
type ItemImage struct {
	ID        uint64    `gorm:"primaryKey"`
	ItemID    uint64    `gorm:"not null;index:idx_item_id"`
	URL       string    `gorm:"type:varchar(500);not null"`
	ThumbURL  string    `gorm:"type:varchar(500);not null"`
	Sort      int       `gorm:"not null;default:0"`
	CreatedAt time.Time `gorm:"type:datetime(3);not null"`
}

// TableName 表名。
func (ItemImage) TableName() string { return "item_images" }

// Tag 自定义标签。
type Tag struct {
	ID        uint64    `gorm:"primaryKey"`
	Name      string    `gorm:"type:varchar(50);not null;uniqueIndex:uk_tag_name"`
	CreatedAt time.Time `gorm:"type:datetime(3);not null"`
}

// TableName 表名。
func (Tag) TableName() string { return "tags" }
