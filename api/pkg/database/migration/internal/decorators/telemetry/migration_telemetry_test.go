package telemetry

import (
	"context"
	"errors"
	"testing"

	migrationtype "frisboo-bank/openapi-generator-service/pkg/database/migration/types/migrationtype"
	environmentenum "frisboo-bank/openapi-generator-service/pkg/environment/models/enums/environment"
	loggercontracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"
	"frisboo-bank/openapi-generator-service/pkg/logger"
	"frisboo-bank/openapi-generator-service/pkg/telemetry/metrics"
	"frisboo-bank/openapi-generator-service/pkg/telemetry/tracer"
	"github.com/stretchr/testify/require"
)

type fakeAdapter struct {
	upErr   error
	downErr error
	resErr  error
	stErr   error
	curErr  error
	cur     int64
	calls   []string
}

func (f *fakeAdapter) Up(ctx context.Context, version uint) error {
	f.calls = append(f.calls, "up")
	return f.upErr
}
func (f *fakeAdapter) Down(ctx context.Context, version uint) error {
	f.calls = append(f.calls, "down")
	return f.downErr
}
func (f *fakeAdapter) Reset(ctx context.Context) error {
	f.calls = append(f.calls, "reset")
	return f.resErr
}
func (f *fakeAdapter) Status(ctx context.Context) error {
	f.calls = append(f.calls, "status")
	return f.stErr
}
func (f *fakeAdapter) CurrentVersion(ctx context.Context) (int64, error) {
	f.calls = append(f.calls, "current_version")
	return f.cur, f.curErr
}
func (f *fakeAdapter) Name() string                      { return "fake" }
func (f *fakeAdapter) Type() migrationtype.MigrationType { return migrationtype.MigrationTypes.GOOSE }
func (f *fakeAdapter) Logger() loggercontracts.Logger {
	return logger.CreateNoopLogger("fake", environmentenum.Environments.TESTING)
}

func TestWrapMigrationAdapterForTracing(t *testing.T) {
	t.Run("tracer is required", func(t *testing.T) {
		fa := &fakeAdapter{}
		require.Panics(t, func() {
			_ = WrapMigrationAdapterForTracing("m", fa, nil)
		})
	})

	t.Run("emits spans for every operation", func(t *testing.T) {
		fa := &fakeAdapter{cur: 7}
		tr, err := tracer.CreateNoopTracer("m", logger.CreateNoopLogger("test", environmentenum.Environments.TESTING))
		require.NoError(t, err)

		dec := WrapMigrationAdapterForTracing("m", fa, tr)
		require.NotNil(t, dec)

		ctx := context.Background()
		require.NoError(t, dec.Up(ctx, 0))
		require.NoError(t, dec.Down(ctx, 1))
		require.NoError(t, dec.Reset(ctx))
		require.NoError(t, dec.Status(ctx))
		v, err := dec.CurrentVersion(ctx)
		require.NoError(t, err)
		require.Equal(t, int64(7), v)
		require.Equal(t, []string{"up", "down", "reset", "status", "current_version"}, fa.calls)
	})

	t.Run("errors propagate and are recorded on the span", func(t *testing.T) {
		fa := &fakeAdapter{upErr: errors.New("boom")}
		tr, err := tracer.CreateNoopTracer("m", logger.CreateNoopLogger("test", environmentenum.Environments.TESTING))
		require.NoError(t, err)

		dec := WrapMigrationAdapterForTracing("m", fa, tr)
		require.Error(t, dec.Up(context.Background(), 0))
	})

	t.Run("promotes Name/Type/Logger from the embedded adapter", func(t *testing.T) {
		fa := &fakeAdapter{}
		tr, _ := tracer.CreateNoopTracer("m", logger.CreateNoopLogger("test", environmentenum.Environments.TESTING))
		dec := WrapMigrationAdapterForTracing("m", fa, tr)
		require.Equal(t, "fake", dec.Name())
		require.Equal(t, migrationtype.MigrationTypes.GOOSE, dec.Type())
	})
}

func TestWrapMigrationAdapterForMetrics(t *testing.T) {
	t.Run("metrics is required", func(t *testing.T) {
		fa := &fakeAdapter{}
		require.Panics(t, func() {
			_ = WrapMigrationAdapterForMetrics("m", fa, nil)
		})
	})

	t.Run("records duration for every operation", func(t *testing.T) {
		fa := &fakeAdapter{cur: 7}
		mt, err := metrics.CreateNoopMetrics("m", logger.CreateNoopLogger("test", environmentenum.Environments.TESTING))
		require.NoError(t, err)

		dec := WrapMigrationAdapterForMetrics("m", fa, mt)
		require.NotNil(t, dec)

		ctx := context.Background()
		require.NoError(t, dec.Up(ctx, 0))
		require.NoError(t, dec.Down(ctx, 1))
		require.NoError(t, dec.Reset(ctx))
		require.NoError(t, dec.Status(ctx))
		v, err := dec.CurrentVersion(ctx)
		require.NoError(t, err)
		require.Equal(t, int64(7), v)
		require.Equal(t, []string{"up", "down", "reset", "status", "current_version"}, fa.calls)
	})

	t.Run("errors propagate", func(t *testing.T) {
		fa := &fakeAdapter{upErr: errors.New("boom")}
		mt, err := metrics.CreateNoopMetrics("m", logger.CreateNoopLogger("test", environmentenum.Environments.TESTING))
		require.NoError(t, err)

		dec := WrapMigrationAdapterForMetrics("m", fa, mt)
		require.Error(t, dec.Up(context.Background(), 0))
	})

	t.Run("promotes Name/Type/Logger from the embedded adapter", func(t *testing.T) {
		fa := &fakeAdapter{}
		mt, _ := metrics.CreateNoopMetrics("m", logger.CreateNoopLogger("test", environmentenum.Environments.TESTING))
		dec := WrapMigrationAdapterForMetrics("m", fa, mt)
		require.Equal(t, "fake", dec.Name())
		require.Equal(t, migrationtype.MigrationTypes.GOOSE, dec.Type())
	})
}
