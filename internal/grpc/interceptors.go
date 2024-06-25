package grpc

import (
	"context"
	"log/slog"
	"strings"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
)

const (
	grpcCode    = "grpc.code"
	grpcTimeMs  = "grpc.time_ms"
	grpcService = "grpc.service"
	grpcMethod  = "grpc.method"
)

func logServiceCalls(l *slog.Logger) logging.Logger {
	return logging.LoggerFunc(
		func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
			f := make(map[string]any, len(fields)/2)
			i := logging.Fields(fields).Iterator()

			for i.Next() {
				k, v := i.At()
				f[k] = v
			}

			if strings.Contains(msg, "finished call") {
				l.Log(
					ctx,
					slog.Level(lvl),
					msg,
					grpcService,
					f[grpcService],
					grpcMethod,
					f[grpcMethod],
					grpcTimeMs,
					f[grpcTimeMs],
					grpcCode,
					f[grpcCode],
				)
			} else {
				l.Log(ctx, slog.Level(lvl), msg, grpcService, f[grpcService], grpcMethod, f[grpcMethod],
					grpcTimeMs, f[grpcTimeMs])
			}
		},
	)
}
