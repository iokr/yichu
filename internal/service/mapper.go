package service

import (
	"fmt"
	"time"

	"github.com/iokr/yichu/internal/dto"
	"github.com/iokr/yichu/internal/model"
	"github.com/iokr/yichu/internal/pkg/dict"
)

const (
	dateLayout     = "2006-01-02"
	dateTimeLayout = "2006-01-02 15:04:05"
)

// mapper 负责 model → DTO，并统一应用价格脱敏。
type mapper struct {
	mask PriceMasker
	cats map[uint64]model.Category
}

func newMapper(cats []model.Category, mask PriceMasker) *mapper {
	m := &mapper{mask: mask, cats: make(map[uint64]model.Category, len(cats))}
	for _, c := range cats {
		m.cats[c.ID] = c
	}
	return m
}

func (m *mapper) catName(id uint64) string {
	if c, ok := m.cats[id]; ok {
		return c.Name
	}
	return ""
}

func deref(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func datePtr(t *time.Time) *string {
	if t == nil {
		return nil
	}
	v := t.Format(dateLayout)
	return &v
}

func timePtr(t *time.Time) *string {
	if t == nil {
		return nil
	}
	v := t.Format(dateLayout)
	return &v
}

func (m *mapper) item(it *model.Item) dto.ItemListItem {
	price, display := m.mask.Price(it.PurchasePrice)
	row := dto.ItemListItem{
		ID:                   it.ID,
		Name:                 it.Name,
		CategoryID:           it.CategoryID,
		Category:             m.catName(it.CategoryID),
		SubCategoryID:        it.SubCategoryID,
		Brand:                deref(it.Brand),
		Color:                dict.ColorOptions(it.Color),
		Size:                 deref(it.Size),
		Season:               dict.SeasonOptions(it.Season),
		Status:               dict.StatusOption(it.Status),
		PurchaseDate:         datePtr(it.PurchaseDate),
		PurchasePrice:        price,
		PurchasePriceDisplay: display,
		Cover:                it.Cover(),
		WearCount:            it.WearCount,
		LastWornAt:           timePtr(it.LastWornAt),
		PriceMasked:          m.mask.Hidden(),
		CreatedAt:            it.CreatedAt.Format(dateTimeLayout),
	}
	if it.SubCategoryID != nil {
		row.SubCategory = m.catName(*it.SubCategoryID)
	}
	return row
}

func (m *mapper) detail(it *model.Item) dto.ItemDetail {
	detail := dto.ItemDetail{
		ItemListItem:    m.item(it),
		Occasion:        dict.OccasionOptions(it.Occasion),
		PurchaseChannel: deref(it.PurchaseChannel),
		PurchaseURL:     deref(it.PurchaseURL),
		StorageLocation: deref(it.StorageLocation),
		Remark:          deref(it.Remark),
		Images:          make([]dto.ItemImageDTO, 0, len(it.Images)),
		Tags:            make([]dto.TagDTO, 0, len(it.Tags)),
		UpdatedAt:       it.UpdatedAt.Format(dateTimeLayout),
	}
	for _, img := range it.Images {
		detail.Images = append(detail.Images, dto.ItemImageDTO{
			ID:       img.ID,
			URL:      img.URL,
			ThumbURL: img.ThumbURL,
			Sort:     img.Sort,
		})
	}
	for _, tag := range it.Tags {
		detail.Tags = append(detail.Tags, dto.TagDTO{ID: tag.ID, Name: tag.Name})
	}
	if it.PurchasePrice != nil && it.WearCount > 0 {
		cost := *it.PurchasePrice / float64(it.WearCount)
		detail.CostPerWear, detail.CostPerWearDisplay = m.mask.Cost(&cost)
	}
	return detail
}

func joinOptions(options []dict.Option) string {
	if len(options) == 0 {
		return ""
	}
	out := ""
	for i, o := range options {
		if i > 0 {
			out += "、"
		}
		out += o.Label
	}
	return out
}

func percent(count, total int64) float64 {
	if total <= 0 {
		return 0
	}
	v := float64(count) / float64(total) * 100
	return float64(int(v*10+0.5)) / 10
}

func costDisplay(mask PriceMasker, price *float64, wear uint32) (float64, string) {
	if price == nil || wear == 0 {
		return 0, ""
	}
	cost := *price / float64(wear)
	if mask.Hidden() {
		return 0, "¥ ***"
	}
	return cost, fmt.Sprintf("¥%.2f", cost)
}
