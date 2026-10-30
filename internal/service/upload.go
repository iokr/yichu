package service

import (
	"context"
	"io"
	"mime/multipart"
	"net/http"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/iokr/yichu/internal/dto"
	"github.com/iokr/yichu/internal/pkg/errcode"
	"github.com/iokr/yichu/internal/pkg/imgproc"
	"github.com/iokr/yichu/internal/pkg/logx"
)

const (
	originalMaxEdge = 1600
	thumbMaxEdge    = 300
	jpegQuality     = 80
)

// UploadService 图片上传、压缩与缩略图。
type UploadService struct {
	storage  Storage
	maxBytes int64
	log      *logx.Logger
}

// NewUploadService 构造服务。
func NewUploadService(storage Storage, maxBytes int64, log *logx.Logger) *UploadService {
	return &UploadService{storage: storage, maxBytes: maxBytes, log: log}
}

// Save 校验并保存图片，返回主图与缩略图地址。
func (s *UploadService) Save(ctx context.Context, header *multipart.FileHeader) (*dto.UploadResult, error) {
	if header == nil {
		return nil, errcode.Param("请选择要上传的图片")
	}
	if s.maxBytes > 0 && header.Size > s.maxBytes {
		return nil, errcode.Param("单张图片不能超过 %d MB", int(s.maxBytes/(1024*1024)))
	}

	ext := strings.ToLower(filepath.Ext(header.Filename))
	switch ext {
	case ".jpg", ".jpeg", ".png", ".webp":
	default:
		return nil, errcode.Param("仅支持 jpg / png / webp 图片")
	}

	file, err := header.Open()
	if err != nil {
		return nil, errcode.Internal(err)
	}
	defer func() { _ = file.Close() }()

	limit := s.maxBytes + 1
	if limit <= 1 {
		limit = 10*1024*1024 + 1
	}
	data, err := io.ReadAll(io.LimitReader(file, limit))
	if err != nil {
		return nil, errcode.Internal(err)
	}
	if s.maxBytes > 0 && int64(len(data)) > s.maxBytes {
		return nil, errcode.Param("单张图片不能超过 %d MB", int(s.maxBytes/(1024*1024)))
	}
	if len(data) == 0 {
		return nil, errcode.Param("图片内容为空")
	}
	switch http.DetectContentType(data) {
	case "image/jpeg", "image/png", "image/webp":
	default:
		return nil, errcode.Param("仅支持 jpg / png / webp 图片")
	}
	if _, _, _, derr := imgproc.Detect(data); derr != nil {
		return nil, errcode.Param("图片已损坏或格式不受支持")
	}

	compressed, width, height, err := imgproc.Compress(data, originalMaxEdge, jpegQuality)
	if err != nil {
		return nil, errcode.Param("图片处理失败，请更换图片重试")
	}
	thumb, err := imgproc.Thumbnail(data, thumbMaxEdge, jpegQuality)
	if err != nil {
		return nil, errcode.Param("缩略图生成失败，请更换图片重试")
	}

	now := time.Now()
	relDir := path.Join(now.Format("2006"), now.Format("01"))
	name := uuid.NewString() + ".jpg"

	mainURL, err := s.storage.Save(ctx, path.Join("original", relDir, name), compressed)
	if err != nil {
		return nil, errcode.Internal(err)
	}
	thumbURL, err := s.storage.Save(ctx, path.Join("thumb", relDir, name), thumb)
	if err != nil {
		_ = s.storage.Delete(ctx, mainURL)
		return nil, errcode.Internal(err)
	}

	s.log.Debugf("图片已保存: %s (%dx%d, %d bytes)", mainURL, width, height, len(compressed))
	return &dto.UploadResult{
		URL:      mainURL,
		ThumbURL: thumbURL,
		Width:    width,
		Height:   height,
		Size:     int64(len(compressed)),
	}, nil
}
