package utils

import (
	"context"
	"time"
)

func WithTimeout(parentCtx context.Context, fn func(context.Context) error, timeout time.Duration) error {
	if timeout <= 0 {
		return fn(parentCtx)
	}

	ctx, cancel := context.WithTimeout(parentCtx, timeout)
	defer cancel()

	return fn(ctx)
}
