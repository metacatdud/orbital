package logger

import (
	"io"
	"log/slog"
	"os"
)

type Level int

// Levels ref: https://pkg.go.dev/log/slog#Level
const (
	LevelDebug Level = -4
	LevelInfo  Level = 0
	LevelWarn  Level = 4
	LevelError Level = 8
)

type Format int

const (
	FormatString Format = 0
	FormatJSON   Format = 1
)

type Logger struct {
	handler *slog.Logger
}

func (l *Logger) Debug(msg string, args ...any) {
	l.handler.Debug(msg, args...)
}

func (l *Logger) Info(msg string, args ...any) {
	l.handler.Info(msg, args...)
}

func (l *Logger) Warn(msg string, args ...any) {
	l.handler.Warn(msg, args...)
}

func (l *Logger) Error(msg string, args ...any) {
	l.handler.Error(msg, args...)
}

func New(lvl Level, fmt Format) *Logger {
	opts := parseSlogOpts(slog.Level(lvl))

	var writer io.Writer = os.Stdout

	var handler slog.Handler = slog.NewTextHandler(writer, opts)
	if fmt == FormatJSON {
		handler = slog.NewJSONHandler(writer, opts)
	}

	return &Logger{
		handler: slog.New(handler),
	}
}

func parseSlogOpts(lvl slog.Level) *slog.HandlerOptions {
	opts := &slog.HandlerOptions{
		Level: lvl,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// Handle time format
			if a.Key == slog.TimeKey {
				t := a.Value.Time()
				a.Value = slog.Int64Value(t.UnixNano() / 1000)
			}
			return a
		},
	}
	return opts
}
