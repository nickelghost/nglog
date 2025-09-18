package nglog

import (
	"context"
	"io"
	"log/slog"
)

// GCPLoggingHandler produces (s)logs that conform to Google Cloud's standards and supports its dedicated fields.
type GCPLoggingHandler struct{ handler slog.Handler }

// NewGCPLoggingHandler sets up the Google Cloud handler.
func NewGCPLoggingHandler(w io.Writer, opts *slog.HandlerOptions) *GCPLoggingHandler {
	opts.ReplaceAttr = func(_ []string, a slog.Attr) slog.Attr {
		switch a.Key {
		case slog.MessageKey:
			a.Key = "message"
		case slog.SourceKey:
			a.Key = "logging.googleapis.com/sourceLocation"
		case slog.LevelKey:
			a.Key = "severity"

			level, _ := a.Value.Any().(slog.Level)
			if level == slog.Level(12) { //nolint:mnd
				a.Value = slog.StringValue("CRITICAL")
			}
		case "trace":
			a.Key = "logging.googleapis.com/trace"
		}

		return a
	}

	return &GCPLoggingHandler{handler: slog.NewJSONHandler(w, opts)}
}

func (h *GCPLoggingHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

func (h *GCPLoggingHandler) Handle(ctx context.Context, rec slog.Record) error {
	return h.handler.Handle(ctx, rec) //nolint:wrapcheck
}

func (h *GCPLoggingHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &GCPLoggingHandler{handler: h.handler.WithAttrs(attrs)}
}

func (h *GCPLoggingHandler) WithGroup(name string) slog.Handler {
	return &GCPLoggingHandler{handler: h.handler.WithGroup(name)}
}
