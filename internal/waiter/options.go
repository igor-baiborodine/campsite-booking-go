package waiter

import (
	"context"
)

type Option func(c *waiterCfg)

func ParentContext(ctx context.Context) Option {
	return func(c *waiterCfg) {
		c.parentCtx = ctx
	}
}

func CatchSignals() Option {
	return func(c *waiterCfg) {
		c.catchSignals = true
	}
}
