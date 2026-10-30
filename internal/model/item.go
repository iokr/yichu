package model

import (
	"time"

	"gorm.io/gorm"
)

// Item 衣物主表。
type Item struct {
	ID              uint64         `gorm:"primaryKey"`
	Name            string         `gorm:"type:varchar(100);not null;index:idx_name"`
	CategoryID      uint64         `gorm:"not null;index:idx_category"`
	SubCategoryID   *uint64        `gorm:"index:idx_sub_category"`
	Brand           *string        `gorm:"type:varchar(50);index:idx_brand"`
	Color           Codes          `gorm:"type:JSON"`
	Size            *string        `gorm:"type:varchar(20)"`
	Season          Codes          `gorm:"type:JSON"`
	Occasion        Codes          `gorm:"type:JSON"`
	PurchaseDate    *time.Time     `gorm:"type:date;index:idx_purchase_date"`
	PurchasePrice   *float64       `gorm:"type:decimal(10,2)"`
	PurchaseChannel *string        `gorm:"type:varchar(50)"`
	PurchaseURL     *string        `gorm:"type:varchar(500)"`
	Status          string         `gorm:"type:varchar(16);not null;default:wearing;index:idx_status"`
	StorageLocation *string        `gorm:"type:varchar(50)"`
	Remark          *string        `gorm:"type:varchar(500)"`
	WearCount       uint32         `gorm:"type:int unsigned;not null;default:0"`
	LastWornAt      *time.Time     `gorm:"type:datetime(3)"`
	CreatedAt       time.Time      `gorm:"type:datetime(3);not null;index:idx_created_at"`
	UpdatedAt       time.Time      `gorm:"type:datetime(3);not null"`
	DeletedAt       gorm.DeletedAt `gorm:"type:datetime(3);index:idx_deleted_at"`

	Images []ItemImage `gorm:"foreignKey:ItemID"`
	Tags   []Tag       `gorm:"many2many:item_tags;"`
}

// TableName 表名。
func (Item) TableName() string { return "items" }

// Cover 返回封面缩略图地址，优先第一张图片。
func (i Item) Cover() string {
	if len(i.Images) == 0 {
		return ""
	}
	if i.Images[0].ThumbURL != "" {
		return i.Images[0].ThumbURL
	}
	return i.Images[0].URL
}
