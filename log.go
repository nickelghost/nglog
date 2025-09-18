package nglog

import (
	"io"
	"log/slog"
	"strings"
)

// GetLogLevel determines the slog log level based on the provided string.
// It returns slog.LevelInfo if the provided string does not match any known level.
// The function is case-insensitive and supports the following levels: debug, info, warn(ing), error.
func GetLogLevel(level string) slog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	}

	return slog.LevelInfo
}

// SetUpLogger sets the default logger as that of the chosen format.
// It supports Google Cloud, JSON and Text formats.
func SetUpLogger(w io.Writer, format string, lvl slog.Level) {
	opts := &slog.HandlerOptions{Level: lvl, AddSource: true}

	switch strings.ToLower(format) {
	case "gcp":
		slog.SetDefault(slog.New(NewGCPLoggingHandler(w, opts)))
	case "json":
		slog.SetDefault(slog.New(slog.NewJSONHandler(w, opts)))
	default:
		slog.SetDefault(slog.New(slog.NewTextHandler(w, opts)))
	}
}
