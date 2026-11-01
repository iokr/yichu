package service

import (
	"context"
	"fmt"
	"time"

	"github.com/iokr/yichu/internal/dto"
	"github.com/iokr/yichu/internal/pkg/dict"
	"github.com/iokr/yichu/internal/pkg/errcode"
	"github.com/iokr/yichu/internal/repository"
)

// StatsService 统计报表。
type StatsService struct {
	items    *repository.ItemRepository
	stats    *repository.StatsRepository
	settings *SettingsService
}

// NewStatsService 构造服务。
func NewStatsService(items *repository.ItemRepository, stats *repository.StatsRepository, settings *SettingsService) *StatsService {
	return &StatsService{items: items, stats: stats, settings: settings}
}

// Overview 统计概览。
func (s *StatsService) Overview(_ context.Context) (*dto.StatsOverview, error) {
	total, err := s.items.Count()
	if err != nil {
		return nil, errcode.Internal(err)
	}
	statusCounts, err := s.items.CountByStatus()
	if err != nil {
		return nil, errcode.Internal(err)
	}
	wearSum, err := s.items.SumWearCount()
	if err != nil {
		return nil, errcode.Internal(err)
	}
	totalAmount, err := s.items.SumPrice(nil)
	if err != nil {
		return nil, errcode.Internal(err)
	}
	pricedCount, err := s.items.PricedCount()
	if err != nil {
		return nil, errcode.Internal(err)
	}

	var avg *float64
	if totalAmount != nil && pricedCount > 0 {
		v := *totalAmount / float64(pricedCount)
		avg = &v
	}

	mask := s.settings.Masker()
	maskedTotal, totalDisplay := mask.Amount(totalAmount)
	maskedAvg, avgDisplay := mask.Amount(avg)

	return &dto.StatsOverview{
		TotalCount:         total,
		WearingCount:       statusCounts[dict.StatusWearing],
		IdleCount:          statusCounts[dict.StatusIdle],
		TotalWearCount:     wearSum,
		TotalAmount:        maskedTotal,
		TotalAmountDisplay: totalDisplay,
		AvgPrice:           maskedAvg,
		AvgPriceDisplay:    avgDisplay,
		PriceMasked:        mask.Hidden(),
	}, nil
}

// ByCategory 按一级分类统计。
func (s *StatsService) ByCategory(_ context.Context) ([]dto.KeyStat, error) {
	total, err := s.items.Count()
	if err != nil {
		return nil, errcode.Internal(err)
	}
	rows, err := s.stats.ByCategory()
	if err != nil {
		return nil, errcode.Internal(err)
	}
	return s.toKeyStats(rows, total), nil
}

// ByBrand 按品牌统计。
func (s *StatsService) ByBrand(_ context.Context) ([]dto.KeyStat, error) {
	total, err := s.items.Count()
	if err != nil {
		return nil, errcode.Internal(err)
	}
	rows, err := s.stats.ByBrand()
	if err != nil {
		return nil, errcode.Internal(err)
	}
	return s.toKeyStats(rows, total), nil
}

// ByYear 按购买年份统计。
func (s *StatsService) ByYear(_ context.Context) ([]dto.KeyStat, error) {
	total, err := s.items.Count()
	if err != nil {
		return nil, errcode.Internal(err)
	}
	rows, err := s.stats.ByYear()
	if err != nil {
		return nil, errcode.Internal(err)
	}
	return s.toKeyStats(rows, total), nil
}

// ByMonth 按购买年月统计，返回该年 12 个月的分布，无购买记录的月份补零。
func (s *StatsService) ByMonth(_ context.Context, year int) ([]dto.KeyStat, error) {
	if year < 1900 || year > 9999 {
		year = time.Now().Year()
	}
	rows, err := s.stats.ByMonth(year)
	if err != nil {
		return nil, errcode.Internal(err)
	}
	found := make(map[string]repository.GroupRow, len(rows))
	var total int64
	for _, row := range rows {
		found[row.Key] = row
		total += row.Count
	}
	filled := make([]repository.GroupRow, 0, 12)
	for month := 1; month <= 12; month++ {
		key := fmt.Sprintf("%04d-%02d", year, month)
		if row, ok := found[key]; ok {
			filled = append(filled, row)
			continue
		}
		filled = append(filled, repository.GroupRow{Key: key, Label: fmt.Sprintf("%d 月", month)})
	}
	return s.toKeyStats(filled, total), nil
}

// BySeason 按季节统计（多选会分别计数）。
func (s *StatsService) BySeason(_ context.Context) ([]dto.CodeStat, error) {
	rows, total, err := s.codeRows()
	if err != nil {
		return nil, err
	}
	counts := make(map[string]int64, len(dict.Seasons))
	for _, row := range rows {
		for _, code := range row.Season {
			counts[code]++
		}
	}
	out := make([]dto.CodeStat, 0, len(dict.Seasons))
	for _, opt := range dict.Seasons {
		out = append(out, dto.CodeStat{
			Code:    opt.Code,
			Label:   opt.Label,
			Count:   counts[opt.Code],
			Percent: percent(counts[opt.Code], total),
		})
	}
	return out, nil
}

// ByColor 按颜色统计。
func (s *StatsService) ByColor(_ context.Context) ([]dto.CodeStat, error) {
	rows, total, err := s.codeRows()
	if err != nil {
		return nil, err
	}
	counts := make(map[string]int64, len(dict.Colors))
	for _, row := range rows {
		for _, code := range row.Color {
			counts[code]++
		}
	}
	out := make([]dto.CodeStat, 0, len(dict.Colors))
	for _, opt := range dict.Colors {
		out = append(out, dto.CodeStat{
			Code:    opt.Code,
			Label:   opt.Label,
			Hex:     opt.Hex,
			Count:   counts[opt.Code],
			Percent: percent(counts[opt.Code], total),
		})
	}
	return out, nil
}

// CostPerWear 单次成本排行。
func (s *StatsService) CostPerWear(_ context.Context, limit int) ([]dto.CostPerWearItem, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	items, err := s.items.CostPerWear(limit)
	if err != nil {
		return nil, errcode.Internal(err)
	}
	mask := s.settings.Masker()
	out := make([]dto.CostPerWearItem, 0, len(items))
	for i := range items {
		it := &items[i]
		price, priceDisplay := mask.Price(it.PurchasePrice)
		var cost *float64
		costDisplay := ""
		if it.PurchasePrice != nil && it.WearCount > 0 {
			v := *it.PurchasePrice / float64(it.WearCount)
			cost, costDisplay = mask.Cost(&v)
		}
		out = append(out, dto.CostPerWearItem{
			ID:            it.ID,
			Name:          it.Name,
			Cover:         it.Cover(),
			WearCount:     it.WearCount,
			PurchasePrice: price,
			PriceDisplay:  priceDisplay,
			CostPerWear:   cost,
			CostDisplay:   costDisplay,
			PriceMasked:   mask.Hidden(),
		})
	}
	return out, nil
}

func (s *StatsService) toKeyStats(rows []repository.GroupRow, total int64) []dto.KeyStat {
	mask := s.settings.Masker()
	out := make([]dto.KeyStat, 0, len(rows))
	for _, row := range rows {
		amount := nullFloat(row.Amount)
		maskedAmount, display := mask.Amount(amount)
		out = append(out, dto.KeyStat{
			Key:           row.Key,
			Label:         row.Label,
			Count:         row.Count,
			Percent:       percent(row.Count, total),
			Amount:        maskedAmount,
			AmountDisplay: display,
		})
	}
	return out
}

func (s *StatsService) codeRows() ([]repository.CodeRow, int64, error) {
	rows, err := s.items.AllCodeRows()
	if err != nil {
		return nil, 0, errcode.Internal(err)
	}
	total, err := s.items.Count()
	if err != nil {
		return nil, 0, errcode.Internal(err)
	}
	return rows, total, nil
}
