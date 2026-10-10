package waiter

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	loggercontracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"
	"frisboo-bank/openapi-generator-service/pkg/utils"
	"frisboo-bank/openapi-generator-service/pkg/validation"
	"frisboo-bank/openapi-generator-service/pkg/waiter/config"
	"frisboo-bank/openapi-generator-service/pkg/waiter/contracts"

	"golang.org/x/sync/errgroup"
)

var _ contracts.Waiter = (*waiter)(nil)

type waiter struct {
	cancelCh       chan struct{}
	cancelOnce     sync.Once
	cleanupTimeout time.Duration
	hooks          map[string]contracts.WaiterHook
	logger         loggercontracts.Logger
	mu             sync.Mutex
	waitErr        error
	waitOnce       sync.Once
}

func NewWaiter(
	cfg *config.WaiterOptions,
	logger loggercontracts.Logger,
) (contracts.Waiter, error) {
	validation.AssertNotNil("cfg", cfg)
	validation.AssertNotNil("logger", logger)

	w := &waiter{
		cancelCh:       make(chan struct{}),
		cleanupTimeout: time.Duration(cfg.CleanupTimeoutMs) * time.Millisecond,
		hooks:          make(map[string]contracts.WaiterHook),
		logger:         logger,
	}

	if cfg.CancelOnShutdownSignal {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh,
			os.Interrupt,
			syscall.SIGINT,
			syscall.SIGTERM,
			syscall.SIGQUIT,
		)

		go func() {
			<-sigCh
			logger.Info("shutdown signal received")
			w.Cancel()
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

func (w *waiter) Wait(ctx context.Context) error {
	w.waitOnce.Do(func() {
		w.waitErr = w.run(ctx)
	})
	return w.waitErr
}

func (w *waiter) run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		select {
		case <-ctx.Done():
		case <-w.cancelCh:
			cancel()
		}
	}()

	w.mu.Lock()
	hooks := make(map[string]contracts.WaiterHook, len(w.hooks))
	for k, v := range w.hooks {
		hooks[k] = v
	}
	w.mu.Unlock()

	waitErr := w.runWait(ctx, hooks)
	cleanupErr := w.runCleanup(ctx, hooks)

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
			return utils.WithTimeout(ctx, cleanupFn, w.cleanupTimeout)
		})
	}

	return group.Wait()
}

func (w *waiter) Cancel() {
	w.cancelOnce.Do(func() {
		close(w.cancelCh)
	})
}
