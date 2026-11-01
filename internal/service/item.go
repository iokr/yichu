package service

import (
	"context"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/iokr/yichu/internal/config"
	"github.com/iokr/yichu/internal/dto"
	"github.com/iokr/yichu/internal/model"
	"github.com/iokr/yichu/internal/pkg/dict"
	"github.com/iokr/yichu/internal/pkg/errcode"
	"github.com/iokr/yichu/internal/repository"
)

// ItemService 衣物业务逻辑。
type ItemService struct {
	items      *repository.ItemRepository
	categories *repository.CategoryRepository
	settings   *SettingsService
	cfg        *config.Config
}

// NewItemService 构造服务。
func NewItemService(items *repository.ItemRepository, categories *repository.CategoryRepository, settings *SettingsService, cfg *config.Config) *ItemService {
	return &ItemService{items: items, categories: categories, settings: settings, cfg: cfg}
}

func (s *ItemService) masker() PriceMasker { return s.settings.Masker() }

// List 分页查询衣物。
func (s *ItemService) List(_ context.Context, q dto.ItemQuery) (*dto.PageResult[dto.ItemListItem], error) {
	q.Normalize()
	query, err := s.repoQuery(q)
	if err != nil {
		return nil, err
	}
	items, total, err := s.items.List(query)
	if err != nil {
		return nil, errcode.Internal(err)
	}
	cats, err := s.categories.All()
	if err != nil {
		return nil, errcode.Internal(err)
	}
	m := newMapper(cats, s.masker())
	list := make([]dto.ItemListItem, 0, len(items))
	for i := range items {
		list = append(list, m.item(&items[i]))
	}
	return &dto.PageResult[dto.ItemListItem]{List: list, Total: total, Page: q.Page, Size: q.Size}, nil
}

// AllDetails 返回全部衣物详情（备份 / 统计导出用）。
func (s *ItemService) AllDetails(_ context.Context) ([]dto.ItemDetail, bool, error) {
	items, err := s.items.All(repository.ItemQuery{Sort: "created_at", Order: "asc"})
	if err != nil {
		return nil, false, errcode.Internal(err)
	}
	cats, err := s.categories.All()
	if err != nil {
		return nil, false, errcode.Internal(err)
	}
	mask := s.masker()
	m := newMapper(cats, mask)
	out := make([]dto.ItemDetail, 0, len(items))
	for i := range items {
		out = append(out, m.detail(&items[i]))
	}
	return out, mask.Hidden(), nil
}

// Get 衣物详情。
func (s *ItemService) Get(_ context.Context, id uint64) (*dto.ItemDetail, error) {
	item, err := s.items.Get(id)
	if err != nil {
		return nil, errcode.Internal(err)
	}
	if item == nil {
		return nil, errcode.NotFound("衣物不存在")
	}
	cats, err := s.categories.All()
	if err != nil {
		return nil, errcode.Internal(err)
	}
	m := newMapper(cats, s.masker())
	detail := m.detail(item)
	return &detail, nil
}

// Create 新增衣物。
func (s *ItemService) Create(ctx context.Context, req dto.ItemRequest) (*dto.ItemDetail, error) {
	item, images, tags, err := s.prepare(req)
	if err != nil {
		return nil, err
	}
	item.Images = images
	if err := s.items.Create(item, tags); err != nil {
		return nil, errcode.Internal(err)
	}
	return s.Get(ctx, item.ID)
}

// Update 全量更新衣物。
func (s *ItemService) Update(ctx context.Context, id uint64, req dto.ItemRequest) (*dto.ItemDetail, error) {
	existing, err := s.items.Get(id)
	if err != nil {
		return nil, errcode.Internal(err)
	}
	if existing == nil {
		return nil, errcode.NotFound("衣物不存在")
	}
	item, images, tags, err := s.prepare(req)
	if err != nil {
		return nil, err
	}
	fields := map[string]any{
		"name":             item.Name,
		"category_id":      item.CategoryID,
		"sub_category_id":  item.SubCategoryID,
		"brand":            item.Brand,
		"color":            item.Color,
		"size":             item.Size,
		"season":           item.Season,
		"occasion":         item.Occasion,
		"purchase_date":    item.PurchaseDate,
		"purchase_price":   item.PurchasePrice,
		"purchase_channel": item.PurchaseChannel,
		"purchase_url":     item.PurchaseURL,
		"status":           item.Status,
		"storage_location": item.StorageLocation,
		"remark":           item.Remark,
	}
	if err := s.items.Update(id, fields, images, tags); err != nil {
		return nil, errcode.Internal(err)
	}
	return s.Get(ctx, id)
}

// Delete 删除衣物（软删除）。
func (s *ItemService) Delete(_ context.Context, id uint64) error {
	item, err := s.items.Get(id)
	if err != nil {
		return errcode.Internal(err)
	}
	if item == nil {
		return errcode.NotFound("衣物不存在")
	}
	if err := s.items.Delete(id); err != nil {
		return errcode.Internal(err)
	}
	return nil
}

// Batch 批量删除或修改状态。
func (s *ItemService) Batch(_ context.Context, req dto.BatchRequest) error {
	if len(req.IDs) == 0 {
		return errcode.Param("请选择要操作的衣物")
	}
	if len(req.IDs) > 200 {
		return errcode.Param("单次最多操作 200 件衣物")
	}
	switch req.Action {
	case dto.BatchActionDelete:
		if err := s.items.BatchDelete(req.IDs); err != nil {
			return errcode.Internal(err)
		}
	case dto.BatchActionStatus:
		status := strings.TrimSpace(req.Status)
		if !dict.IsStatus(status) {
			return errcode.Param("无效的状态：%s", status)
		}
		if err := s.items.BatchStatus(req.IDs, status); err != nil {
			return errcode.Internal(err)
		}
	default:
		return errcode.Param("不支持的批量操作：%s", req.Action)
	}
	return nil
}

func (s *ItemService) repoQuery(q dto.ItemQuery) (repository.ItemQuery, error) {
	purchaseFrom, err := parseDateParam("购买开始日期", q.PurchaseFrom)
	if err != nil {
		return repository.ItemQuery{}, err
	}
	purchaseTo, err := parseDateParam("购买结束日期", q.PurchaseTo)
	if err != nil {
		return repository.ItemQuery{}, err
	}
	if purchaseFrom != nil && purchaseTo != nil && purchaseTo.Before(*purchaseFrom) {
		return repository.ItemQuery{}, errcode.Param("购买时间范围有误：结束日期不能早于开始日期")
	}
	return repository.ItemQuery{
		CategoryID:    q.CategoryID,
		SubCategoryID: q.SubCategoryID,
		TagID:         q.TagID,
		Status:        q.Status,
		Keyword:       q.Keyword,
		Brand:         q.Brand,
		Colors:        q.Colors(),
		Seasons:       q.Seasons(),
		Occasions:     q.Occasions(),
		MinPrice:      q.MinPrice,
		MaxPrice:      q.MaxPrice,
		Year:          q.Year,
		PurchaseFrom:  purchaseFrom,
		PurchaseTo:    purchaseTo,
		Sort:          q.Sort,
		Order:         q.Order,
		Page:          q.Page,
		Size:          q.Size,
	}, nil
}

// prepare 校验请求并转换为模型，返回图片与标签。
func (s *ItemService) prepare(req dto.ItemRequest) (*model.Item, []model.ItemImage, []string, error) {
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, nil, nil, errcode.Param("衣物名称不能为空")
	}
	if utf8.RuneCountInString(name) > 100 {
		return nil, nil, nil, errcode.Param("衣物名称不能超过 100 个字符")
	}
	if req.CategoryID == 0 {
		return nil, nil, nil, errcode.Param("请选择一级分类")
	}
	cat, err := s.categories.Get(req.CategoryID)
	if err != nil {
		return nil, nil, nil, errcode.Internal(err)
	}
	if cat == nil || !cat.IsTop() {
		return nil, nil, nil, errcode.Param("一级分类不存在")
	}

	var subID *uint64
	if req.SubCategoryID != nil && *req.SubCategoryID > 0 {
		sub, serr := s.categories.Get(*req.SubCategoryID)
		if serr != nil {
			return nil, nil, nil, errcode.Internal(serr)
		}
		if sub == nil || sub.ParentID != cat.ID {
			return nil, nil, nil, errcode.Param("二级分类与一级分类不匹配")
		}
		subID = req.SubCategoryID
	}

	colors := dict.Normalize(req.Color)
	if bad := dict.Validate(colors, dict.IsColor); bad != "" {
		return nil, nil, nil, errcode.Param("无效的颜色：%s", bad)
	}
	seasons := dict.Normalize(req.Season)
	if bad := dict.Validate(seasons, dict.IsSeason); bad != "" {
		return nil, nil, nil, errcode.Param("无效的季节：%s", bad)
	}
	occasions := dict.Normalize(req.Occasion)
	if bad := dict.Validate(occasions, dict.IsOccasion); bad != "" {
		return nil, nil, nil, errcode.Param("无效的场合：%s", bad)
	}

	status := strings.TrimSpace(req.Status)
	if status == "" {
		status = dict.StatusWearing
	}
	if !dict.IsStatus(status) {
		return nil, nil, nil, errcode.Param("无效的状态：%s", status)
	}

	var purchaseDate *time.Time
	if raw := strings.TrimSpace(req.PurchaseDate); raw != "" {
		t, perr := time.ParseInLocation(dateLayout, raw, time.Local)
		if perr != nil {
			return nil, nil, nil, errcode.Param("购买日期格式应为 YYYY-MM-DD")
		}
		purchaseDate = &t
	}

	if req.PurchasePrice != nil {
		if *req.PurchasePrice < 0 {
			return nil, nil, nil, errcode.Param("价格不能为负数")
		}
		if *req.PurchasePrice > 99999999.99 {
			return nil, nil, nil, errcode.Param("价格超出可记录范围")
		}
	}

	brand, err := limitedPtr("品牌", req.Brand, 50)
	if err != nil {
		return nil, nil, nil, err
	}
	size, err := limitedPtr("尺码", req.Size, 20)
	if err != nil {
		return nil, nil, nil, err
	}
	channel, err := limitedPtr("购买渠道", req.PurchaseChannel, 50)
	if err != nil {
		return nil, nil, nil, err
	}
	purchaseURL, err := limitedPtr("订单链接", req.PurchaseURL, 500)
	if err != nil {
		return nil, nil, nil, err
	}
	if purchaseURL != nil && !strings.HasPrefix(*purchaseURL, "http://") && !strings.HasPrefix(*purchaseURL, "https://") {
		return nil, nil, nil, errcode.Param("订单链接需以 http:// 或 https:// 开头")
	}
	location, err := limitedPtr("存放位置", req.StorageLocation, 50)
	if err != nil {
		return nil, nil, nil, err
	}
	remark, err := limitedPtr("备注", req.Remark, 500)
	if err != nil {
		return nil, nil, nil, err
	}

	maxCount := s.cfg.Storage.MaxCount
	if len(req.Images) > maxCount {
		return nil, nil, nil, errcode.Param("最多上传 %d 张图片", maxCount)
	}
	images := make([]model.ItemImage, 0, len(req.Images))
	for i, img := range req.Images {
		url := strings.TrimSpace(img.URL)
		thumb := strings.TrimSpace(img.ThumbURL)
		if url == "" {
			return nil, nil, nil, errcode.Param("第 %d 张图片地址为空", i+1)
		}
		if !strings.HasPrefix(url, "/uploads/") {
			return nil, nil, nil, errcode.Param("第 %d 张图片地址不合法", i+1)
		}
		if thumb == "" {
			thumb = url
		}
		if !strings.HasPrefix(thumb, "/uploads/") {
			return nil, nil, nil, errcode.Param("第 %d 张缩略图地址不合法", i+1)
		}
		images = append(images, model.ItemImage{URL: url, ThumbURL: thumb, Sort: i})
	}

	tagNames := make([]string, 0, len(req.Tags))
	for _, raw := range req.Tags {
		tag := strings.TrimSpace(raw)
		if tag == "" {
			continue
		}
		if utf8.RuneCountInString(tag) > 20 {
			return nil, nil, nil, errcode.Param("单个标签不能超过 20 个字符")
		}
		tagNames = append(tagNames, tag)
	}
	if len(tagNames) > 20 {
		return nil, nil, nil, errcode.Param("标签最多 20 个")
	}

	item := &model.Item{
		Name:            name,
		CategoryID:      cat.ID,
		SubCategoryID:   subID,
		Brand:           brand,
		Color:           codesOrNil(colors),
		Size:            size,
		Season:          codesOrNil(seasons),
		Occasion:        codesOrNil(occasions),
		PurchaseDate:    purchaseDate,
		PurchasePrice:   req.PurchasePrice,
		PurchaseChannel: channel,
		PurchaseURL:     purchaseURL,
		Status:          status,
		StorageLocation: location,
		Remark:          remark,
	}
	return item, images, tagNames, nil
}

// parseDateParam 解析可选的日期查询参数，空字符串表示不限制。
func parseDateParam(field, raw string) (*time.Time, error) {
	value := strings.TrimSpace(raw)
	if value == "" {
		return nil, nil
	}
	t, err := time.ParseInLocation(dateLayout, value, time.Local)
	if err != nil {
		return nil, errcode.Param("%s格式应为 YYYY-MM-DD", field)
	}
	return &t, nil
}

func limitedPtr(field, raw string, max int) (*string, error) {
	value := strings.TrimSpace(raw)
	if value == "" {
		return nil, nil
	}
	if utf8.RuneCountInString(value) > max {
		return nil, errcode.Param("%s不能超过 %d 个字符", field, max)
	}
	return &value, nil
}

func codesOrNil(codes []string) model.Codes {
	if len(codes) == 0 {
		return nil
	}
	return model.Codes(codes)
}
