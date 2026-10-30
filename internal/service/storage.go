package service

import (
	"context"
	"errors"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// Storage 图片存储抽象，后续可替换为 OSS / MinIO 实现。
type Storage interface {
	// Save 保存文件并返回可公开访问的 URL 路径。
	Save(ctx context.Context, relPath string, data []byte) (string, error)
	// Delete 按 URL 删除文件。
	Delete(ctx context.Context, url string) error
}

// LocalStorage 本地文件存储，目录结构 uploads/original 与 uploads/thumb。
type LocalStorage struct {
	root string
}

// NewLocalStorage 创建本地存储并确保目录存在。
func NewLocalStorage(root string) (*LocalStorage, error) {
	abs, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}
	for _, sub := range []string{"original", "thumb"} {
		if err := os.MkdirAll(filepath.Join(abs, sub), 0o755); err != nil {
			return nil, err
		}
	}
	return &LocalStorage{root: abs}, nil
}

// Root 返回存储根目录（绝对路径）。
func (s *LocalStorage) Root() string { return s.root }

// Save 实现 Storage。
func (s *LocalStorage) Save(_ context.Context, relPath string, data []byte) (string, error) {
	rel, err := safeRelPath(relPath)
	if err != nil {
		return "", err
	}
	full := filepath.Join(s.root, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		return "", err
	}
	if err := os.WriteFile(full, data, 0o644); err != nil {
		return "", err
	}
	return "/uploads/" + rel, nil
}

// Delete 实现 Storage。
func (s *LocalStorage) Delete(_ context.Context, url string) error {
	rel := strings.TrimPrefix(strings.TrimSpace(url), "/uploads/")
	rel, err := safeRelPath(rel)
	if err != nil {
		return err
	}
	return os.Remove(filepath.Join(s.root, filepath.FromSlash(rel)))
}

// safeRelPath 防止路径穿越。
func safeRelPath(relPath string) (string, error) {
	rel := strings.TrimLeft(filepath.ToSlash(relPath), "/")
	if rel == "" {
		return "", errors.New("文件路径为空")
	}
	clean := path.Clean(rel)
	if clean == "." || clean == ".." || strings.HasPrefix(clean, "../") {
		return "", errors.New("非法文件路径")
	}
	return clean, nil
}
