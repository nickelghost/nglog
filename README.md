# nglog

Common Go logging helpers built on top of the standard library's `log/slog` package. Provides utilities for level parsing and structured log formatting, with first-class support for Google Cloud Logging.

## Installation

```sh
go get github.com/nickelghost/nglog
```

## Usage

### Setting up the default logger

`SetUpLogger` configures the global `slog` default logger. It accepts an `io.Writer`, a format string, and a log level. Source location is always included.

Supported formats:

| Format        | Description                                                        |
| ------------- | ------------------------------------------------------------------ |
| `"gcp"`       | JSON with field names remapped to Google Cloud Logging conventions |
| `"json"`      | Standard `slog` JSON output                                        |
| anything else | Human-readable text output (default)                               |

```go
import (
    "os"
    "github.com/nickelghost/nglog"
)

func main() {
    lvl := nglog.GetLogLevel(os.Getenv("LOG_LEVEL"))
    nglog.SetUpLogger(os.Stdout, os.Getenv("LOG_FORMAT"), lvl)

    slog.Info("server starting", "port", 8080)
}
```

### Parsing log levels from configuration

`GetLogLevel` converts a string (e.g. from an environment variable or config file) into an `slog.Level`. It is case-insensitive and falls back to `slog.LevelInfo` for unrecognised values.

Supported values: `debug`, `info`, `warn`, `warning`, `err`, `error`, `crit`, `critical`.

```go
lvl := nglog.GetLogLevel("warn") // returns slog.LevelWarn
```

### Critical log level

The package defines `LevelCritical` (`slog.Level(12)`), a severity above `slog.LevelError`, for use when a critical condition must be distinguished from a regular error.

```go
slog.Log(context.Background(), nglog.LevelCritical, "unrecoverable failure")
```

### Using the GCP handler directly

`NewGCPLoggingHandler` returns an `*slog.JSONHandler` that remaps field names to match [Google Cloud Logging structured logging](https://cloud.google.com/logging/docs/structured-logging) expectations:

| slog field | GCP field                               |
| ---------- | --------------------------------------- |
| `msg`      | `message`                               |
| `source`   | `logging.googleapis.com/sourceLocation` |
| `level`    | `severity`                              |
| `trace`    | `logging.googleapis.com/trace`          |

```go
handler := nglog.NewGCPLoggingHandler(os.Stdout, &slog.HandlerOptions{
    Level:     slog.LevelDebug,
    AddSource: true,
})
logger := slog.New(handler)
logger.Info("request handled", "trace", "projects/my-project/traces/abc123")
```
