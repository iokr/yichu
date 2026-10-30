package service

import (
	"context"

	"github.com/iokr/yichu/internal/dto"
	"github.com/iokr/yichu/internal/pkg/errcode"
	"github.com/iokr/yichu/internal/repository"
)

// TagService 标签查询。
type TagService struct {
	tags *repository.TagRepository
}

// NewTagService 构造服务。
func NewTagService(tags *repository.TagRepository) *TagService {
	return &TagService{tags: tags}
}

// List 返回全部标签及其关联衣物数。
func (s *TagService) List(_ context.Context) ([]dto.TagDTO, error) {
	list, err := s.tags.All()
	if err != nil {
		return nil, errcode.Internal(err)
	}
	counts, cerr := s.tags.Counts()
	if cerr != nil {
		return nil, errcode.Internal(cerr)
	}
	out := make([]dto.TagDTO, 0, len(list))
	for _, tag := range list {
		out = append(out, dto.TagDTO{ID: tag.ID, Name: tag.Name, Count: counts[tag.ID]})
	}
	return out, nil
}
