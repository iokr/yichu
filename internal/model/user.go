package model

import "time"

// User 单账号管理员。
type User struct {
	ID           uint64    `gorm:"primaryKey"`
	Username     string    `gorm:"type:varchar(64);not null;uniqueIndex:uk_username"`
	PasswordHash string    `gorm:"type:varchar(100);not null"`
	Nickname     string    `gorm:"type:varchar(50);not null;default:''"`
	CreatedAt    time.Time `gorm:"type:datetime(3);not null"`
}

// TableName 表名。
func (User) TableName() string { return "users" }
