package service

import (
	"context"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/iokr/yichu/internal/dto"
	"github.com/iokr/yichu/internal/model"
	"github.com/iokr/yichu/internal/pkg/dict"
	"github.com/iokr/yichu/internal/pkg/errcode"
	"github.com/iokr/yichu/internal/repository"
)

// WearService 穿着打卡。
type WearService struct {
	items *repository.ItemRepository
	logs  *repository.WearLogRepository
}

// NewWearService 构造服务。
func NewWearService(items *repository.ItemRepository, logs *repository.WearLogRepository) *WearService {
	return &WearService{items: items, logs: logs}
}

// Create 新增穿着记录，同时累加穿着次数。
func (s *WearService) Create(_ context.Context, itemID uint64, req dto.WearLogRequest) (*dto.WearLogItem, error) {
	item, err := s.items.Get(itemID)
	if err != nil {
		return nil, errcode.Internal(err)
	}
	if item == nil {
		return nil, errcode.NotFound("衣物不存在")
	}

	wornDate := time.Now()
	if raw := strings.TrimSpace(req.WornDate); raw != "" {
		t, perr := time.ParseInLocation(dateLayout, raw, time.Local)
		if perr != nil {
			return nil, errcode.Param("穿着日期格式应为 YYYY-MM-DD")
		}
		wornDate = t
	}
	today := time.Now()
	if wornDate.After(time.Date(today.Year(), today.Month(), today.Day(), 23, 59, 59, 0, time.Local)) {
		return nil, errcode.Param("穿着日期不能晚于今天")
	}

	var occasion *string
	if raw := strings.TrimSpace(req.Occasion); raw != "" {
		if !dict.IsOccasion(raw) {
			return nil, errcode.Param("无效的场合：%s", raw)
		}
		occasion = &raw
	}
	remark := strings.TrimSpace(req.Remark)
	if utf8.RuneCountInString(remark) > 255 {
		return nil, errcode.Param("备注不能超过 255 个字符")
	}
	var remarkPtr *string
	if remark != "" {
		remarkPtr = &remark
	}

	log := &model.WearLog{
		ItemID:   itemID,
		WornDate: wornDate,
		Occasion: occasion,
		Remark:   remarkPtr,
	}
	if err := s.logs.Add(log); err != nil {
		return nil, errcode.Internal(err)
	}
	return toWearLogItem(log), nil
}

// List 查询某件衣物的穿着记录。
func (s *WearService) List(_ context.Context, itemID uint64) ([]dto.WearLogItem, error) {
	item, err := s.items.Get(itemID)
	if err != nil {
		return nil, errcode.Internal(err)
	}
	if item == nil {
		return nil, errcode.NotFound("衣物不存在")
	}
	logs, err := s.logs.ListByItem(itemID, 0)
	if err != nil {
		return nil, errcode.Internal(err)
	}
	return s.MapLogs(logs), nil
}

// MapLogs 批量转换（备份复用）。
func (s *WearService) MapLogs(logs []model.WearLog) []dto.WearLogItem {
	out := make([]dto.WearLogItem, 0, len(logs))
	for i := range logs {
		out = append(out, *toWearLogItem(&logs[i]))
	}
	return out
}

// All 返回全部穿着记录。
func (s *WearService) All(_ context.Context) ([]model.WearLog, error) {
	logs, err := s.logs.All()
	if err != nil {
		return nil, errcode.Internal(err)
	}
	return logs, nil
}

func toWearLogItem(log *model.WearLog) *dto.WearLogItem {
	item := &dto.WearLogItem{
		ID:        log.ID,
		ItemID:    log.ItemID,
		WornDate:  log.WornDate.Format(dateLayout),
		Remark:    deref(log.Remark),
		CreatedAt: log.CreatedAt.Format(dateTimeLayout),
	}
	if log.Occasion != nil && *log.Occasion != "" {
		var opt dict.Option
		if found, ok := dict.FindOccasion(*log.Occasion); ok {
			opt = found
		} else {
			opt = dict.Option{Code: *log.Occasion, Label: *log.Occasion}
		}
		item.Occasion = &opt
	}
	return item
}
