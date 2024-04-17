package log

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
)

type Logger interface {
	Handler() slog.Handler
	With(args ...any) Logger
	WithGroup(name string) Logger
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
	Panic(msg string, args ...any)
	PanicContext(ctx context.Context, msg string, args ...any)
}

type Config struct {
	slog.HandlerOptions
	w      io.Writer
	json   bool
	prefix string
}

type logger struct {
	*slog.Logger
}

func (l *logger) Panic(msg string, args ...any) {
	l.Error(msg, args...)
	panic(fmt.Sprint(append([]any{msg}, args...)...))
}
func (l *logger) PanicContext(_ context.Context, msg string, args ...any) {
	l.Panic(msg, args...)
}

func (l *logger) With(args ...any) Logger      { return &logger{l.Logger.With(args...)} }
func (l *logger) WithGroup(name string) Logger { return &logger{l.Logger.WithGroup(name)} }

var defaultLogger = NewDefaultLogger()

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

		defaultLogger = l
	}()

	for _, o := range opts {
		o(&config)
	}

	if config.json {
		return &logger{slog.New(slog.NewJSONHandler(config.w, &config.HandlerOptions))}
	}

	return &logger{slog.New(slog.NewTextHandler(config.w, &config.HandlerOptions))}
}

func Handler() slog.Handler                              { return defaultLogger.Handler() }
func With(args ...any) Logger                            { return defaultLogger.With(args...) }
func WithGroup(name string) Logger                       { return defaultLogger.WithGroup(name) }
func Enabled(ctx context.Context, level slog.Level) bool { return defaultLogger.Enabled(ctx, level) }
func Log(ctx context.Context, level slog.Level, msg string, args ...any) {
	defaultLogger.Log(ctx, level, msg, args...)
}
func LogAttrs(ctx context.Context, level slog.Level, msg string, attrs ...slog.Attr) {
	defaultLogger.LogAttrs(ctx, level, msg, attrs...)
}
func Debug(msg string, args ...any) { defaultLogger.Debug(msg, args...) }
func DebugContext(ctx context.Context, msg string, args ...any) {
	defaultLogger.DebugContext(ctx, msg, args...)
}
func Info(msg string, args ...any) { defaultLogger.Info(msg, args...) }
func InfoContext(ctx context.Context, msg string, args ...any) {
	defaultLogger.InfoContext(ctx, msg, args...)
}
func Warn(msg string, args ...any) { defaultLogger.Warn(msg, args...) }
func WarnContext(ctx context.Context, msg string, args ...any) {
	defaultLogger.WarnContext(ctx, msg, args...)
}
func Error(msg string, args ...any) { defaultLogger.Error(msg, args...) }
func ErrorContext(ctx context.Context, msg string, args ...any) {
	defaultLogger.ErrorContext(ctx, msg, args...)
}
func Panic(msg string, args ...any) { defaultLogger.Panic(msg, args...) }
func PanicContext(ctx context.Context, msg string, args ...any) {
	defaultLogger.PanicContext(ctx, msg, args...)
}
