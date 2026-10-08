package telemetry

import (
	"context"
	"time"

	"frisboo-bank/openapi-generator-service/pkg/database/migration/contracts"
	metricscontracts "frisboo-bank/openapi-generator-service/pkg/telemetry/metrics/contracts"
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
// span and a duration metric for every migration operation.
func WrapMigrationAdapterForTelemetry(
	name string,
	delegate contracts.MigrationAdapter,
	tracer tracercontracts.Tracer,
	metrics metricscontracts.Metrics,
) contracts.MigrationAdapter {
	validation.AssertNotEmpty("name", name)
	validation.AssertNotNil("delegate", delegate)

	return &migrationAdapterTelemetry{
		MigrationAdapter: delegate,
		name:             name,
		tracer:           tracer,
		metrics:          metrics,
	}
}

func (m *migrationAdapterTelemetry) record(ctx context.Context, op string, fn func(context.Context) error) error {
	start := time.Now()

	var span tracercontracts.TracerSpan
	if m.tracer != nil {
		ctx, span = m.tracer.Start(ctx, "migration."+op)
	}

	err := fn(ctx)
	if span != nil {
		defer span.End()
		if err != nil {
			span.RecordError(err)
		}
	}
	if m.metrics != nil {
		m.metrics.RecordDuration("migration.operation", time.Since(start), "client", m.name, "op", op, "error", err != nil)
	}
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

	var span tracercontracts.TracerSpan
	if m.tracer != nil {
		ctx, span = m.tracer.Start(ctx, "migration.current_version")
	}

	version, err := m.MigrationAdapter.CurrentVersion(ctx)
	if span != nil {
		defer span.End()
		if err != nil {
			span.RecordError(err)
		}
	}
	if m.metrics != nil {
		m.metrics.RecordDuration("migration.operation", time.Since(start), "client", m.name, "op", "current_version", "error", err != nil)
	}
	return version, err
}
