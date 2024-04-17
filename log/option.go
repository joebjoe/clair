package log

import (
	"io"
	"log/slog"
	"path"
	"path/filepath"
	"regexp"
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

var pkgRex = regexp.MustCompile(`pkg/mod/(.*)$`)

func WithSource(c *Config) {
	c.AddSource = true
	WithReplaceAttr(func(groups []string, a slog.Attr) slog.Attr {
		if a.Key != slog.SourceKey {
			return a
		}

		slogSrc := a.Value.Any().(*slog.Source)
		pkg, file := filepath.Split(slogSrc.File)
		pkg = strings.TrimPrefix(strings.TrimSuffix(pkg, string(filepath.Separator)), string(filepath.Separator))
		src := source{
			Directory: pkg,
			File:      file,
			Function:  path.Base(slogSrc.Function),
			Line:      slogSrc.Line,
		}

		// file replacing
		matches := pkgRex.FindStringSubmatch(src.Directory)
		if matches != nil {
			src.Directory = matches[1]
		}

		a.Value = slog.AnyValue(src)

		return a
	})(c)
}

func WithTimeFormat(format string) Option {
	return WithReplaceAttr(func(_ []string, a slog.Attr) slog.Attr {
		if a.Key == slog.TimeKey && !a.Value.Time().IsZero() {
			a.Value = slog.StringValue(a.Value.Time().Format(format))
		}

		return a
	})
}

func WithWriter(w io.Writer) Option {
	return func(c *Config) { c.w = w }
}

type source struct {
	File      string `json:"file"`
	Function  string `json:"function"`
	Directory string `json:"directory"`
	Line      int    `json:"line"`
}
