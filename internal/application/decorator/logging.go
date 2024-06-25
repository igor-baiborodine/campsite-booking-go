package decorator

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/igor-baiborodine/campsite-booking-go/internal/application/handler"
)

type loggingCommandHandler[C any] struct {
	base   handler.Command[C]
	logger *slog.Logger
}

func ApplyCommandDecorator[C any](
	handler handler.Command[C],
	logger *slog.Logger,
) handler.Command[C] {
	return loggingCommandHandler[C]{
		base:   handler,
		logger: logger,
	}
}

func (d loggingCommandHandler[C]) Handle(ctx context.Context, cmd C) (err error) {
	handlerName := extractHandlerName(cmd)
	logger := d.logger.With("command", handlerName, "command_body", fmt.Sprintf("%#v", cmd))

	logger.Debug("executing")
	defer func() {
		if err == nil {
			logger.Info("executed successfully")
		} else {
			logger.Error("failed to execute", slog.Any("error", err))
		}
	}()

	return d.base.Handle(ctx, cmd)
}

type loggingQueryHandler[C any, R any] struct {
	base   handler.Query[C, R]
	logger *slog.Logger
}

func ApplyQueryDecorator[C any, R any](
	handler handler.Query[C, R],
	logger *slog.Logger,
) handler.Query[C, R] {
	return loggingQueryHandler[C, R]{
		base:   handler,
		logger: logger,
	}
}

func (d loggingQueryHandler[C, R]) Handle(ctx context.Context, cmd C) (result R, err error) {
	handlerName := extractHandlerName(cmd)
	logger := d.logger.With("query", handlerName, "query_body", fmt.Sprintf("%#v", cmd))

	logger.Debug("executing")
	defer func() {
		if err == nil {
			logger.Info(
				"executed successfully",
				slog.Any("result", fmt.Sprintf("%v", result)),
			)
		} else {
			logger.Error("failed to execute", slog.Any("error", err))
		}
	}()

	return d.base.Handle(ctx, cmd)
}

func extractHandlerName(handler any) string {
	return strings.Split(fmt.Sprintf("%T", handler), ".")[1]
}
