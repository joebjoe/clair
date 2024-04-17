package log

import (
	"context"
	"io"
	"log/slog"
	"os"
)

type Logger interface {
	Handler() slog.Handler
	With(args ...any) *slog.Logger
	WithGroup(name string) *slog.Logger
	Enabled(ctx context.Context, level slog.Level) bool
	Log(ctx context.Context, level slog.Level, msg string, args ...any)
	LogAttrs(ctx context.Context, level slog.Level, msg string, attrs ...slog.Attr)
	Debug(msg string, args ...any)
	DebugContext(ctx context.Context, msg string, args ...any)
	Info(msg string, args ...any)
	InfoContext(ctx context.Context, msg string, args ...any)
	Warn(msg string, args ...any)
	WarnContext(ctx context.Context, msg string, args ...any)
	Error(msg string, args ...any)
	ErrorContext(ctx context.Context, msg string, args ...any)
}

type Config struct {
	slog.HandlerOptions
	w      io.Writer
	json   bool
	prefix string
}

var logger = NewDefaultLogger()

func InitDefaultLogger()       { _ = NewDefaultLogger() }
func Init(opts ...Option)      { _ = New(opts...) }
func NewDefaultLogger() Logger { return New() }
func New(opts ...Option) (l Logger) {
	config := Config{
		HandlerOptions: slog.HandlerOptions{
			AddSource:   false,
			Level:       slog.LevelInfo,
			ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr { return a },
		},
		json: false,
		w:    os.Stderr,
	}

	defer func() {
		if l == nil {
			return
		}

		if config.prefix != "" {
			l = l.With(slog.String("prefix", config.prefix))
		}

		logger = l
	}()

	for _, o := range opts {
		o(&config)
	}

	if config.json {
		return slog.New(slog.NewJSONHandler(config.w, &config.HandlerOptions))
	}

	return slog.New(slog.NewTextHandler(config.w, &config.HandlerOptions))
}

func Handler() slog.Handler                              { return logger.Handler() }
func With(args ...any) *slog.Logger                      { return logger.With(args...) }
func WithGroup(name string) *slog.Logger                 { return logger.WithGroup(name) }
func Enabled(ctx context.Context, level slog.Level) bool { return logger.Enabled(ctx, level) }
func Log(ctx context.Context, level slog.Level, msg string, args ...any) {
	logger.Log(ctx, level, msg, args...)
}
func LogAttrs(ctx context.Context, level slog.Level, msg string, attrs ...slog.Attr) {
	logger.LogAttrs(ctx, level, msg, attrs...)
}
func Debug(msg string, args ...any) { logger.Debug(msg, args...) }
func DebugContext(ctx context.Context, msg string, args ...any) {
	logger.DebugContext(ctx, msg, args...)
}
func Info(msg string, args ...any)                             { logger.Info(msg, args...) }
func InfoContext(ctx context.Context, msg string, args ...any) { logger.InfoContext(ctx, msg, args...) }
func Warn(msg string, args ...any)                             { logger.Warn(msg, args...) }
func WarnContext(ctx context.Context, msg string, args ...any) { logger.WarnContext(ctx, msg, args...) }
func Error(msg string, args ...any)                            { logger.Error(msg, args...) }
func ErrorContext(ctx context.Context, msg string, args ...any) {
	logger.ErrorContext(ctx, msg, args...)
}
