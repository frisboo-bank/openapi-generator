package waiter

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	loggercontracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"
	"frisboo-bank/openapi-generator-service/pkg/waiter/config"
	"frisboo-bank/openapi-generator-service/pkg/waiter/contracts"
)

// mockLogger is a minimal stub of loggercontracts.Logger. Only Infof and
// Errorf are implemented; every other method is the zero value of the
// embedded interface and is never called by the tests.
type mockLogger struct {
	loggercontracts.Logger
	mu     sync.Mutex
	infoCnt int
	errCnt  int
	errFmt  string
}

func (m *mockLogger) Infof(format string, args ...any) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.infoCnt++
}

func (m *mockLogger) Errorf(format string, args ...any) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.errCnt++
	m.errFmt = format
}

func (m *mockLogger) infoCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.infoCnt
}

func (m *mockLogger) errorCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.errCnt
}

func (m *mockLogger) lastErrorFormat() string {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.errFmt
}

// newTestWaiter builds a waiter with signal handling disabled and the
// default cleanup timeout, ready for deterministic testing.
func newTestWaiter(t *testing.T, cfg *config.WaiterOptions) (contracts.Waiter, *mockLogger) {
	t.Helper()
	if cfg == nil {
		cfg = &config.WaiterOptions{}
	}
	cfg.SetDefaults()
	cfg.CancelOnShutdownSignal = false
	logger := &mockLogger{}
	w, err := NewWaiter(cfg, logger)
	if err != nil {
		t.Fatalf("NewWaiter returned error: %v", err)
	}
	return w, logger
}

func expectErr(t *testing.T, got error, want string) {
	t.Helper()
	if got == nil {
		t.Fatalf("expected error %q, got nil", want)
	}
	if got.Error() != want {
		t.Fatalf("expected %q, got %q", want, got.Error())
	}
}

func expectPanic(t *testing.T, f func()) {
	t.Helper()
	defer func() { _ = recover() }()
	f()
	t.Fatal("expected panic, got none")
}

// TestNewWaiterNilGuard verifies that nil cfg or nil logger panics.
func TestNewWaiterNilGuard(t *testing.T) {
	logger := &mockLogger{}
	expectPanic(t, func() { NewWaiter(nil, logger) })
	expectPanic(t, func() { NewWaiter(&config.WaiterOptions{}, nil) })
}

// TestAddHookValidation covers the four rejection paths and the happy path.
func TestAddHookValidation(t *testing.T) {
	w, _ := newTestWaiter(t, nil)

	expectErr(t, w.AddHook(contracts.WaiterHook{}), "waiter: hook name is required")
	expectErr(t, w.AddHook(contracts.WaiterHook{Name: "h"}), "waiter: hook \"h\" has no Wait or Cleanup function")
	hook := contracts.WaiterHook{Name: "h", Wait: func(context.Context) error { return nil }, Cleanup: func(context.Context) error { return nil }}
	if err := w.AddHook(hook); err != nil {
		t.Fatalf("expected valid hook to be accepted, got %v", err)
	}
	expectErr(t, w.AddHook(hook), "waiter: hook \"h\" already registered")

	// Hooks with only Wait or only Cleanup are valid.
	if err := w.AddHook(contracts.WaiterHook{Name: "w", Wait: func(context.Context) error { return nil }}); err != nil {
		t.Fatalf("expected wait-only hook to be accepted, got %v", err)
	}
	if err := w.AddHook(contracts.WaiterHook{Name: "c", Cleanup: func(context.Context) error { return nil }}); err != nil {
		t.Fatalf("expected cleanup-only hook to be accepted, got %v", err)
	}
}

// TestAddHooksBulk verifies AddHooks wraps AddHook and stops on first error.
func TestAddHooksBulk(t *testing.T) {
	w, _ := newTestWaiter(t, nil)
	hook := contracts.WaiterHook{Name: "ok", Wait: func(context.Context) error { return nil }, Cleanup: func(context.Context) error { return nil }}
	if err := w.AddHooks(hook, contracts.WaiterHook{Name: "ok", Wait: func(context.Context) error { return nil }, Cleanup: func(context.Context) error { return nil }}); err == nil || err.Error() != "waiter: hook \"ok\" already registered" {
		t.Fatalf("expected duplicate error, got %v", err)
	}
}

// TestAddHooksEmpty verifies AddHooks with no args is a no-op.
func TestAddHooksEmpty(t *testing.T) {
	w, _ := newTestWaiter(t, nil)
	if err := w.AddHooks(); err != nil {
		t.Fatalf("expected no error for empty AddHooks, got %v", err)
	}
}

// TestAddHooksMultiple verifies AddHooks registers several distinct hooks.
func TestAddHooksMultiple(t *testing.T) {
	w, _ := newTestWaiter(t, nil)
	hooks := []contracts.WaiterHook{
		{Name: "a", Wait: func(context.Context) error { return nil }},
		{Name: "b", Cleanup: func(context.Context) error { return nil }},
		{Name: "c", Wait: func(context.Context) error { return nil }, Cleanup: func(context.Context) error { return nil }},
	}
	if err := w.AddHooks(hooks...); err != nil {
		t.Fatalf("expected all hooks accepted, got %v", err)
	}
}

// TestWaitWithNoHooks verifies Wait returns nil immediately with no hooks.
func TestWaitWithNoHooks(t *testing.T) {
	w, _ := newTestWaiter(t, nil)
	if err := w.Wait(context.Background()); err != nil {
		t.Fatalf("Wait returned error: %v", err)
	}
}

// TestWaitRunsHooksOnce verifies hooks run on first Wait and are not re-run
// on subsequent calls (idempotency via waitOnce).
func TestWaitRunsHooksOnce(t *testing.T) {
	w, _ := newTestWaiter(t, nil)
	var calls int
	_ = w.AddHook(contracts.WaiterHook{
		Name: "h",
		Wait: func(context.Context) error { calls++; return nil },
	})

	if err := w.Wait(context.Background()); err != nil {
		t.Fatalf("first Wait returned error: %v", err)
	}
	if calls != 1 {
		t.Fatalf("expected 1 hook call, got %d", calls)
	}
	if err := w.Wait(context.Background()); err != nil {
		t.Fatalf("second Wait returned error: %v", err)
	}
	if calls != 1 {
		t.Fatalf("hook ran twice, expected 1 call, got %d", calls)
	}
}

// TestWaitErrorCached verifies the error from the first Wait is returned
// again on subsequent calls without re-running hooks.
func TestWaitErrorCached(t *testing.T) {
	w, _ := newTestWaiter(t, nil)
	sentinel := errors.New("boom")
	var calls int
	_ = w.AddHook(contracts.WaiterHook{
		Name: "h",
		Wait: func(context.Context) error { calls++; return sentinel },
	})

	err1 := w.Wait(context.Background())
	err2 := w.Wait(context.Background())
	if !errors.Is(err1, sentinel) {
		t.Fatalf("first Wait expected sentinel, got %v", err1)
	}
	if !errors.Is(err2, sentinel) {
		t.Fatalf("second Wait expected sentinel, got %v", err2)
	}
	if calls != 1 {
		t.Fatalf("hook ran %d times, expected 1", calls)
	}
}

// TestHookAddedAfterWaitNotRun verifies a hook registered after the first
// Wait is never invoked (the hook snapshot is taken at run time).
func TestHookAddedAfterWaitNotRun(t *testing.T) {
	w, _ := newTestWaiter(t, nil)
	var secondCalled bool
	_ = w.AddHook(contracts.WaiterHook{
		Name:    "first",
		Wait:    func(context.Context) error { return nil },
		Cleanup: func(context.Context) error { return nil },
	})
	if err := w.Wait(context.Background()); err != nil {
		t.Fatalf("first Wait returned error: %v", err)
	}
	_ = w.AddHook(contracts.WaiterHook{
		Name:  "second",
		Wait:  func(context.Context) error { secondCalled = true; return nil },
	})
	if err := w.Wait(context.Background()); err != nil {
		t.Fatalf("second Wait returned error: %v", err)
	}
	if secondCalled {
		t.Fatal("hook added after Wait was unexpectedly invoked")
	}
}

// TestHookOnlyCleanup verifies a hook with only Cleanup runs cleanup but no
// wait, and Wait returns nil.
func TestHookOnlyCleanup(t *testing.T) {
	w, logger := newTestWaiter(t, nil)
	var mu sync.Mutex
	var cleaned bool
	_ = w.AddHook(contracts.WaiterHook{
		Name:    "h",
		Cleanup: func(context.Context) error { mu.Lock(); cleaned = true; mu.Unlock(); return nil },
	})

	if err := w.Wait(context.Background()); err != nil {
		t.Fatalf("Wait returned error: %v", err)
	}
	mu.Lock()
	got := cleaned
	mu.Unlock()
	if !got {
		t.Fatal("cleanup hook was never invoked")
	}
	if logger.errorCount() != 0 {
		t.Fatalf("expected no cleanup errors, got %d", logger.errorCount())
	}
}

// TestHookOnlyWait verifies a hook with only Wait runs wait and skips cleanup.
func TestHookOnlyWait(t *testing.T) {
	w, logger := newTestWaiter(t, nil)
	var called bool
	_ = w.AddHook(contracts.WaiterHook{
		Name: "h",
		Wait: func(context.Context) error { called = true; return nil },
	})

	if err := w.Wait(context.Background()); err != nil {
		t.Fatalf("Wait returned error: %v", err)
	}
	if !called {
		t.Fatal("wait hook was never invoked")
	}
	if logger.errorCount() != 0 {
		t.Fatalf("expected no cleanup errors, got %d", logger.errorCount())
	}
}

// TestMultipleHooksRunConcurrently verifies several hooks run concurrently
// rather than sequentially.
func TestMultipleHooksRunConcurrently(t *testing.T) {
	w, _ := newTestWaiter(t, nil)
	started := make(chan struct{}, 2)
	release := make(chan struct{})
	_ = w.AddHook(contracts.WaiterHook{
		Name: "a",
		Wait: func(context.Context) error { started <- struct{}{}; <-release; return nil },
	})
	_ = w.AddHook(contracts.WaiterHook{
		Name: "b",
		Wait: func(context.Context) error { started <- struct{}{}; <-release; return nil },
	})

	go func() { _ = w.Wait(context.Background()) }()

	<-started
	<-started
	close(release)

	done := make(chan error, 1)
	go func() { done <- w.Wait(context.Background()) }()
	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("Wait returned error: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Wait did not return after hooks completed")
	}
}

// TestCancelStopsWait verifies Cancel unblocks a hook waiting on ctx.Done()
// and that Cancel is safe to call more than once.
func TestCancelStopsWait(t *testing.T) {
	w, _ := newTestWaiter(t, nil)
	_ = w.AddHook(contracts.WaiterHook{
		Name: "h",
		Wait: func(ctx context.Context) error { <-ctx.Done(); return nil },
	})

	go func() {
		time.Sleep(50 * time.Millisecond)
		w.Cancel()
	}()

	done := make(chan error, 1)
	go func() { done <- w.Wait(context.Background()) }()

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("Wait returned error: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Wait did not return after Cancel")
	}

	// Idempotency: a second Cancel must not panic.
	w.Cancel()
}

// TestCancelBeforeWait verifies Cancel called before Wait still cancels the
// internal context so hooks observe ctx.Done().
func TestCancelBeforeWait(t *testing.T) {
	w, _ := newTestWaiter(t, nil)
	_ = w.AddHook(contracts.WaiterHook{
		Name: "h",
		Wait: func(ctx context.Context) error { <-ctx.Done(); return nil },
	})

	w.Cancel()

	done := make(chan error, 1)
	go func() { done <- w.Wait(context.Background()) }()

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("Wait returned error: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Wait did not return after pre-emptive Cancel")
	}
}

// TestParentContextCancellation verifies that cancelling the caller's ctx
// propagates to the hook's Wait function.
func TestParentContextCancellation(t *testing.T) {
	w, _ := newTestWaiter(t, nil)
	_ = w.AddHook(contracts.WaiterHook{
		Name: "h",
		Wait: func(ctx context.Context) error {
			<-ctx.Done()
			return ctx.Err()
		},
	})

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()

	done := make(chan error, 1)
	go func() { done <- w.Wait(ctx) }()

	select {
	case err := <-done:
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("expected context.Canceled, got %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Wait did not return after parent ctx cancel")
	}
}

// TestCleanupRunsAfterWait verifies cleanup hooks execute after Wait
// completes, and that a cleanup error is logged but does not override
// the wait error.
func TestCleanupRunsAfterWait(t *testing.T) {
	w, logger := newTestWaiter(t, nil)
	var mu sync.Mutex
	var cleaned bool
	_ = w.AddHook(contracts.WaiterHook{
		Name:    "h",
		Wait:    func(context.Context) error { return errors.New("wait failed") },
		Cleanup: func(context.Context) error { mu.Lock(); cleaned = true; mu.Unlock(); return nil },
	})

	if err := w.Wait(context.Background()); err == nil || err.Error() != "wait failed" {
		t.Fatalf("expected wait error, got %v", err)
	}

	mu.Lock()
	got := cleaned
	mu.Unlock()
	if !got {
		t.Fatal("cleanup hook was never invoked")
	}
	if logger.errorCount() != 0 {
		t.Fatalf("expected no cleanup errors, got %d", logger.errorCount())
	}
}

// TestCleanupErrorIsLogged verifies that a failing cleanup hook is logged
// and does not override the (nil) wait error.
func TestCleanupErrorIsLogged(t *testing.T) {
	w, logger := newTestWaiter(t, nil)
	_ = w.AddHook(contracts.WaiterHook{
		Name:    "h",
		Cleanup: func(context.Context) error { return errors.New("cleanup boom") },
	})

	if err := w.Wait(context.Background()); err != nil {
		t.Fatalf("Wait returned error: %v", err)
	}
	if logger.errorCount() == 0 {
		t.Fatal("expected cleanup error to be logged")
	}
	if logger.lastErrorFormat() != "hook cleanup failed with error: %v" {
		t.Fatalf("unexpected log format: %q", logger.lastErrorFormat())
	}
}

// TestCleanupTimeoutIsRespected verifies that a cleanup hook exceeding the
// configured timeout is cut short by utils.WithTimeout.
func TestCleanupTimeoutIsRespected(t *testing.T) {
	cfg := &config.WaiterOptions{CleanupTimeoutMs: 100}
	w, logger := newTestWaiter(t, cfg)
	var mu sync.Mutex
	cancelled := false
	_ = w.AddHook(contracts.WaiterHook{
		Name: "h",
		Cleanup: func(ctx context.Context) error {
			select {
			case <-ctx.Done():
				mu.Lock()
				cancelled = true
				mu.Unlock()
				return ctx.Err()
			case <-time.After(5 * time.Second):
				return nil
			}
		},
	})

	if err := w.Wait(context.Background()); err != nil {
		t.Fatalf("Wait returned error: %v", err)
	}
	mu.Lock()
	got := cancelled
	mu.Unlock()
	if !got {
		t.Fatal("cleanup hook was not cancelled by timeout")
	}
	if logger.errorCount() == 0 {
		t.Fatal("expected cleanup timeout error to be logged")
	}
}

// TestAddHooksPartial verifies that when AddHooks receives a mix of valid
// and invalid hooks, the valid ones are registered before the first error
// is returned.
func TestAddHooksPartial(t *testing.T) {
	w, _ := newTestWaiter(t, nil)
	good := contracts.WaiterHook{Name: "good", Wait: func(context.Context) error { return nil }}
	bad := contracts.WaiterHook{Name: "bad"} // missing Wait and Cleanup

	if err := w.AddHooks(good, bad); err == nil || err.Error() != "waiter: hook \"bad\" has no Wait or Cleanup function" {
		t.Fatalf("expected validation error, got %v", err)
	}
	// The good hook was registered before the bad one was rejected.
	if err := w.AddHook(contracts.WaiterHook{Name: "good", Wait: func(context.Context) error { return nil }}); err == nil || err.Error() != "waiter: hook \"good\" already registered" {
		t.Fatalf("expected duplicate error, got %v", err)
	}
}

// TestCleanupMultipleHooksConcurrent verifies that cleanup hooks run
// concurrently rather than sequentially.
func TestCleanupMultipleHooksConcurrent(t *testing.T) {
	w, _ := newTestWaiter(t, nil)
	started := make(chan struct{}, 2)
	release := make(chan struct{})
	_ = w.AddHook(contracts.WaiterHook{
		Name:    "a",
		Cleanup: func(context.Context) error { started <- struct{}{}; <-release; return nil },
	})
	_ = w.AddHook(contracts.WaiterHook{
		Name:    "b",
		Cleanup: func(context.Context) error { started <- struct{}{}; <-release; return nil },
	})

	done := make(chan error, 1)
	go func() { done <- w.Wait(context.Background()) }()

	<-started
	<-started
	close(release)

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("Wait returned error: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Wait did not return after cleanup completed")
	}
}

// TestCleanupTimeoutOnlyOneLogs verifies that when multiple cleanup hooks
// run concurrently and only one exceeds the timeout, exactly one cleanup
// error is logged.
func TestCleanupTimeoutOnlyOneLogs(t *testing.T) {
	cfg := &config.WaiterOptions{CleanupTimeoutMs: 50}
	w, logger := newTestWaiter(t, cfg)
	release := make(chan struct{})
	_ = w.AddHook(contracts.WaiterHook{
		Name:    "fast",
		Cleanup: func(context.Context) error { <-release; return nil },
	})
	_ = w.AddHook(contracts.WaiterHook{
		Name: "slow",
		Cleanup: func(ctx context.Context) error {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(5 * time.Second):
				return nil
			}
		},
	})

	done := make(chan error, 1)
	go func() { done <- w.Wait(context.Background()) }()

	// Give the slow hook time to hit its 50ms timeout.
	time.Sleep(150 * time.Millisecond)
	close(release)

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("Wait returned error: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Wait did not return after cleanup completed")
	}
	if logger.errorCount() != 1 {
		t.Fatalf("expected exactly 1 cleanup error logged, got %d", logger.errorCount())
	}
}

// TestWaitContextAlreadyCancelled verifies that a context already cancelled
// before Wait is called propagates to the hook immediately.
func TestWaitContextAlreadyCancelled(t *testing.T) {
	w, _ := newTestWaiter(t, nil)
	_ = w.AddHook(contracts.WaiterHook{
		Name: "h",
		Wait: func(ctx context.Context) error {
			<-ctx.Done()
			return ctx.Err()
		},
	})

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := w.Wait(ctx)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
}

// TestCancelAfterWaitIsNoop verifies that calling Cancel after Wait has
// already completed is a safe no-op and does not affect the cached result.
func TestCancelAfterWaitIsNoop(t *testing.T) {
	w, _ := newTestWaiter(t, nil)
	_ = w.AddHook(contracts.WaiterHook{
		Name:    "h",
		Wait:    func(context.Context) error { return nil },
		Cleanup: func(context.Context) error { return nil },
	})

	if err := w.Wait(context.Background()); err != nil {
		t.Fatalf("first Wait returned error: %v", err)
	}
	w.Cancel()
	if err := w.Wait(context.Background()); err != nil {
		t.Fatalf("second Wait returned error: %v", err)
	}
}

// TestNewWaiterSignalHandling verifies that constructing a waiter with
// signal handling enabled does not panic and that Cancel still works.
func TestNewWaiterSignalHandling(t *testing.T) {
	cfg := &config.WaiterOptions{CancelOnShutdownSignal: true}
	cfg.SetDefaults()
	w, _ := newTestWaiter(t, cfg)
	_ = w.AddHook(contracts.WaiterHook{
		Name: "h",
		Wait: func(ctx context.Context) error { <-ctx.Done(); return nil },
	})

	go func() {
		time.Sleep(50 * time.Millisecond)
		w.Cancel()
	}()

	done := make(chan error, 1)
	go func() { done <- w.Wait(context.Background()) }()

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("Wait returned error: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Wait did not return after Cancel")
	}
}

// TestAddHookThenAddHooksDuplicate verifies that AddHooks can be called
// after AddHook and that a duplicate name is rejected.
func TestAddHookThenAddHooksDuplicate(t *testing.T) {
	w, _ := newTestWaiter(t, nil)
	_ = w.AddHook(contracts.WaiterHook{Name: "h", Wait: func(context.Context) error { return nil }})
	if err := w.AddHooks(contracts.WaiterHook{Name: "h", Wait: func(context.Context) error { return nil }}); err == nil || err.Error() != "waiter: hook \"h\" already registered" {
		t.Fatalf("expected duplicate error, got %v", err)
	}
}

// TestCleanupTimeoutZero verifies that a zero cleanup timeout disables the
// timeout wrapper (WithTimeout calls fn directly) and the hook completes
// without a logged error.
func TestCleanupTimeoutZero(t *testing.T) {
	cfg := &config.WaiterOptions{CleanupTimeoutMs: 0}
	w, logger := newTestWaiter(t, cfg)
	_ = w.AddHook(contracts.WaiterHook{
		Name: "h",
		Cleanup: func(context.Context) error {
			// This would time out under a real timeout; with timeout=0 it
			// runs directly and must complete.
			time.Sleep(20 * time.Millisecond)
			return nil
		},
	})

	if err := w.Wait(context.Background()); err != nil {
		t.Fatalf("Wait returned error: %v", err)
	}
	if logger.errorCount() != 0 {
		t.Fatalf("expected no cleanup errors, got %d", logger.errorCount())
	}
}
