package model

// Setting 键值配置，运行时覆盖 config.yaml。
type Setting struct {
	Key   string `gorm:"primaryKey;type:varchar(64)"`
	Value string `gorm:"type:varchar(255);not null;default:''"`
}

// TableName 表名。
func (Setting) TableName() string { return "settings" }

// 设置项 key。
const (
	SettingHidePrice     = "hide_price"
	SettingSchemaVersion = "schema_version"
)

// SchemaVersion 当前数据库结构版本（V1.0 由 AutoMigrate 维护）。
const SchemaVersion = "1.0"
