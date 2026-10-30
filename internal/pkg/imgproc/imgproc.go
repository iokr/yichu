package imgproc

import (
	"bytes"
	"fmt"
	"image"

	"github.com/disintegration/imaging"
	_ "golang.org/x/image/webp" // 注册 webp 解码器
)

// Detect 返回图片格式（jpeg/png/webp...）与尺寸。
func Detect(data []byte) (string, int, int, error) {
	cfg, format, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		return "", 0, 0, fmt.Errorf("无法识别的图片格式: %w", err)
	}
	return format, cfg.Width, cfg.Height, nil
}

// Decode 解码图片并按 EXIF 自动校正方向。
func Decode(data []byte) (image.Image, error) {
	img, err := imaging.Decode(bytes.NewReader(data), imaging.AutoOrientation(true))
	if err != nil {
		return nil, fmt.Errorf("图片解码失败: %w", err)
	}
	return img, nil
}

// resizeLongEdge 把长边缩放到 edge 以内，不放大。
func resizeLongEdge(img image.Image, edge int) image.Image {
	if edge <= 0 {
		return img
	}
	b := img.Bounds()
	w, h := b.Dx(), b.Dy()
	long := w
	if h > long {
		long = h
	}
	if long <= edge {
		return img
	}
	if w >= h {
		return imaging.Resize(img, edge, 0, imaging.Lanczos)
	}
	return imaging.Resize(img, 0, edge, imaging.Lanczos)
}

// EncodeJPEG 把图片编码为 JPEG。
func EncodeJPEG(img image.Image, quality int) ([]byte, error) {
	if quality <= 0 || quality > 100 {
		quality = 80
	}
	var buf bytes.Buffer
	if err := imaging.Encode(&buf, img, imaging.JPEG, imaging.JPEGQuality(quality)); err != nil {
		return nil, fmt.Errorf("图片编码失败: %w", err)
	}
	return buf.Bytes(), nil
}

// Compress 压缩主图：长边不超过 maxEdge，JPEG 质量 quality。
func Compress(data []byte, maxEdge, quality int) ([]byte, int, int, error) {
	img, err := Decode(data)
	if err != nil {
		return nil, 0, 0, err
	}
	img = resizeLongEdge(img, maxEdge)
	out, err := EncodeJPEG(img, quality)
	if err != nil {
		return nil, 0, 0, err
	}
	b := img.Bounds()
	return out, b.Dx(), b.Dy(), nil
}

// Thumbnail 生成缩略图：长边不超过 edge。
func Thumbnail(data []byte, edge, quality int) ([]byte, error) {
	img, err := Decode(data)
	if err != nil {
		return nil, err
	}
	img = resizeLongEdge(img, edge)
	return EncodeJPEG(img, quality)
}
