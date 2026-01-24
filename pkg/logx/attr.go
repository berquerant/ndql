package logx

import (
	"encoding/json"
	"fmt"
	"log/slog"
)

func Err(err error) any {
	if err != nil {
		return slog.String("err", err.Error())
	}
	return slog.String("err", "nil")
}

func JSON(k string, v any) any {
	b, err := json.Marshal(v)
	if err == nil {
		return slog.String(k, string(b))
	}
	return slog.String(k, "not a json")
}

func Format(k string, v any, format string) any { return slog.String(k, fmt.Sprintf(format, v)) }
func String(k string, v any) any                { return Format(k, v, "%s") }
func Value(k string, v any) any                 { return Format(k, v, "%v") }
func Verbose(k string, v any) any               { return Format(k, v, "%#v") }

func OnTrace(f func()) {
	if level == LevelTrace {
		f()
	}
}

func Error(err error, msg string, args ...any) { slog.Error(msg, append(args, Err(err))...) }
