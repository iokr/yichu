package service

import "fmt"

// PriceMasker 负责把真实金额转换为对外契约字段。
// 脱敏时真实金额返回 nil，展示值固定为 "¥ ***"。
type PriceMasker struct {
	hidden bool
}

// NewPriceMasker 构造脱敏器。
func NewPriceMasker(hidden bool) PriceMasker { return PriceMasker{hidden: hidden} }

// Hidden 是否处于脱敏态。
func (m PriceMasker) Hidden() bool { return m.hidden }

// Price 返回可下发的价格与展示值。
func (m PriceMasker) Price(v *float64) (*float64, string) {
	if v == nil {
		return nil, ""
	}
	if m.hidden {
		return nil, "¥ ***"
	}
	return v, fmt.Sprintf("¥%.2f", *v)
}

// Amount 统计金额：脱敏态返回 (nil, "***")。
func (m PriceMasker) Amount(v *float64) (*float64, string) {
	if m.hidden {
		return nil, "***"
	}
	if v == nil {
		return nil, "¥0.00"
	}
	return v, fmt.Sprintf("¥%.2f", *v)
}

// Cost 单次成本展示。
func (m PriceMasker) Cost(v *float64) (*float64, string) {
	if v == nil {
		return nil, ""
	}
	if m.hidden {
		return nil, "¥ ***"
	}
	return v, fmt.Sprintf("¥%.2f", *v)
}

// CSV 导出的价格列。
func (m PriceMasker) CSV(v *float64) string {
	if v == nil {
		return ""
	}
	if m.hidden {
		return "***"
	}
	return fmt.Sprintf("%.2f", *v)
}
