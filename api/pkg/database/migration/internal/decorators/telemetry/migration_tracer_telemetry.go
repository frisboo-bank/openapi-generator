package telemetry

import (
	"context"

	"frisboo-bank/openapi-generator-service/pkg/database/migration/contracts"
	tracercontracts "frisboo-bank/openapi-generator-service/pkg/telemetry/tracer/contracts"
	"frisboo-bank/openapi-generator-service/pkg/validation"
)

type migrationTracer struct {
	contracts.Migration
	name   string
	tracer tracercontracts.Tracer
}

func decorateMigrationForTracing(
	name string,
	delegate contracts.Migration,
	tracer tracercontracts.Tracer,
) contracts.Migration {
	validation.AssertNotEmpty("name", name)
	validation.AssertNotNil("delegate", delegate)
	validation.AssertNotNil("tracer", tracer)

	return &migrationTracer{
		Migration: delegate,
		name:      name,
		tracer:    tracer,
	}
}

func (m *migrationTracer) Up(ctx context.Context, version uint) error {
	return m.record(ctx, "up", func(ctx context.Context) error {
		return m.Migration.Up(ctx, version)
	})
}

func (m *migrationTracer) Down(ctx context.Context, version uint) error {
	return m.record(ctx, "down", func(ctx context.Context) error {
		return m.Migration.Down(ctx, version)
	})
}

func (m *migrationTracer) Reset(ctx context.Context) error {
	return m.record(ctx, "reset", func(ctx context.Context) error {
		return m.Migration.Reset(ctx)
	})
}

func (m *migrationTracer) Status(ctx context.Context) error {
	return m.record(ctx, "status", func(ctx context.Context) error {
		return m.Migration.Status(ctx)
	})
}

func (m *migrationTracer) CurrentVersion(ctx context.Context) (int64, error) {
	var version int64

	err := m.record(ctx, "current_version", func(ctx context.Context) error {
		var err error
		version, err = m.Migration.CurrentVersion(ctx)
		return err
	})

	return version, err
}

func (m *migrationTracer) record(ctx context.Context, op string, fn func(context.Context) error) error {
	ctx, span := m.tracer.Start(ctx, "migration."+m.name+"."+op)
	defer span.End()

	err := fn(ctx)
	if err != nil {
		span.RecordError(err)
	}
	return err
}
