package waiter

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	loggerContracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"
	"frisboo-bank/openapi-generator-service/pkg/utils"
	"frisboo-bank/openapi-generator-service/pkg/validation"
	"frisboo-bank/openapi-generator-service/pkg/waiter/contracts"
	"frisboo-bank/openapi-generator-service/pkg/waiter/models"

	"golang.org/x/sync/errgroup"
)

var _ contracts.Waiter = (*waiter)(nil)

type waiter struct {
	cancel         context.CancelFunc
	cleanupTimeout time.Duration
	ctx            context.Context
	hooks          map[string]contracts.WaiterHook
	logger         loggerContracts.Logger
	mu             sync.Mutex
	waitOnce       sync.Once
}

func NewWaiter(
	cfg *models.WaiterOptions,
	logger loggerContracts.Logger,
) (contracts.Waiter, error) {
	return NewWaiterWithContext(cfg, logger, context.Background())
}

func NewWaiterWithContext(
	cfg *models.WaiterOptions,
	logger loggerContracts.Logger,
	parentCtx context.Context,
) (contracts.Waiter, error) {
	validation.AssertNotNil("cfg", cfg)
	validation.AssertNotNil("logger", logger)

	ctx, cancel := context.WithCancel(parentCtx)

	w := &waiter{
		cancel:         cancel,
		cleanupTimeout: time.Duration(cfg.CleanupTimeoutMs) * time.Millisecond,
		ctx:            ctx,
		hooks:          make(map[string]contracts.WaiterHook),
		logger:         logger,
	}

	if cfg.CancelOnShutdownSignal {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh,
			os.Interrupt,
			syscall.SIGINT,
			syscall.SIGTERM,
			syscall.SIGQUIT)

		// The parent context is still active. When a signal arrives, cancel it.
		go func() {
			<-sigCh
			logger.Info("shutdown signal received")
			cancel()
			signal.Stop(sigCh)
			close(sigCh)
		}()
	}

	return w, nil
}

func (w *waiter) AddHooks(hooks ...contracts.WaiterHook) error {
	w.mu.Lock()
	defer w.mu.Unlock()

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

func (w *waiter) Wait() error {
	var err error
	w.waitOnce.Do(func() {
		err = w.run()
	})
	return err
}

func (w *waiter) run() error {
	w.mu.Lock()
	hooks := make(map[string]contracts.WaiterHook, len(w.hooks))
	for k, v := range w.hooks {
		hooks[k] = v
	}
	w.mu.Unlock()

	defer w.cancel()

	waitErr := w.runWait(hooks)
	cleanupErr := w.runCleanup(hooks)

	if cleanupErr != nil {
		w.logger.Errorf("hook cleanup failed with error: %v", cleanupErr)
	}

	return waitErr
}

func (w *waiter) runWait(hooks map[string]contracts.WaiterHook) error {
	group := errgroup.Group{}

	for name, hook := range hooks {
		fn := hook.Wait
		if fn == nil {
			continue
		}

		hookName := name
		waitFn := hook.Wait

		w.logger.Infof("start waiting for hook: %q", hookName)

		group.Go(func() error {
			return waitFn(w.ctx)
		})
	}

	return group.Wait()
}

func (w *waiter) runCleanup(hooks map[string]contracts.WaiterHook) error {
	group := errgroup.Group{}

	for name, hook := range hooks {
		if hook.Cleanup == nil {
			continue
		}

		hookName := name
		cleanupFn := hook.Cleanup

		w.logger.Infof("start cleaning for hook: %q", hookName)

		group.Go(func() error {
			return utils.WithTimeout(w.ctx, cleanupFn, w.cleanupTimeout)
		})
	}

	return group.Wait()
}

func (w *waiter) Cancel() {
	w.cancel()
}
