package nglog

import (
	"io"
	"log/slog"
	"strings"
)

// LevelCritical is a custom log level that represents critical errors. It is set to 12, which is higher
// than the standard error level (10) in slog. This is due to slog not having a built-in critical level.
const LevelCritical = slog.Level(12)

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
	case "err", "error":
		return slog.LevelError
	case "crit", "critical":
		return LevelCritical
	}

	return slog.LevelInfo
}

// SetUpLogger sets the default logger as that of the chosen format.
// Supports Google Cloud, JSON and Text formats. Adds source information if addSource is true.
func SetUpLogger(w io.Writer, format string, lvl slog.Level, addSource bool) {
	opts := &slog.HandlerOptions{Level: lvl, AddSource: addSource}

	switch strings.ToLower(format) {
	case "gcp":
		slog.SetDefault(slog.New(NewGCPLoggingHandler(w, opts)))
	case "json":
		slog.SetDefault(slog.New(slog.NewJSONHandler(w, opts)))
	default:
		slog.SetDefault(slog.New(slog.NewTextHandler(w, opts)))
	}
}
