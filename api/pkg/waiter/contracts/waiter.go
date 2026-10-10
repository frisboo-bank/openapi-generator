package contracts

import (
	"context"
	"errors"
)

// ErrForcedShutdown is returned by Start and Stop when the consumer forced
// shutdown because graceful shutdown did not complete in time.
var ErrForcedShutdown = errors.New("waiter: forced shutdown")

type (
	WaitFunc    func(ctx context.Context) error
	CleanupFunc func(ctx context.Context) error

	WaiterHook struct {
		Name    string
		Wait    WaitFunc
		Cleanup CleanupFunc
	}

	Waiter interface {
		AddHooks(hooks ...WaiterHook) error
		AddHook(hook WaiterHook) error

		// Start begins waiting for the application to become ready.
		// It blocks until all wait hooks finish or the context is done.
		Start(ctx context.Context) error

		// Stop requests shutdown.
		//
		// It blocks until shutdown completes gracefully or ctx is done.
		// If ctx is done first, the waiter is forced.
		Stop(ctx context.Context) error
	}
)
