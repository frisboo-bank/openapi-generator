package waiter

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	loggermocks "frisboo-bank/openapi-generator-service/mocks/pkg/logger"
	"frisboo-bank/openapi-generator-service/pkg/waiter/config"
	"frisboo-bank/openapi-generator-service/pkg/waiter/contracts"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setup(t *testing.T, cfg *config.WaiterOptions) contracts.Waiter {
	t.Helper()

	if cfg == nil {
		cfg = &config.WaiterOptions{}
	}
	cfg.SetDefaults()

	ctrl := gomock.NewController(t)
	logger := loggermocks.NewMockLogger(ctrl)

	logger.EXPECT().Info(gomock.Any()).AnyTimes()
	logger.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()
	logger.EXPECT().Errorf(gomock.Any(), gomock.Any()).AnyTimes()

	w, err := NewWaiter(cfg, logger)
	require.NoError(t, err)

	return w
}

func TestWaiter_Success(t *testing.T) {
	w := setup(t, nil)

	var calls atomic.Int32
	err := w.AddHook(contracts.WaiterHook{
		Name: "hook1",
		Wait: func(ctx context.Context) error { calls.Add(1); return nil },
	})
	require.NoError(t, err)

	require.NoError(t, w.Start(context.Background()))
	assert.Equal(t, int32(1), calls.Load())

	require.NoError(t, w.Start(context.Background()))
	assert.Equal(t, int32(1), calls.Load(), "hook should not run again on subsequent calls")
}

func TestWaiter_Error(t *testing.T) {
	w := setup(t, nil)

	var calls atomic.Int32
	sentinel := errors.New("boom")
	err := w.AddHook(contracts.WaiterHook{
		Name: "hook1",
		Wait: func(ctx context.Context) error { calls.Add(1); return sentinel },
	})
	require.NoError(t, err)

	err1 := w.Start(context.Background())
	err2 := w.Start(context.Background())

	assert.ErrorIs(t, err1, sentinel)
	assert.ErrorIs(t, err2, sentinel)
	assert.Equal(t, int32(1), calls.Load(), "hook should only run once even if it errors")
}

func TestWait_NoHooks(t *testing.T) {
	w := setup(t, nil)

	err := w.Start(context.Background())
	assert.NoError(t, err)
}

func TestWait_RunsWaitAndCleanup(t *testing.T) {
	w := setup(t, nil)

	var wait1, cleanup1, wait2, cleanup2 atomic.Int32

	err := w.AddHooks(
		contracts.WaiterHook{
			Name:    "hook1",
			Wait:    func(ctx context.Context) error { wait1.Add(1); return nil },
			Cleanup: func(ctx context.Context) error { cleanup1.Add(1); return nil },
		},
		contracts.WaiterHook{
			Name:    "hook2",
			Wait:    func(ctx context.Context) error { wait2.Add(1); return nil },
			Cleanup: func(ctx context.Context) error { cleanup2.Add(1); return nil },
		},
	)
	require.NoError(t, err)

	require.NoError(t, w.Start(context.Background()))

	assert.Equal(t, int32(1), wait1.Load(), "hook1 Wait should run once")
	assert.Equal(t, int32(1), cleanup1.Load(), "hook1 Cleanup should run once")
	assert.Equal(t, int32(1), wait2.Load(), "hook2 Wait should run once")
	assert.Equal(t, int32(1), cleanup2.Load(), "hook2 Cleanup should run once")
}

func TestAddHook_Validation(t *testing.T) {
	w := setup(t, nil)

	err := w.AddHook(contracts.WaiterHook{})
	assert.EqualError(t, err, "waiter: hook name is required")

	err = w.AddHook(contracts.WaiterHook{Name: "empty"})
	assert.EqualError(t, err, `waiter: hook "empty" has no Wait or Cleanup function`)

	hook := contracts.WaiterHook{
		Name:    "hook1",
		Wait:    func(ctx context.Context) error { return nil },
		Cleanup: func(ctx context.Context) error { return nil },
	}
	err = w.AddHook(hook)
	assert.NoError(t, err)

	err = w.AddHook(hook)
	assert.EqualError(t, err, `waiter: hook "hook1" already registered`)

	err = w.AddHook(contracts.WaiterHook{Name: "w", Wait: func(context.Context) error { return nil }})
	require.NoError(t, err)

	err = w.AddHook(contracts.WaiterHook{Name: "c", Cleanup: func(context.Context) error { return nil }})
	require.NoError(t, err)
}

func TestAddHooks_StopsOnError(t *testing.T) {
	w := setup(t, nil)

	err := w.AddHooks(
		contracts.WaiterHook{Name: "good", Wait: func(context.Context) error { return nil }},
		contracts.WaiterHook{Name: "bad"},
	)
	assert.EqualError(t, err, `waiter: hook "bad" has no Wait or Cleanup function`)

	err = w.AddHook(contracts.WaiterHook{Name: "good", Wait: func(context.Context) error { return nil }})
	assert.EqualError(t, err, `waiter: hook "good" already registered`)
}




func TestCancel_GracefulShutdownRunsCleanupWithLiveCtx(t *testing.T) {
	w := setup(t, nil)

	waitStarted := make(chan struct{})

	var waitCalls atomic.Int32
	var cleanupCalls atomic.Int32
	var cleanupCtxLive atomic.Bool

	require.NoError(t, w.AddHook(contracts.WaiterHook{
		Name: "hook1",
		Wait: func(ctx context.Context) error {
			waitCalls.Add(1)
			close(waitStarted)

			<-ctx.Done()
			return ctx.Err()
		},
		Cleanup: func(ctx context.Context) error {
			cleanupCalls.Add(1)
			cleanupCtxLive.Store(ctx.Err() == nil)
			return nil
		},
	}))

	waitErrCh := make(chan error, 1)
	go func() {
		waitErrCh <- w.Start(context.Background())
	}()

	<-waitStarted

	gracefulCtx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	require.NoError(t, w.Stop(gracefulCtx))

	assert.Equal(t, int32(1), waitCalls.Load(), "wait should run once")
	assert.Equal(t, int32(1), cleanupCalls.Load(), "cleanup should run once")
	assert.True(t, cleanupCtxLive.Load(), "cleanup should receive a live context during graceful shutdown")

	select {
	case err := <-waitErrCh:
		require.ErrorIs(t, err, context.Canceled)
	case <-time.After(time.Second):
		t.Fatal("Wait did not return after graceful shutdown")
	}
}

func TestCancel_ForceShutdownWhenGracefulShutdownDoesNotFinish(t *testing.T) {
	w := setup(t, nil)

	waitStarted := make(chan struct{})
	cleanupStarted := make(chan struct{})
	cleanupDone := make(chan struct{})
	cleanupRelease := make(chan struct{})

	var waitCalls atomic.Int32
	var cleanupCalls atomic.Int32

	require.NoError(t, w.AddHook(contracts.WaiterHook{
		Name: "hook1",
		Wait: func(ctx context.Context) error {
			waitCalls.Add(1)
			close(waitStarted)

			<-ctx.Done()
			return ctx.Err()
		},
		Cleanup: func(ctx context.Context) error {
			cleanupCalls.Add(1)
			close(cleanupStarted)
			defer close(cleanupDone)

			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-cleanupRelease:
				return nil
			}
		},
	}))

	waitErrCh := make(chan error, 1)
	go func() {
		waitErrCh <- w.Start(context.Background())
	}()

	<-waitStarted

	forceCtx, force := context.WithCancel(context.Background())

	cancelErrCh := make(chan error, 1)
	go func() {
		cancelErrCh <- w.Stop(forceCtx)
	}()

	// Graceful shutdown has started, and cleanup is now running.
	select {
	case <-cleanupStarted:
	case <-time.After(time.Second):
		t.Fatal("cleanup did not start during graceful shutdown")
	}

	// Consumer decides graceful shutdown is taking too long and forces it.
	force()

	select {
	case err := <-cancelErrCh:
		require.ErrorIs(t, err, ErrForcedShutdown)
		require.ErrorIs(t, err, context.Canceled)
	case <-time.After(time.Second):
		t.Fatal("Cancel did not return after force")
	}

	select {
	case err := <-waitErrCh:
		require.ErrorIs(t, err, ErrForcedShutdown)
	case <-time.After(time.Second):
		t.Fatal("Wait did not return after forced shutdown")
	}

	// Cleanup should observe the forced cancellation.
	select {
	case <-cleanupDone:
	case <-time.After(time.Second):
		t.Fatal("cleanup did not return after forced shutdown")
	}

	assert.Equal(t, int32(1), waitCalls.Load(), "wait should run once")
	assert.Equal(t, int32(1), cleanupCalls.Load(), "cleanup should start once")
}
