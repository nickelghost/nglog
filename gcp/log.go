package ngloggcp

import (
	"context"
	"io"
	"log/slog"
)

// GetLoggerCustomSetup returns custom setup for ngtel's logger setup that adds support for GCP's log format.
func GetLoggerCustomSetup(w io.Writer, format string) func(opts *slog.HandlerOptions) bool {
	return func(opts *slog.HandlerOptions) bool {
		if format == "google_cloud" {
			slog.SetDefault(slog.New(NewCloudLoggingHandler(w, opts)))

			return true
		}

		return false
	}
}

// CloudLoggingHandler produces (s)logs that conform to Google Cloud's standards and supports its dedicated fields.
type CloudLoggingHandler struct{ handler slog.Handler }

// NewCloudLoggingHandler sets up the Google Cloud handler.
func NewCloudLoggingHandler(w io.Writer, opts *slog.HandlerOptions) *CloudLoggingHandler {
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

	return &CloudLoggingHandler{handler: slog.NewJSONHandler(w, opts)}
}

func (h *CloudLoggingHandler) Enabled(ctx context.Context, level slog.Level) bool { //nolint:revive
	return h.handler.Enabled(ctx, level)
}

func (h *CloudLoggingHandler) Handle(ctx context.Context, rec slog.Record) error { //nolint:revive
	return h.handler.Handle(ctx, rec) //nolint:wrapcheck
}

func (h *CloudLoggingHandler) WithAttrs(attrs []slog.Attr) slog.Handler { //nolint:revive
	return &CloudLoggingHandler{handler: h.handler.WithAttrs(attrs)}
}

func (h *CloudLoggingHandler) WithGroup(name string) slog.Handler { //nolint:revive
	return &CloudLoggingHandler{handler: h.handler.WithGroup(name)}
}
