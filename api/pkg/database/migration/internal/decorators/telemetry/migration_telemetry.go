package telemetry

import (
	"context"
	"time"

	"frisboo-bank/openapi-generator-service/pkg/database/migration/contracts"
	metricscontracts "frisboo-bank/openapi-generator-service/pkg/telemetry/metrics/contracts"
	tracercontracts "frisboo-bank/openapi-generator-service/pkg/telemetry/tracer/contracts"
	"frisboo-bank/openapi-generator-service/pkg/validation"
)

// migrationAdapterTracing decorates a MigrationAdapter with tracing spans.
// Only call WrapMigrationAdapterForTracing when a non-nil tracer is available;
// otherwise leave the adapter unwrapped. No noop fallback is needed.
var _ contracts.MigrationAdapter = (*migrationAdapterTracing)(nil)

type migrationAdapterTracing struct {
	contracts.MigrationAdapter
	name   string
	tracer tracercontracts.Tracer
}

// WrapMigrationAdapterForTracing wraps a MigrationAdapter so that every
// operation emits a tracing span. The tracer must be non-nil; callers that
// have no tracer simply don't call this wrapper.
func WrapMigrationAdapterForTracing(
	name string,
	delegate contracts.MigrationAdapter,
	tracer tracercontracts.Tracer,
) contracts.MigrationAdapter {
	validation.AssertNotEmpty("name", name)
	validation.AssertNotNil("delegate", delegate)
	validation.AssertNotNil("tracer", tracer)

	return &migrationAdapterTracing{
		MigrationAdapter: delegate,
		name:             name,
		tracer:           tracer,
	}
}

func (m *migrationAdapterTracing) record(ctx context.Context, op string, fn func(context.Context) error) error {
	ctx, span := m.tracer.Start(ctx, "migration."+op)
	defer span.End()

	err := fn(ctx)
	if err != nil {
		span.RecordError(err)
	}
	return err
}

func (m *migrationAdapterTracing) Up(ctx context.Context, version uint) error {
	return m.record(ctx, "up", func(ctx context.Context) error {
		return m.MigrationAdapter.Up(ctx, version)
	})
}

func (m *migrationAdapterTracing) Down(ctx context.Context, version uint) error {
	return m.record(ctx, "down", func(ctx context.Context) error {
		return m.MigrationAdapter.Down(ctx, version)
	})
}

func (m *migrationAdapterTracing) Reset(ctx context.Context) error {
	return m.record(ctx, "reset", func(ctx context.Context) error {
		return m.MigrationAdapter.Reset(ctx)
	})
}

func (m *migrationAdapterTracing) Status(ctx context.Context) error {
	return m.record(ctx, "status", func(ctx context.Context) error {
		return m.MigrationAdapter.Status(ctx)
	})
}

func (m *migrationAdapterTracing) CurrentVersion(ctx context.Context) (int64, error) {
	ctx, span := m.tracer.Start(ctx, "migration.current_version")
	defer span.End()

	version, err := m.MigrationAdapter.CurrentVersion(ctx)
	if err != nil {
		span.RecordError(err)
	}
	return version, err
}

// migrationAdapterMetrics decorates a MigrationAdapter with duration metrics.
// Only call WrapMigrationAdapterForMetrics when a non-nil metrics is available;
// otherwise leave the adapter unwrapped. No noop fallback is needed.
var _ contracts.MigrationAdapter = (*migrationAdapterMetrics)(nil)

type migrationAdapterMetrics struct {
	contracts.MigrationAdapter
	name    string
	metrics metricscontracts.Metrics
}

// WrapMigrationAdapterForMetrics wraps a MigrationAdapter so that every
// operation records a duration metric. The metrics must be non-nil; callers
// that have no metrics simply don't call this wrapper.
func WrapMigrationAdapterForMetrics(
	name string,
	delegate contracts.MigrationAdapter,
	metrics metricscontracts.Metrics,
) contracts.MigrationAdapter {
	validation.AssertNotEmpty("name", name)
	validation.AssertNotNil("delegate", delegate)
	validation.AssertNotNil("metrics", metrics)

	return &migrationAdapterMetrics{
		MigrationAdapter: delegate,
		name:             name,
		metrics:          metrics,
	}
}

func (m *migrationAdapterMetrics) record(ctx context.Context, op string, fn func(context.Context) error) error {
	start := time.Now()
	err := fn(ctx)
	m.metrics.RecordDuration("migration.operation", time.Since(start), "client", m.name, "op", op, "error", err != nil)
	return err
}

func (m *migrationAdapterMetrics) Up(ctx context.Context, version uint) error {
	return m.record(ctx, "up", func(ctx context.Context) error {
		return m.MigrationAdapter.Up(ctx, version)
	})
}

func (m *migrationAdapterMetrics) Down(ctx context.Context, version uint) error {
	return m.record(ctx, "down", func(ctx context.Context) error {
		return m.MigrationAdapter.Down(ctx, version)
	})
}

func (m *migrationAdapterMetrics) Reset(ctx context.Context) error {
	return m.record(ctx, "reset", func(ctx context.Context) error {
		return m.MigrationAdapter.Reset(ctx)
	})
}

func (m *migrationAdapterMetrics) Status(ctx context.Context) error {
	return m.record(ctx, "status", func(ctx context.Context) error {
		return m.MigrationAdapter.Status(ctx)
	})
}

func (m *migrationAdapterMetrics) CurrentVersion(ctx context.Context) (int64, error) {
	start := time.Now()
	version, err := m.MigrationAdapter.CurrentVersion(ctx)
	m.metrics.RecordDuration("migration.operation", time.Since(start), "client", m.name, "op", "current_version", "error", err != nil)
	return version, err
}
