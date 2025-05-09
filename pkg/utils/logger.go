package utils

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var logLevelMap = map[string]slog.Level{
	"debug": slog.LevelDebug,
	"info":  slog.LevelInfo,
	"warn":  slog.LevelWarn,
	"error": slog.LevelError,
}

// 自定义格式处理器，以 [时间] [级别] [消息] 格式输出日志
type CustomFormatHandler struct {
	mu     sync.Mutex
	w      io.Writer
	level  slog.Level
	attrs  []slog.Attr
	groups []string
}

// NewCustomFormatHandler 创建一个新的自定义格式处理器
func NewCustomFormatHandler(w io.Writer, opts *slog.HandlerOptions) *CustomFormatHandler {
	if opts == nil {
		opts = &slog.HandlerOptions{}
	}

	// 获取日志级别，如果opts.Level是nil则默认为Info
	var level slog.Level
	if opts.Level != nil {
		level = opts.Level.Level()
	} else {
		level = slog.LevelInfo
	}

	return &CustomFormatHandler{
		w:     w,
		level: level,
	}
}

// Enabled 实现 slog.Handler 接口
func (h *CustomFormatHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= h.level
}

// Handle 实现 slog.Handler 接口，以自定义格式输出日志
func (h *CustomFormatHandler) Handle(ctx context.Context, record slog.Record) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	// 时间格式
	timeStr := record.Time.Format("2006-01-02 15:04:05.000")

	// 日志级别
	var levelStr string
	switch {
	case record.Level >= slog.LevelError:
		levelStr = "ERROR"
	case record.Level >= slog.LevelWarn:
		levelStr = "WARN "
	case record.Level >= slog.LevelInfo:
		levelStr = "INFO "
	default:
		levelStr = "DEBUG"
	}

	// 构建日志行
	logLine := fmt.Sprintf("[%s] [%s] %s", timeStr, levelStr, record.Message)

	// 处理其他属性
	var attrs []string
	record.Attrs(func(attr slog.Attr) bool {
		attrs = append(attrs, fmt.Sprintf("%s=%v", attr.Key, attr.Value))
		return true
	})

	if len(attrs) > 0 {
		logLine += " " + strings.Join(attrs, " ")
	}

	// 写入日志
	_, err := fmt.Fprintln(h.w, logLine)
	return err
}

// WithAttrs 实现 slog.Handler 接口
func (h *CustomFormatHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	h2 := *h
	h2.attrs = append(h.attrs[:], attrs...)
	return &h2
}

// WithGroup 实现 slog.Handler 接口
func (h *CustomFormatHandler) WithGroup(name string) slog.Handler {
	h2 := *h
	h2.groups = append(h.groups[:], name)
	return &h2
}

// MultiHandler 实现了 slog.Handler 接口，将日志同时发送到多个处理器
type MultiHandler struct {
	handlers []slog.Handler
}

// Enabled 实现 slog.Handler 接口
func (h *MultiHandler) Enabled(ctx context.Context, level slog.Level) bool {
	// 如果任何一个处理器启用了该级别，则返回 true
	for _, handler := range h.handlers {
		if handler.Enabled(ctx, level) {
			return true
		}
	}
	return false
}

// Handle 实现 slog.Handler 接口
func (h *MultiHandler) Handle(ctx context.Context, record slog.Record) error {
	// 将记录发送到所有处理器
	for _, handler := range h.handlers {
		if handler.Enabled(ctx, record.Level) {
			if err := handler.Handle(ctx, record); err != nil {
				return err
			}
		}
	}
	return nil
}

// WithAttrs 实现 slog.Handler 接口
func (h *MultiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newHandlers := make([]slog.Handler, len(h.handlers))
	for i, handler := range h.handlers {
		newHandlers[i] = handler.WithAttrs(attrs)
	}
	return &MultiHandler{handlers: newHandlers}
}

// WithGroup 实现 slog.Handler 接口
func (h *MultiHandler) WithGroup(name string) slog.Handler {
	newHandlers := make([]slog.Handler, len(h.handlers))
	for i, handler := range h.handlers {
		newHandlers[i] = handler.WithGroup(name)
	}
	return &MultiHandler{handlers: newHandlers}
}

// SetupLogger 设置日志输出
func SetupLogger(logLevel string, logFile string) error {
	// 创建标准错误输出的处理器，使用自定义格式
	stdHandler := NewCustomFormatHandler(os.Stderr, &slog.HandlerOptions{
		Level: logLevelMap[logLevel],
	})

	// 如果没有指定日志文件，则仅使用标准错误处理器
	if logFile == "" {
		slog.SetDefault(slog.New(stdHandler))
		return nil
	}

	// 确保日志文件所在目录存在
	logDir := filepath.Dir(logFile)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return err
	}

	// 打开日志文件，如果不存在则创建，追加写入模式
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	// 创建文件输出的处理器，使用自定义格式
	fileHandler := NewCustomFormatHandler(file, &slog.HandlerOptions{
		Level: logLevelMap[logLevel],
	})

	// 创建多输出处理器
	multiHandler := &MultiHandler{
		handlers: []slog.Handler{stdHandler, fileHandler},
	}

	// 设置全局日志处理器
	slog.SetDefault(slog.New(multiHandler))
	return nil
}

// InitDefaultLogger 初始化默认日志处理器
func InitDefaultLogger(level slog.Level) {
	slog.SetDefault(slog.New(NewCustomFormatHandler(os.Stderr, &slog.HandlerOptions{
		Level: level,
	})))
}
