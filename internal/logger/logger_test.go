package logger

import (
	"context"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogger_New(t *testing.T) {
	// given
	config := LogConfig{
		Environment: "", // default
		LogLevel:    "", // default
	}
	// when
	got := New(config)
	// then
	assert.NotNil(t, got)
	assert.True(t, got.Enabled(context.TODO(), slog.LevelInfo))
}
