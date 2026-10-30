package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// Codes 是存 MySQL JSON 的字符串数组（颜色/季节/场合 code）。
type Codes []string

// Value 实现 driver.Valuer。
func (c Codes) Value() (driver.Value, error) {
	if len(c) == 0 {
		return nil, nil
	}
	b, err := json.Marshal([]string(c))
	if err != nil {
		return nil, err
	}
	return string(b), nil
}

// Scan 实现 sql.Scanner。
func (c *Codes) Scan(v any) error {
	if v == nil {
		*c = nil
		return nil
	}
	var raw []byte
	switch t := v.(type) {
	case []byte:
		raw = t
	case string:
		raw = []byte(t)
	default:
		return fmt.Errorf("无法解析 JSON 数组字段: %T", v)
	}
	if len(raw) == 0 {
		*c = nil
		return nil
	}
	var out []string
	if err := json.Unmarshal(raw, &out); err != nil {
		return fmt.Errorf("解析 JSON 数组失败: %w", err)
	}
	*c = out
	return nil
}

// GormDataType 告诉 GORM 使用 JSON 列类型。
func (Codes) GormDataType() string { return "json" }

// GormDBDataType 兼容不同数据库的列类型。
func (Codes) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	switch db.Dialector.Name() {
	case "mysql":
		return "JSON"
	case "postgres":
		return "JSONB"
	default:
		return "TEXT"
	}
}
