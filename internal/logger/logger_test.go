package logger

import (
	"context"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogger_New(t *testing.T) {
	tests := map[string]struct {
		logLevel Level
		want     slog.Level
	}{
		"DefaultLogLevel": {
			logLevel: "",
			want:     slog.LevelInfo,
		},
		"DebugLogLevel": {
			logLevel: "DEBUG",
			want:     slog.LevelDebug,
		},
		"InfoLogLevel": {
			logLevel: "INFO",
			want:     slog.LevelInfo,
		},
		"WarnLogLevel": {
			logLevel: "WARN",
			want:     slog.LevelWarn,
		},
		"ErrorLogLevel": {
			logLevel: "ERROR",
			want:     slog.LevelError,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// given
			config := LogConfig{
				Environment: "", // default
				LogLevel:    tc.logLevel,
			}
			// when
			got := New(config)
			// then
			assert.NotNil(t, got)
			enabled := got.Enabled(context.TODO(), tc.want)
			assert.Truef(t, enabled, "New() enabled log level = %t, want true for %s", enabled, tc.want)
		})
	}
}
