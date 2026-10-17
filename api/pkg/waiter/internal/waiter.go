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
	"frisboo-bank/openapi-generator-service/pkg/validation"
	"frisboo-bank/openapi-generator-service/pkg/waiter/config"
	"frisboo-bank/openapi-generator-service/pkg/waiter/contracts"
	"frisboo-bank/openapi-generator-service/pkg/waiter/types/hookstatus"

	"golang.org/x/sync/errgroup"
)

var _ contracts.Waiter = (*waiter)(nil)

type waiter struct {
	stopCh   chan struct{}
	stopOnce sync.Once

	hooks  map[string]contracts.WaiterHook
	logger loggercontracts.Logger

	mu sync.Mutex
}

func NewWaiter(
	cfg *config.WaiterOptions,
	logger loggercontracts.Logger,
) (contracts.Waiter, error) {
	validation.AssertNotNil("cfg", cfg)
	validation.AssertNotNil("logger", logger)

	w := &waiter{
		stopCh: make(chan struct{}),
		hooks:  make(map[string]contracts.WaiterHook),
		logger: logger,
	}

	if cfg.CancelOnShutdownSignal {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

		grace := time.Duration(cfg.CleanupTimeoutMs) * time.Millisecond

		go func() {
			<-sigCh
			logger.Info("shutdown signal received, starting graceful shutdown")

			done := make(chan struct{})
			go func() {
				select {
				case <-sigCh:
					logger.Info("shutdown signal received again, forcing exit")
					os.Exit(1)
				case <-done:
					return
				}
			}()

			graceCtx, graceCancel := context.WithTimeout(context.Background(), grace)
			if err := w.Stop(graceCtx); err != nil {
				logger.Errorf("graceful shutdown failed: %v", err)
			}
			graceCancel()

			signal.Stop(sigCh)
			close(done)
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
	validation.AssertNotNil("hook", hook)
	validation.AssertNotEmpty("hook.Name", hook.Name)

	w.mu.Lock()
	defer w.mu.Unlock()

	if hook.Start == nil && hook.Stop == nil {
		return fmt.Errorf("waiter: hook %q has no Start or Stop function", hook.Name)
	}

	if _, exists := w.hooks[hook.Name]; exists {
		return fmt.Errorf("waiter: hook %q already registered", hook.Name)
	}

	w.hooks[hook.Name] = hook
	return nil
}

func (w *waiter) Start(ctx context.Context) error {
	waitCtx, waitCancel := context.WithCancel(ctx)
	defer waitCancel()

	go func() {
		select {
		case <-w.stopCh:
			waitCancel()
		case <-waitCtx.Done():
			return
		}
	}()

	w.mu.Lock()
	hooks := make(map[string]contracts.WaiterHook, len(w.hooks))
	for name, hook := range w.hooks {
		hooks[name] = hook
	}
	w.mu.Unlock()

	group, groupCtx := errgroup.WithContext(waitCtx)

	for name, hook := range hooks {
		if hook.Start == nil {
			continue
		}

		w.logger.Infof("starting hook: %q", name)

		group.Go(func() error {
			if err := hook.Start(groupCtx); err != nil {
				w.logger.Errorf("hook %q start failed: %v", name, err)
				return err
			}
			return nil
		})
	}

	return group.Wait()
}

func (w *waiter) Stop(ctx context.Context) error {
	w.stopOnce.Do(func() {
		close(w.stopCh)
	})

	w.mu.Lock()
	hooks := make(map[string]contracts.WaiterHook, len(w.hooks))
	for name, hook := range w.hooks {
		hooks[name] = hook
	}
	w.mu.Unlock()

	group, groupCtx := errgroup.WithContext(ctx)

	for name, hook := range hooks {
		if hook.Stop == nil {
			continue
		}

		if !hook.IsRunning() {
			w.logger.Infof("hook %q is not running, skipping stop", name)
			continue
		}

		w.logger.Infof("stopping hook: %q", name)

		group.Go(func() error {
			if err := hook.Stop(groupCtx); err != nil {
				w.logger.Errorf("hook %q stop failed: %v", name, err)
				return err
			}
			return nil
		})
	}

	return group.Wait()
}

func (w *waiter) Status() map[string]hookstatus.HookStatus {
	w.mu.Lock()
	defer w.mu.Unlock()

	statuses := make(map[string]hookstatus.HookStatus)

	for name, hook := range w.hooks {
		if hook.Status != nil {
			statuses[name] = hook.Status()
		}
	}

	return statuses
}
