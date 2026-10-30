package service

import (
	"context"
	"database/sql"
	"time"

	"github.com/iokr/yichu/internal/dto"
	"github.com/iokr/yichu/internal/pkg/dict"
	"github.com/iokr/yichu/internal/pkg/errcode"
	"github.com/iokr/yichu/internal/repository"
)

// DashboardService 仪表盘聚合。
type DashboardService struct {
	items      *repository.ItemRepository
	stats      *repository.StatsRepository
	categories *repository.CategoryRepository
	settings   *SettingsService
}

// NewDashboardService 构造服务。
func NewDashboardService(items *repository.ItemRepository, stats *repository.StatsRepository, categories *repository.CategoryRepository, settings *SettingsService) *DashboardService {
	return &DashboardService{items: items, stats: stats, categories: categories, settings: settings}
}

// Get 返回仪表盘数据。
func (s *DashboardService) Get(_ context.Context) (*dto.DashboardResponse, error) {
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
	year := time.Now().Year()
	yearAmount, err := s.items.SumPrice(&year)
	if err != nil {
		return nil, errcode.Internal(err)
	}
	catRows, err := s.stats.ByCategory()
	if err != nil {
		return nil, errcode.Internal(err)
	}
	cats, err := s.categories.All()
	if err != nil {
		return nil, errcode.Internal(err)
	}
	recent, err := s.items.Recent(6)
	if err != nil {
		return nil, errcode.Internal(err)
	}

	mask := s.settings.Masker()
	m := newMapper(cats, mask)

	catStats := make([]dto.KeyStat, 0, len(catRows))
	for _, row := range catRows {
		amount := nullFloat(row.Amount)
		maskedAmount, display := mask.Amount(amount)
		catStats = append(catStats, dto.KeyStat{
			Key:           row.Key,
			Label:         row.Label,
			Count:         row.Count,
			Percent:       percent(row.Count, total),
			Amount:        maskedAmount,
			AmountDisplay: display,
		})
	}

	statusStats := make([]dto.CodeStat, 0, len(dict.Statuses))
	for _, opt := range dict.Statuses {
		statusStats = append(statusStats, dto.CodeStat{
			Code:    opt.Code,
			Label:   opt.Label,
			Count:   statusCounts[opt.Code],
			Percent: percent(statusCounts[opt.Code], total),
		})
	}

	recentItems := make([]dto.ItemListItem, 0, len(recent))
	for i := range recent {
		recentItems = append(recentItems, m.item(&recent[i]))
	}

	maskedTotal, totalDisplay := mask.Amount(totalAmount)
	maskedYear, yearDisplay := mask.Amount(yearAmount)

	return &dto.DashboardResponse{
		TotalCount:         total,
		WearingCount:       statusCounts[dict.StatusWearing],
		IdleCount:          statusCounts[dict.StatusIdle],
		TotalWearCount:     wearSum,
		CategoryStats:      catStats,
		StatusStats:        statusStats,
		Year:               year,
		YearAmount:         maskedYear,
		YearAmountDisplay:  yearDisplay,
		TotalAmount:        maskedTotal,
		TotalAmountDisplay: totalDisplay,
		PriceMasked:        mask.Hidden(),
		Recent:             recentItems,
	}, nil
}

func nullFloat(v sql.NullFloat64) *float64 {
	if !v.Valid {
		return nil
	}
	out := v.Float64
	return &out
}
