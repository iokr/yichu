package logx

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// Level 日志级别。
type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
)

// ParseLevel 解析级别字符串，未知值回退 info。
func ParseLevel(s string) Level {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "debug":
		return LevelDebug
	case "warn", "warning":
		return LevelWarn
	case "error":
		return LevelError
	default:
		return LevelInfo
	}
}

func levelTag(l Level) string {
	switch l {
	case LevelDebug:
		return "[debug] "
	case LevelWarn:
		return "[warn] "
	case LevelError:
		return "[error] "
	default:
		return "[info] "
	}
}

// dailyWriter 按天切割的日志写入器：当前写入 yichu.log，
// 跨天时把旧文件重命名为 yichu-YYYY-MM-DD.log 后继续写。
type dailyWriter struct {
	mu   sync.Mutex
	dir  string
	name string
	path string
	day  string
	f    *os.File
}

func newDailyWriter(file string) (*dailyWriter, error) {
	if strings.TrimSpace(file) == "" {
		file = "./logs/yichu.log"
	}
	dir := filepath.Dir(file)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("创建日志目录失败: %w", err)
	}
	w := &dailyWriter{dir: dir, name: filepath.Base(file), path: file}
	if err := w.rotateLocked(); err != nil {
		return nil, err
	}
	return w, nil
}

func (w *dailyWriter) rotateLocked() error {
	today := time.Now().Format("2006-01-02")
	if w.f != nil && w.day == today {
		return nil
	}
	if w.f != nil {
		_ = w.f.Close()
		w.f = nil
	}
	if w.day != "" {
		ext := filepath.Ext(w.name)
		base := strings.TrimSuffix(w.name, ext)
		old := filepath.Join(w.dir, fmt.Sprintf("%s-%s%s", base, w.day, ext))
		_ = os.Rename(w.path, old)
	}
	f, err := os.OpenFile(w.path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("打开日志文件失败: %w", err)
	}
	w.f = f
	w.day = today
	return nil
}

func (w *dailyWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if err := w.rotateLocked(); err != nil {
		return 0, err
	}
	return w.f.Write(p)
}

// Close 关闭底层文件。
func (w *dailyWriter) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.f == nil {
		return nil
	}
	err := w.f.Close()
	w.f = nil
	return err
}

// Logger 带级别过滤与按天切割的日志器。
type Logger struct {
	std    *log.Logger
	level  Level
	writer *dailyWriter
}

// New 创建日志器。
func New(file, level string) (*Logger, error) {
	w, err := newDailyWriter(file)
	if err != nil {
		return nil, err
	}
	return &Logger{
		std:    log.New(w, "", log.LstdFlags|log.Lmicroseconds),
		level:  ParseLevel(level),
		writer: w,
	}, nil
}

// Std 返回标准库 log.Logger，便于第三方组件复用同一输出。
func (l *Logger) Std() *log.Logger { return l.std }

// Writer 返回底层写入器（可交给 gin 等使用）。
func (l *Logger) Writer() io.Writer { return l.writer }

// Enabled 判断级别是否输出。
func (l *Logger) Enabled(lv Level) bool { return lv >= l.level }

func (l *Logger) logf(lv Level, format string, args ...any) {
	if l == nil || !l.Enabled(lv) {
		return
	}
	l.std.Printf(levelTag(lv)+format, args...)
}

// Debugf 调试日志。
func (l *Logger) Debugf(format string, args ...any) { l.logf(LevelDebug, format, args...) }

// Infof 普通日志。
func (l *Logger) Infof(format string, args ...any) { l.logf(LevelInfo, format, args...) }

// Warnf 警告日志。
func (l *Logger) Warnf(format string, args ...any) { l.logf(LevelWarn, format, args...) }

// Errorf 错误日志。
func (l *Logger) Errorf(format string, args ...any) { l.logf(LevelError, format, args...) }

// Close 关闭日志文件。
func (l *Logger) Close() error {
	if l == nil || l.writer == nil {
		return nil
	}
	return l.writer.Close()
}
