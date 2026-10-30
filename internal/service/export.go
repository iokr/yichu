package service

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"

	"github.com/iokr/yichu/internal/dto"
	"github.com/iokr/yichu/internal/model"
	"github.com/iokr/yichu/internal/pkg/dict"
	"github.com/iokr/yichu/internal/pkg/errcode"
)

// ExportResult 导出结果。
type ExportResult struct {
	Data        []byte
	Filename    string
	ContentType string
	Masked      bool
}

var exportHeader = []string{
	"ID", "名称", "一级分类", "二级分类", "品牌", "颜色", "尺码", "季节", "场合",
	"购买日期", "价格", "购买渠道", "订单链接", "状态", "存放位置", "穿着次数", "最近穿着", "标签", "备注", "添加时间",
}

// Export 导出 CSV / Excel，价格列受脱敏约束。
func (s *ItemService) Export(_ context.Context, q dto.ItemQuery, format string) (*ExportResult, error) {
	q.Normalize()
	query := s.repoQuery(q)
	query.Page = 1
	query.Size = 0

	items, err := s.items.All(query)
	if err != nil {
		return nil, errcode.Internal(err)
	}
	cats, err := s.categories.All()
	if err != nil {
		return nil, errcode.Internal(err)
	}
	mask := s.masker()
	m := newMapper(cats, mask)

	rows := make([][]string, 0, len(items))
	for i := range items {
		it := &items[i]
		subName := ""
		if it.SubCategoryID != nil {
			subName = m.catName(*it.SubCategoryID)
		}
		rows = append(rows, []string{
			strconv.FormatUint(it.ID, 10),
			it.Name,
			m.catName(it.CategoryID),
			subName,
			deref(it.Brand),
			joinOptions(dict.ColorOptions(it.Color)),
			deref(it.Size),
			joinOptions(dict.SeasonOptions(it.Season)),
			joinOptions(dict.OccasionOptions(it.Occasion)),
			formatDate(it.PurchaseDate),
			mask.CSV(it.PurchasePrice),
			deref(it.PurchaseChannel),
			deref(it.PurchaseURL),
			dict.StatusLabel(it.Status),
			deref(it.StorageLocation),
			strconv.FormatUint(uint64(it.WearCount), 10),
			formatDate(it.LastWornAt),
			joinTagNames(it.Tags),
			deref(it.Remark),
			it.CreatedAt.Format(dateTimeLayout),
		})
	}

	stamp := time.Now().Format("20060102-150405")
	switch strings.ToLower(strings.TrimSpace(format)) {
	case "excel", "xlsx":
		data, xerr := writeXLSX(exportHeader, rows)
		if xerr != nil {
			return nil, errcode.Internal(xerr)
		}
		return &ExportResult{
			Data:        data,
			Filename:    fmt.Sprintf("yichu-items-%s.xlsx", stamp),
			ContentType: "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
			Masked:      mask.Hidden(),
		}, nil
	case "", "csv":
		data, cerr := writeCSV(exportHeader, rows)
		if cerr != nil {
			return nil, errcode.Internal(cerr)
		}
		return &ExportResult{
			Data:        data,
			Filename:    fmt.Sprintf("yichu-items-%s.csv", stamp),
			ContentType: "text/csv; charset=utf-8",
			Masked:      mask.Hidden(),
		}, nil
	default:
		return nil, errcode.Param("不支持的导出格式：%s", format)
	}
}

func writeCSV(header []string, rows [][]string) ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteString("\xEF\xBB\xBF") // BOM，保证 Excel 打开中文不乱码
	w := csv.NewWriter(&buf)
	if err := w.Write(header); err != nil {
		return nil, err
	}
	for _, row := range rows {
		if err := w.Write(row); err != nil {
			return nil, err
		}
	}
	w.Flush()
	if err := w.Error(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func writeXLSX(header []string, rows [][]string) ([]byte, error) {
	f := excelize.NewFile()
	defer func() { _ = f.Close() }()

	const sheet = "衣物"
	index, err := f.NewSheet(sheet)
	if err != nil {
		return nil, err
	}
	_ = f.DeleteSheet("Sheet1")
	f.SetActiveSheet(index)

	all := make([][]string, 0, len(rows)+1)
	all = append(all, header)
	all = append(all, rows...)
	for i := range all {
		cell, cerr := excelize.CoordinatesToCellName(1, i+1)
		if cerr != nil {
			return nil, cerr
		}
		row := all[i]
		if serr := f.SetSheetRow(sheet, cell, &row); serr != nil {
			return nil, serr
		}
	}
	style, _ := f.NewStyle(&excelize.Style{Font: &excelize.Font{Bold: true}})
	lastCell, _ := excelize.CoordinatesToCellName(len(header), 1)
	_ = f.SetCellStyle(sheet, "A1", lastCell, style)
	_ = f.SetColWidth(sheet, "A", "A", 6)
	_ = f.SetColWidth(sheet, "B", "B", 28)
	_ = f.SetColWidth(sheet, "J", "J", 12)
	_ = f.SetColWidth(sheet, "K", "K", 10)
	_ = f.SetColWidth(sheet, "M", "M", 32)
	_ = f.SetColWidth(sheet, "R", "R", 20)

	buf, werr := f.WriteToBuffer()
	if werr != nil {
		return nil, werr
	}
	return buf.Bytes(), nil
}

func formatDate(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format(dateLayout)
}

func joinTagNames(tags []model.Tag) string {
	if len(tags) == 0 {
		return ""
	}
	names := make([]string, 0, len(tags))
	for _, tag := range tags {
		names = append(names, tag.Name)
	}
	return strings.Join(names, "、")
}
