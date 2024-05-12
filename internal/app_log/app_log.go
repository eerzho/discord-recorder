package app_log

import (
	"log/slog"
	"os"
)

const (
	debug = "debug"
	info  = "info"
	warn  = "warn"
	err   = "error"
)

var logger *slog.Logger

func Setup(logLevel string) {
	var lvl slog.Level

	switch logLevel {
	case debug:
		lvl = slog.LevelDebug
	case info:
		lvl = slog.LevelInfo
	case warn:
		lvl = slog.LevelWarn
	case err:
		lvl = slog.LevelError
	default:
		lvl = slog.LevelInfo
	}

	logger = slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: lvl}),
	)
}

func Logger() *slog.Logger {
	return logger
}
