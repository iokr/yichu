package service

import (
	"context"
	"strings"
	"unicode/utf8"

	"github.com/iokr/yichu/internal/dto"
	"github.com/iokr/yichu/internal/model"
	"github.com/iokr/yichu/internal/pkg/errcode"
	"github.com/iokr/yichu/internal/repository"
)

// CategoryService 分类维护。
type CategoryService struct {
	cats *repository.CategoryRepository
}

// NewCategoryService 构造服务。
func NewCategoryService(cats *repository.CategoryRepository) *CategoryService {
	return &CategoryService{cats: cats}
}

// Tree 返回一级分类 + 二级分类树，并带上关联衣物数。
func (s *CategoryService) Tree(_ context.Context) ([]dto.CategoryNode, error) {
	list, err := s.cats.All()
	if err != nil {
		return nil, errcode.Internal(err)
	}
	counts, err := s.cats.ItemCounts()
	if err != nil {
		return nil, errcode.Internal(err)
	}
	subCounts, err := s.cats.SubItemCounts()
	if err != nil {
		return nil, errcode.Internal(err)
	}

	tops := make([]dto.CategoryNode, 0, 4)
	index := make(map[uint64]int, 4)
	for _, c := range list {
		if !c.IsTop() {
			continue
		}
		tops = append(tops, dto.CategoryNode{
			ID:        c.ID,
			ParentID:  c.ParentID,
			Name:      c.Name,
			Sort:      c.Sort,
			IsSystem:  c.IsSystem,
			ItemCount: counts[c.ID],
			Children:  make([]dto.CategoryNode, 0, 8),
		})
		index[c.ID] = len(tops) - 1
	}
	for _, c := range list {
		if c.IsTop() {
			continue
		}
		i, ok := index[c.ParentID]
		if !ok {
			continue
		}
		tops[i].Children = append(tops[i].Children, dto.CategoryNode{
			ID:        c.ID,
			ParentID:  c.ParentID,
			Name:      c.Name,
			Sort:      c.Sort,
			IsSystem:  c.IsSystem,
			ItemCount: subCounts[c.ID],
		})
	}
	return tops, nil
}

// Create 新增分类：ParentID=0 为一级分类，否则挂在对应一级分类下。
func (s *CategoryService) Create(_ context.Context, req dto.CategoryRequest) (*dto.CategoryNode, error) {
	name, err := normalizeCategoryName(req.Name)
	if err != nil {
		return nil, err
	}
	if req.ParentID != 0 {
		parent, gerr := s.cats.Get(req.ParentID)
		if gerr != nil {
			return nil, errcode.Internal(gerr)
		}
		if parent == nil || !parent.IsTop() {
			return nil, errcode.Param("二级分类只能挂在一级分类下")
		}
	}
	dup, derr := s.cats.NameExists(req.ParentID, name, 0)
	if derr != nil {
		return nil, errcode.Internal(derr)
	}
	if dup {
		return nil, errcode.Param("同层已存在名为「%s」的分类", name)
	}

	list, lerr := s.cats.All()
	if lerr != nil {
		return nil, errcode.Internal(lerr)
	}
	sort := req.Sort
	if sort <= 0 {
		maxSort := 0
		for _, c := range list {
			if c.ParentID == req.ParentID && c.Sort > maxSort {
				maxSort = c.Sort
			}
		}
		sort = maxSort + 1
	}

	cat := &model.Category{Name: name, ParentID: req.ParentID, Sort: sort}
	if err := s.cats.Create(cat); err != nil {
		return nil, errcode.Internal(err)
	}
	return &dto.CategoryNode{
		ID:       cat.ID,
		ParentID: cat.ParentID,
		Name:     cat.Name,
		Sort:     cat.Sort,
	}, nil
}

// Update 修改分类名称/排序；一级分类不可移动到其它分类下。
func (s *CategoryService) Update(_ context.Context, id uint64, req dto.CategoryRequest) (*dto.CategoryNode, error) {
	cat, err := s.cats.Get(id)
	if err != nil {
		return nil, errcode.Internal(err)
	}
	if cat == nil {
		return nil, errcode.NotFound("分类不存在")
	}
	name, err := normalizeCategoryName(req.Name)
	if err != nil {
		return nil, err
	}

	parentID := cat.ParentID
	if cat.IsTop() {
		if req.ParentID != 0 {
			return nil, errcode.Forbidden("一级分类不能移动到其它分类下")
		}
	} else if req.ParentID != cat.ParentID {
		parent, gerr := s.cats.Get(req.ParentID)
		if gerr != nil {
			return nil, errcode.Internal(gerr)
		}
		if parent == nil || !parent.IsTop() {
			return nil, errcode.Param("二级分类只能挂在一级分类下")
		}
		parentID = req.ParentID
	}

	dup, derr := s.cats.NameExists(parentID, name, cat.ID)
	if derr != nil {
		return nil, errcode.Internal(derr)
	}
	if dup {
		return nil, errcode.Param("同层已存在名为「%s」的分类", name)
	}

	cat.Name = name
	cat.ParentID = parentID
	if req.Sort > 0 {
		cat.Sort = req.Sort
	}
	if err := s.cats.Update(cat); err != nil {
		return nil, errcode.Internal(err)
	}
	return &dto.CategoryNode{
		ID:       cat.ID,
		ParentID: cat.ParentID,
		Name:     cat.Name,
		Sort:     cat.Sort,
		IsSystem: cat.IsSystem,
	}, nil
}

// Delete 删除分类：一级分类不可删；有衣物关联时禁止删除。
func (s *CategoryService) Delete(_ context.Context, id uint64) error {
	cat, err := s.cats.Get(id)
	if err != nil {
		return errcode.Internal(err)
	}
	if cat == nil {
		return errcode.NotFound("分类不存在")
	}
	if cat.IsTop() {
		return errcode.Forbidden("一级分类不可删除")
	}
	count, cerr := s.cats.ItemCount(id)
	if cerr != nil {
		return errcode.Internal(cerr)
	}
	if count > 0 {
		return errcode.Forbidden("该分类下还有 %d 件衣物，请先调整后再删除", count)
	}
	if err := s.cats.Delete(id); err != nil {
		return errcode.Internal(err)
	}
	return nil
}

// Sort 批量调整排序。
func (s *CategoryService) Sort(_ context.Context, items []dto.CategorySortItem) error {
	if len(items) == 0 {
		return errcode.Param("排序列表不能为空")
	}
	if len(items) > 200 {
		return errcode.Param("单次最多调整 200 个分类")
	}
	all, err := s.cats.All()
	if err != nil {
		return errcode.Internal(err)
	}
	exists := make(map[uint64]struct{}, len(all))
	for _, c := range all {
		exists[c.ID] = struct{}{}
	}
	sorts := make([]repository.SortItem, 0, len(items))
	for _, item := range items {
		if _, ok := exists[item.ID]; !ok {
			return errcode.Param("分类 %d 不存在", item.ID)
		}
		if item.Sort < 0 {
			return errcode.Param("排序号不能为负数")
		}
		sorts = append(sorts, repository.SortItem{ID: item.ID, Sort: item.Sort})
	}
	if err := s.cats.UpdateSorts(sorts); err != nil {
		return errcode.Internal(err)
	}
	return nil
}

func normalizeCategoryName(raw string) (string, error) {
	name := strings.TrimSpace(raw)
	if name == "" {
		return "", errcode.Param("分类名称不能为空")
	}
	if utf8.RuneCountInString(name) > 50 {
		return "", errcode.Param("分类名称不能超过 50 个字符")
	}
	return name, nil
}
