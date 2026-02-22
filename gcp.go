package nglog

import (
	"io"
	"log/slog"
)

// NewGCPLoggingHandler sets up the Google Cloud handler.
func NewGCPLoggingHandler(w io.Writer, opts *slog.HandlerOptions) *slog.JSONHandler {
	// Create a shallow copy of the options to avoid mutating the caller's struct.
	// We want to wrap the ReplaceAttr logic without affecting the original opts.
	newOpts := *opts

	newOpts.ReplaceAttr = func(groups []string, a slog.Attr) slog.Attr {
		if opts.ReplaceAttr != nil {
			a = opts.ReplaceAttr(groups, a)
		}

		switch a.Key {
		case slog.MessageKey:
			a.Key = "message"
		case slog.SourceKey:
			a.Key = "logging.googleapis.com/sourceLocation"
		case slog.LevelKey:
			a.Key = "severity"

			level, ok := a.Value.Any().(slog.Level)
			if ok && level == LevelCritical {
				a.Value = slog.StringValue("CRITICAL")
			}
		case "trace":
			a.Key = "logging.googleapis.com/trace"
		}

		return a
	}

	return slog.NewJSONHandler(w, &newOpts)
}
