package log

import (
	"io"
	"log/slog"
	"path/filepath"
	"strings"
	"time"
)

type Option func(c *Config)

func WithJSON(c *Config) { c.json = true }

func WithLevel(lvl slog.Leveler) Option {
	return func(c *Config) { c.Level = lvl }
}

func WithLevelListener(listener func() slog.Level, tick time.Duration) Option {
	lvl := &slog.LevelVar{}

	go func() {
		wait := time.After(tick)
		lvl.Set(listener())
		<-wait
	}()

	return WithLevel(lvl)
}

func WithPrefix(prefix string) Option {
	return func(c *Config) { c.prefix = prefix }
}

func WithReplaceAttr(repl ...func(groups []string, a slog.Attr) slog.Attr) Option {
	return func(c *Config) {
		if c.ReplaceAttr == nil {
			c.ReplaceAttr = func(groups []string, a slog.Attr) slog.Attr { return a }
		}

		for _, next := range repl {
			base := c.ReplaceAttr
			c.ReplaceAttr = func(groups []string, a slog.Attr) slog.Attr {
				return next(groups, base(groups, a))
			}
		}
	}
}

func WithSource(c *Config) { c.AddSource = true }

func WithSourceDepth(n int) Option {
	return func(c *Config) {
		WithSource(c)
		WithReplaceAttr(func(_ []string, a slog.Attr) slog.Attr {
			if a.Key != slog.SourceKey {
				return a
			}

			src, ok := a.Value.Any().(slog.Source)
			if !ok {
				return a
			}

			pathParts := strings.Split(src.File, string(filepath.Separator))
			if len(pathParts) < n {
				src.File = filepath.Join(pathParts[len(pathParts)-n:]...)
				a.Value = slog.AnyValue(src)
			}

			return a
		})
	}
}

func WithTimeFormat(fmt string) Option {
	return WithReplaceAttr(func(_ []string, a slog.Attr) slog.Attr {
		if a.Key == slog.TimeKey && !a.Value.Time().IsZero() {
			a.Value = slog.StringValue(a.Value.Time().Format(fmt))
		}

		return a
	})
}

func WithWriter(w io.Writer) Option {
	return func(c *Config) { c.w = w }
}
