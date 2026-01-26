package logx

import (
	"context"
	"io"
	"log/slog"
)

const (
	LevelTrace = slog.Level(-8)
)

func TraceWith(logger *slog.Logger, msg string, args ...any) {
	logger.Log(context.Background(), LevelTrace, msg, args...)
}

func Trace(msg string, args ...any) {
	TraceWith(slog.Default(), msg, args...)
}

var (
	level = slog.LevelInfo
)

func Setup(w io.Writer, debug, trace, quiet bool) {
	if debug {
		level = slog.LevelDebug
	}
	if trace {
		level = LevelTrace
	}
	if quiet {
		level = slog.LevelError
	}
	handler := slog.NewTextHandler(w, &slog.HandlerOptions{
		Level: level,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.LevelKey {
				switch a.Value.Any().(slog.Level) {
				case LevelTrace:
					a.Value = slog.StringValue("TRACE")
				}
			}
			if a.Key == slog.TimeKey {
				return slog.Attr{}
			}
			return a
		},
	})
	slog.SetDefault(slog.New(handler))
}
