package telemetry

import (
	"context"
	"time"

	"frisboo-bank/openapi-generator-service/pkg/database/migration/contracts"
	metricspkg "frisboo-bank/openapi-generator-service/pkg/telemetry/metrics"
	metricscontracts "frisboo-bank/openapi-generator-service/pkg/telemetry/metrics/contracts"
	tracerpkg "frisboo-bank/openapi-generator-service/pkg/telemetry/tracer"
	tracercontracts "frisboo-bank/openapi-generator-service/pkg/telemetry/tracer/contracts"
	"frisboo-bank/openapi-generator-service/pkg/validation"
)

// migrationAdapterTelemetry decorates a MigrationAdapter with tracing and
// metrics, mirroring the sql_client decorator pattern. Name/Type/Logger are
// promoted from the embedded adapter; the mutating and read operations are
// instrumented.
var _ contracts.MigrationAdapter = (*migrationAdapterTelemetry)(nil)

type migrationAdapterTelemetry struct {
	contracts.MigrationAdapter
	name    string
	tracer  tracercontracts.Tracer
	metrics metricscontracts.Metrics
}

// WrapMigrationAdapterForTelemetry returns a MigrationAdapter that emits a
// span and a duration metric for every migration operation. A nil tracer or
// metrics is replaced with a noop adapter, so callers can pass nil without
// opting out at the call site.
func WrapMigrationAdapterForTelemetry(
	name string,
	delegate contracts.MigrationAdapter,
	tracer tracercontracts.Tracer,
	metrics metricscontracts.Metrics,
) contracts.MigrationAdapter {
	validation.AssertNotEmpty("name", name)
	validation.AssertNotNil("delegate", delegate)

	logger := delegate.Logger()
	if tracer == nil {
		tracer, _ = tracerpkg.CreateNoopTracer(name, logger)
	}
	if metrics == nil {
		metrics, _ = metricspkg.CreateNoopMetrics(name, logger)
	}

	return &migrationAdapterTelemetry{
		MigrationAdapter: delegate,
		name:             name,
		tracer:           tracer,
		metrics:          metrics,
	}
}

func (m *migrationAdapterTelemetry) record(ctx context.Context, op string, fn func(context.Context) error) error {
	start := time.Now()
	ctx, span := m.tracer.Start(ctx, "migration."+op)
	defer span.End()

	err := fn(ctx)
	if err != nil {
		span.RecordError(err)
	}
	m.metrics.RecordDuration("migration.operation", time.Since(start), "client", m.name, "op", op, "error", err != nil)
	return err
}

func (m *migrationAdapterTelemetry) Up(ctx context.Context, version uint) error {
	return m.record(ctx, "up", func(ctx context.Context) error {
		return m.MigrationAdapter.Up(ctx, version)
	})
}

func (m *migrationAdapterTelemetry) Down(ctx context.Context, version uint) error {
	return m.record(ctx, "down", func(ctx context.Context) error {
		return m.MigrationAdapter.Down(ctx, version)
	})
}

func (m *migrationAdapterTelemetry) Reset(ctx context.Context) error {
	return m.record(ctx, "reset", func(ctx context.Context) error {
		return m.MigrationAdapter.Reset(ctx)
	})
}

func (m *migrationAdapterTelemetry) Status(ctx context.Context) error {
	return m.record(ctx, "status", func(ctx context.Context) error {
		return m.MigrationAdapter.Status(ctx)
	})
}

func (m *migrationAdapterTelemetry) CurrentVersion(ctx context.Context) (int64, error) {
	start := time.Now()
	ctx, span := m.tracer.Start(ctx, "migration.current_version")
	defer span.End()

	version, err := m.MigrationAdapter.CurrentVersion(ctx)
	if err != nil {
		span.RecordError(err)
	}
	m.metrics.RecordDuration("migration.operation", time.Since(start), "client", m.name, "op", "current_version", "error", err != nil)
	return version, err
}
