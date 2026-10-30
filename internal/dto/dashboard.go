package dto

import (
	"github.com/iokr/yichu/internal/pkg/dict"
)

// CodeStat 仅数量型统计（季节/颜色/状态）。
type CodeStat struct {
	Code    string  `json:"code"`
	Label   string  `json:"label"`
	Hex     string  `json:"hex,omitempty"`
	Count   int64   `json:"count"`
	Percent float64 `json:"percent"`
}

// KeyStat 带金额的统计（分类/品牌/年份）。
type KeyStat struct {
	Key           string   `json:"key"`
	Label         string   `json:"label"`
	Count         int64    `json:"count"`
	Percent       float64  `json:"percent"`
	Amount        *float64 `json:"amount"`
	AmountDisplay string   `json:"amount_display"`
}

// DashboardResponse 仪表盘数据。
type DashboardResponse struct {
	TotalCount         int64          `json:"total_count"`
	WearingCount       int64          `json:"wearing_count"`
	IdleCount          int64          `json:"idle_count"`
	TotalWearCount     int64          `json:"total_wear_count"`
	CategoryStats      []KeyStat      `json:"category_stats"`
	StatusStats        []CodeStat     `json:"status_stats"`
	Year               int            `json:"year"`
	YearAmount         *float64       `json:"year_amount"`
	YearAmountDisplay  string         `json:"year_amount_display"`
	TotalAmount        *float64       `json:"total_amount"`
	TotalAmountDisplay string         `json:"total_amount_display"`
	PriceMasked        bool           `json:"price_masked"`
	Recent             []ItemListItem `json:"recent"`
}

// StatsOverview 统计概览。
type StatsOverview struct {
	TotalCount         int64    `json:"total_count"`
	WearingCount       int64    `json:"wearing_count"`
	IdleCount          int64    `json:"idle_count"`
	TotalWearCount     int64    `json:"total_wear_count"`
	TotalAmount        *float64 `json:"total_amount"`
	TotalAmountDisplay string   `json:"total_amount_display"`
	AvgPrice           *float64 `json:"avg_price"`
	AvgPriceDisplay    string   `json:"avg_price_display"`
	PriceMasked        bool     `json:"price_masked"`
}

// CostPerWearItem 单次成本排行项。
type CostPerWearItem struct {
	ID            uint64   `json:"id"`
	Name          string   `json:"name"`
	Cover         string   `json:"cover"`
	WearCount     uint32   `json:"wear_count"`
	PurchasePrice *float64 `json:"purchase_price"`
	PriceDisplay  string   `json:"purchase_price_display"`
	CostPerWear   *float64 `json:"cost_per_wear"`
	CostDisplay   string   `json:"cost_per_wear_display"`
	PriceMasked   bool     `json:"price_masked"`
}

// DictOption 别名，避免前端类型推断困难。
type DictOption = dict.Option
