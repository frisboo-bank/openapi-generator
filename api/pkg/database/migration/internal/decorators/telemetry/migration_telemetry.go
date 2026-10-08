package telemetry

import (
	"context"
	"time"

	"frisboo-bank/openapi-generator-service/pkg/database/migration/contracts"
	tracercontracts "frisboo-bank/openapi-generator-service/pkg/telemetry/tracer/contracts"

	metricscontracts "frisboo-bank/openapi-generator-service/pkg/telemetry/metrics/contracts"
	"frisboo-bank/openapi-generator-service/pkg/validation"
)

var _ contracts.Migration = (*migrationTelemetry)(nil)

type migrationTelemetry struct {
	contracts.Migration
	name    string
	tracer  tracercontracts.Tracer
	metrics metricscontracts.Metrics
}

func WrapMigrationForTelemetry(
	name string,
	delegate contracts.Migration,
	tracer tracercontracts.Tracer,
	metrics metricscontracts.Metrics,
) contracts.Migration {
	validation.AssertNotEmpty("name", name)
	validation.AssertNotNil("delegate", delegate)
	validation.AssertNotNil("tracer", tracer)
	validation.AssertNotNil("metrics", metrics)

	return &migrationTelemetry{
		Migration: delegate,
		name:      name,
		tracer:    tracer,
		metrics:   metrics,
	}
}

func (m *migrationTelemetry) CurrentVersion(ctx context.Context) (int64, error) {
	start := time.Now()
	ctx, span := m.tracer.Start(ctx, "migration.current_version")
	defer span.End()

	version, err := m.Migration.CurrentVersion(ctx)
	if err != nil {
		span.RecordError(err)
	}

	m.metrics.RecordDuration("migration.operation", time.Since(start), "client", m.name, "op", "current_version", "error", err != nil)
	return version, err
}

func (m *migrationTelemetry) Down(ctx context.Context, version uint) error {
	return m.record(ctx, "down", func(ctx context.Context) error {
		return m.Migration.Down(ctx, version)
	})
}

func (m *migrationTelemetry) Reset(ctx context.Context) error {
	return m.record(ctx, "reset", func(ctx context.Context) error {
		return m.Migration.Reset(ctx)
	})
}

func (m *migrationTelemetry) Status(ctx context.Context) error {
	return m.record(ctx, "status", func(ctx context.Context) error {
		return m.Migration.Status(ctx)
	})
}

func (m *migrationTelemetry) Up(ctx context.Context, version uint) error {
	return m.record(ctx, "up", func(ctx context.Context) error {
		return m.Migration.Up(ctx, version)
	})
}

func (m *migrationTelemetry) record(ctx context.Context, op string, fn func(context.Context) error) error {
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
