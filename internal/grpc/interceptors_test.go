package grpc

import (
	"bytes"
	"context"
	"log/slog"
	"testing"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/igor-baiborodine/campsite-booking-go/internal/logger"
	"github.com/stretchr/testify/assert"
)

func TestInterceptors_InterceptorLogger(t *testing.T) {
	// given
	var buf bytes.Buffer
	slog.SetDefault(logger.NewDefault(&buf, nil))
	// when
	interceptorLogger().Log(context.TODO(), logging.LevelInfo, "Test message")
	// then
	got := buf.String()
	want := slog.LevelInfo.String() + " Test message\n"
	assert.Contains(t, got, want,
		"interceptorLogger.Log() got = %s, want %s", got, want)
}
