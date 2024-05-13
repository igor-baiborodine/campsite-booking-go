package logger

import (
	std "log"
	"log/slog"
	"os"
)

type SilentLogger struct{}

func (*SilentLogger) Fatalf(format string, v ...interface{}) {
	std.Fatalf(format, v...)
}

func (*SilentLogger) Printf(format string, v ...interface{}) {
	std.Printf(format, v...)
}

type Level string

type LogConfig struct {
	Environment string
	LogLevel    Level
}

const (
	DEBUG Level = "DEBUG"
	INFO  Level = "INFO"
	WARN  Level = "WARN"
	ERROR Level = "ERROR"
)

func New(cfg LogConfig) *slog.Logger {
	level := new(slog.LevelVar)
	level.Set(logLevelToSlog(cfg.LogLevel))

	opts := &slog.HandlerOptions{
		Level: level,
	}
	switch cfg.Environment {
	case "production":
		return slog.New(slog.NewJSONHandler(os.Stdout, opts))
	default:
		return slog.New(slog.NewTextHandler(os.Stdout, opts))
	}
}

func logLevelToSlog(level Level) slog.Level {
	switch level {
	case ERROR:
		return slog.LevelError
	case WARN:
		return slog.LevelWarn
	case INFO:
		return slog.LevelInfo
	case DEBUG:
		return slog.LevelDebug
	default:
		return slog.LevelInfo
	}
}
