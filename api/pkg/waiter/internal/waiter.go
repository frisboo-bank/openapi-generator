package waiter

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	loggercontracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"
	"frisboo-bank/openapi-generator-service/pkg/validation"
	"frisboo-bank/openapi-generator-service/pkg/waiter/config"
	"frisboo-bank/openapi-generator-service/pkg/waiter/contracts"

	"golang.org/x/sync/errgroup"
)

var (
	// ErrForcedShutdown is returned by Wait and Cancel when the consumer
	// forced shutdown because graceful shutdown did not complete in time.
	ErrForcedShutdown = errors.New("waiter: forced shutdown")
)

var _ contracts.Waiter = (*waiter)(nil)

type waiter struct {
	cancelCh   chan struct{}
	forceCh    chan struct{}
	finishedCh chan struct{}

	cancelOnce sync.Once
	forceOnce  sync.Once
	startOnce  sync.Once

	forced atomic.Bool

	hooks  map[string]contracts.WaiterHook
	logger loggercontracts.Logger

	mu            sync.Mutex
	cleanupCancel context.CancelFunc
	waitErr       error
}

func NewWaiter(
	cfg *config.WaiterOptions,
	logger loggercontracts.Logger,
) (contracts.Waiter, error) {
	validation.AssertNotNil("cfg", cfg)
	validation.AssertNotNil("logger", logger)

	w := &waiter{
		cancelCh: make(chan struct{}),
		forceCh:  make(chan struct{}),
		finishedCh: make(chan struct{}),
		hooks:    make(map[string]contracts.WaiterHook),
		logger:   logger,
	}

	if cfg.CancelOnShutdownSignal {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh,
			os.Interrupt,
			syscall.SIGINT,
			syscall.SIGTERM,
			syscall.SIGQUIT,
		)

		// The signal path supplies the grace deadline. A cleanup hook that
		// honours ctx (e.g. http.Server.Shutdown) drains until the deadline,
		// after which the waiter is forced.
		grace := time.Duration(cfg.CleanupTimeoutMs) * time.Millisecond

		go func() {
			<-sigCh
			logger.Info("shutdown signal received")
			graceCtx, graceCancel := context.WithTimeout(context.Background(), grace)
			_ = w.Stop(graceCtx)
			graceCancel()
			signal.Stop(sigCh)
			close(sigCh)
		}()
	}

	return w, nil
}

func (w *waiter) AddHooks(hooks ...contracts.WaiterHook) error {
	for _, h := range hooks {
		if err := w.AddHook(h); err != nil {
			return err
		}
	}
	return nil
}

func (w *waiter) AddHook(hook contracts.WaiterHook) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if hook.Name == "" {
		return fmt.Errorf("waiter: hook name is required")
	}

	if hook.Wait == nil && hook.Cleanup == nil {
		return fmt.Errorf("waiter: hook %q has no Wait or Cleanup function", hook.Name)
	}

	if _, exists := w.hooks[hook.Name]; exists {
		return fmt.Errorf("waiter: hook %q already registered", hook.Name)
	}

	w.hooks[hook.Name] = hook
	return nil
}

func (w *waiter) Start(ctx context.Context) error {
	w.start(ctx)

	select {
	case <-w.forceCh:
		return ErrForcedShutdown
	default:
	}

	select {
	case <-w.forceCh:
		return ErrForcedShutdown
	case <-w.finishedCh:
		if w.forced.Load() {
			return ErrForcedShutdown
		}
		return w.waitErr
	}
}

func (w *waiter) start(ctx context.Context) {
	w.startOnce.Do(func() {
		w.mu.Lock()
		hooks := make(map[string]contracts.WaiterHook, len(w.hooks))
		for name, hook := range w.hooks {
			hooks[name] = hook
		}
		w.mu.Unlock()

		go func() {
			w.waitErr = w.run(ctx, hooks)
			close(w.finishedCh)
		}()
	})
}

// Stop requests shutdown and blocks until it completes gracefully or ctx is
// done. If ctx finishes first the waiter is forced: cleanup contexts are
// cancelled, Start returns ErrForcedShutdown, and Stop returns an error
// wrapping ErrForcedShutdown.
func (w *waiter) Stop(ctx context.Context) error {
	w.cancelOnce.Do(func() {
		close(w.cancelCh)
	})

	// If nobody called Wait yet, Cancel can start the shutdown lifecycle.
	w.start(context.Background())

	// If already finished, prefer graceful result.
	select {
	case <-w.finishedCh:
		return w.shutdownError()
	default:
	}

	select {
	case <-w.finishedCh:
		return w.shutdownError()
	case <-ctx.Done():
		w.force()
		return fmt.Errorf("%w: %w", ErrForcedShutdown, ctx.Err())
	}
}

func (w *waiter) shutdownError() error {
	if errors.Is(w.waitErr, context.Canceled) ||
		errors.Is(w.waitErr, context.DeadlineExceeded) {
		return nil
	}

	return w.waitErr
}

func (w *waiter) force() {
	w.forceOnce.Do(func() {
		w.forced.Store(true)
		close(w.forceCh)

		w.mu.Lock()
		if w.cleanupCancel != nil {
			w.cleanupCancel()
		}
		w.mu.Unlock()
	})
}

func (w *waiter) run(
	ctx context.Context,
	hooks map[string]contracts.WaiterHook,
) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		select {
		case <-ctx.Done():
		case <-w.cancelCh:
			cancel()
		case <-w.forceCh:
			cancel()
		}
	}()

	// Cleanup runs against a fresh context so it is not pre-cancelled by the
	// wait phase. It is cancellable by force() so cooperative cleanup hooks can
	// abort during forced shutdown. The deadline is supplied by the consumer
	// via the ctx passed to Cancel; there is no per-hook timeout.
	cleanupCtx, cleanupCancel := context.WithCancel(context.Background())

	w.mu.Lock()
	w.cleanupCancel = cleanupCancel
	w.mu.Unlock()

	// If force happened before the cleanup context was registered, cancel it.
	if w.forced.Load() {
		cleanupCancel()
	}

	defer func() {
		cleanupCancel()

		w.mu.Lock()
		w.cleanupCancel = nil
		w.mu.Unlock()
	}()

	waitErr := w.runWait(ctx, hooks)

	// If forced before cleanup started, abandon cleanup.
	select {
	case <-w.forceCh:
		return waitErr
	default:
	}

	cleanupErr := w.runCleanup(cleanupCtx, hooks)
	if cleanupErr != nil {
		w.logger.Errorf("hook cleanup failed with error: %v", cleanupErr)
	}

	return waitErr
}

func (w *waiter) runWait(ctx context.Context, hooks map[string]contracts.WaiterHook) error {
	group := errgroup.Group{}

	for name, hook := range hooks {
		if hook.Wait == nil {
			continue
		}

		hookName := name
		waitFn := hook.Wait

		w.logger.Infof("start waiting for hook: %q", hookName)

		group.Go(func() error {
			return waitFn(ctx)
		})
	}

	return group.Wait()
}

func (w *waiter) runCleanup(ctx context.Context, hooks map[string]contracts.WaiterHook) error {
	group := errgroup.Group{}

	for name, hook := range hooks {
		if hook.Cleanup == nil {
			continue
		}

		hookName := name
		cleanupFn := hook.Cleanup

		w.logger.Infof("start cleaning for hook: %q", hookName)

		group.Go(func() error {
			return cleanupFn(ctx)
		})
	}

	return group.Wait()
}
