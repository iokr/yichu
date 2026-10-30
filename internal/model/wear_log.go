package model

import "time"

// WearLog 穿着记录。
type WearLog struct {
	ID        uint64    `gorm:"primaryKey"`
	ItemID    uint64    `gorm:"not null;index:idx_wear_item"`
	WornDate  time.Time `gorm:"type:date;not null;index:idx_worn_date"`
	Occasion  *string   `gorm:"type:varchar(20)"`
	Remark    *string   `gorm:"type:varchar(255)"`
	CreatedAt time.Time `gorm:"type:datetime(3);not null"`
}

// TableName 表名。
func (WearLog) TableName() string { return "wear_logs" }
