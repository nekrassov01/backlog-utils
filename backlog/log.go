package backlog

import (
	"io"
	"log/slog"
)

// NewLogger creates a new slog.Logger with the specified writer and log level.
func NewLogger(w io.Writer, level string) *slog.Logger {
	var lv slog.Level
	if err := lv.UnmarshalText([]byte(level)); err != nil {
		lv = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level: lv,
	}
	l := slog.New(slog.NewTextHandler(w, opts))

	return l
}
