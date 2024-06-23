package handler

import (
	"context"
)

type Command[C any] interface {
	Handle(ctx context.Context, cmd C) error
}

type Query[Q any, R any] interface {
	Handle(ctx context.Context, qry Q) (R, error)
}
