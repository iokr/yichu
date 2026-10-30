package dto

import (
	"strings"

	"github.com/iokr/yichu/internal/pkg/dict"
)

// ItemImageInput 新增/编辑时提交的图片。
type ItemImageInput struct {
	URL      string `json:"url"`
	ThumbURL string `json:"thumb_url"`
}

// ItemImageDTO 详情中的图片。
type ItemImageDTO struct {
	ID       uint64 `json:"id"`
	URL      string `json:"url"`
	ThumbURL string `json:"thumb_url"`
	Sort     int    `json:"sort"`
}

// TagDTO 标签。
type TagDTO struct {
	ID    uint64 `json:"id"`
	Name  string `json:"name"`
	Count int64  `json:"count,omitempty"`
}

// ItemRequest 新增/更新衣物请求。
type ItemRequest struct {
	Name            string           `json:"name" binding:"required"`
	CategoryID      uint64           `json:"category_id"`
	SubCategoryID   *uint64          `json:"sub_category_id"`
	Brand           string           `json:"brand"`
	Color           []string         `json:"color"`
	Size            string           `json:"size"`
	Season          []string         `json:"season"`
	Occasion        []string         `json:"occasion"`
	PurchaseDate    string           `json:"purchase_date"`
	PurchasePrice   *float64         `json:"purchase_price"`
	PurchaseChannel string           `json:"purchase_channel"`
	PurchaseURL     string           `json:"purchase_url"`
	Status          string           `json:"status"`
	StorageLocation string           `json:"storage_location"`
	Remark          string           `json:"remark"`
	Images          []ItemImageInput `json:"images"`
	Tags            []string         `json:"tags"`
}

// ItemListItem 列表项。
type ItemListItem struct {
	ID                   uint64        `json:"id"`
	Name                 string        `json:"name"`
	CategoryID           uint64        `json:"category_id"`
	Category             string        `json:"category"`
	SubCategoryID        *uint64       `json:"sub_category_id"`
	SubCategory          string        `json:"sub_category"`
	Brand                string        `json:"brand"`
	Color                []dict.Option `json:"color"`
	Size                 string        `json:"size"`
	Season               []dict.Option `json:"season"`
	Status               dict.Option   `json:"status"`
	PurchaseDate         *string       `json:"purchase_date"`
	PurchasePrice        *float64      `json:"purchase_price"`
	PurchasePriceDisplay string        `json:"purchase_price_display"`
	Cover                string        `json:"cover"`
	WearCount            uint32        `json:"wear_count"`
	LastWornAt           *string       `json:"last_worn_at"`
	PriceMasked          bool          `json:"price_masked"`
	CreatedAt            string        `json:"created_at"`
}

// ItemDetail 详情，列表字段内联。
type ItemDetail struct {
	ItemListItem
	Occasion           []dict.Option  `json:"occasion"`
	PurchaseChannel    string         `json:"purchase_channel"`
	PurchaseURL        string         `json:"purchase_url"`
	StorageLocation    string         `json:"storage_location"`
	Remark             string         `json:"remark"`
	Images             []ItemImageDTO `json:"images"`
	Tags               []TagDTO       `json:"tags"`
	CostPerWear        *float64       `json:"cost_per_wear"`
	CostPerWearDisplay string         `json:"cost_per_wear_display"`
	UpdatedAt          string         `json:"updated_at"`
}

// ItemQuery 列表查询参数。
type ItemQuery struct {
	CategoryID    uint64   `form:"category_id"`
	SubCategoryID uint64   `form:"sub_category_id"`
	TagID         uint64   `form:"tag_id"`
	Status        string   `form:"status"`
	Keyword       string   `form:"keyword"`
	Brand         string   `form:"brand"`
	Color         string   `form:"color"`
	Season        string   `form:"season"`
	Occasion      string   `form:"occasion"`
	MinPrice      *float64 `form:"min_price"`
	MaxPrice      *float64 `form:"max_price"`
	Year          *int     `form:"year"`
	Sort          string   `form:"sort"`
	Order         string   `form:"order"`
	Page          int      `form:"page"`
	Size          int      `form:"size"`
}

// Normalize 归一化分页与排序参数。
func (q *ItemQuery) Normalize() {
	if q.Page <= 0 {
		q.Page = 1
	}
	if q.Size <= 0 {
		q.Size = 20
	}
	if q.Size > 100 {
		q.Size = 100
	}
	q.Keyword = strings.TrimSpace(q.Keyword)
	q.Brand = strings.TrimSpace(q.Brand)
	q.Status = strings.TrimSpace(q.Status)
	q.Sort = strings.TrimSpace(q.Sort)
	if strings.ToLower(strings.TrimSpace(q.Order)) == "asc" {
		q.Order = "asc"
	} else {
		q.Order = "desc"
	}
}

// Colors 颜色筛选 code 列表。
func (q ItemQuery) Colors() []string { return SplitCodes(q.Color) }

// Seasons 季节筛选 code 列表。
func (q ItemQuery) Seasons() []string { return SplitCodes(q.Season) }

// Occasions 场合筛选 code 列表。
func (q ItemQuery) Occasions() []string { return SplitCodes(q.Occasion) }

// SplitCodes 把 "white,beige" 拆成 code 数组。
func SplitCodes(raw string) []string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil
	}
	parts := strings.FieldsFunc(raw, func(r rune) bool { return r == ',' })
	return dict.Normalize(parts)
}
