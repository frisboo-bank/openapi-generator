package waiter_test

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"

	loggermocks "frisboo-bank/openapi-generator-service/mocks/pkg/logger"
	"frisboo-bank/openapi-generator-service/pkg/waiter/config"
	"frisboo-bank/openapi-generator-service/pkg/waiter/contracts"
	waiter "frisboo-bank/openapi-generator-service/pkg/waiter/internal"

	"github.com/go-jose/go-jose/v4/testutils/assert"
	"github.com/go-jose/go-jose/v4/testutils/require"
	"github.com/golang/mock/gomock"
)

func newWaiter(t *testing.T, cfg *config.WaiterOptions) contracts.Waiter {
	t.Helper()

	if cfg == nil {
		cfg = &config.WaiterOptions{}
	}
	cfg.SetDefaults()

	ctrl := gomock.NewController(t)
	logger := loggermocks.NewMockLogger(ctrl)

	logger.EXPECT().Info(gomock.Any()).AnyTimes()
	logger.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()

	w, err := waiter.NewWaiter(cfg, logger)
	require.NoError(t, err)

	return w
}

func TestWaiter_Success(t *testing.T) {
	instance := newWaiter(t, nil)

	var calls atomic.Int32

	err := instance.AddHook(contracts.WaiterHook{
		Name: "hook1",
		Wait: func(ctx context.Context) error {
			calls.Add(1)
			return nil
		},
	})
	assert.NoError(t, err)

	err = instance.Wait(context.Background())
	assert.NoError(t, err)

	assert.Equal(t, int32(1), calls.Load())
}

func TestWaiter_Error(t *testing.T) {
	instance := newWaiter(t, nil)

	sentinel := errors.New("boom")

	err := instance.AddHook(contracts.WaiterHook{
		Name: "hook1",
		Wait: func(ctx context.Context) error {
			return sentinel
		},
	})
	assert.NoError(t, err)

	err = instance.Wait(context.Background())

	assert.Error(t, err)
}
