package telemetry

import (
	"context"
	"fmt"
	"runtime/debug"

	"frisboo-bank/openapi-generator-service/pkg/database/migration/contracts"
	"frisboo-bank/openapi-generator-service/pkg/database/migration/types/migrationtype"
	loggercontracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"
	tracercontracts "frisboo-bank/openapi-generator-service/pkg/telemetry/tracer/contracts"
	"frisboo-bank/openapi-generator-service/pkg/validation"
)

type migrationTracer struct {
	delegate contracts.Migration
	name     string
	tracer   tracercontracts.Tracer
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
		delegate: delegate,
		name:     name,
		tracer:   tracer,
	}
}

func (m *migrationTracer) Up(ctx context.Context, version uint) error {
	return m.record(ctx, "up", func(ctx context.Context, _ tracercontracts.TracerSpan) error {
		return m.delegate.Up(ctx, version)
	})
}

func (m *migrationTracer) Down(ctx context.Context, version uint) error {
	return m.record(ctx, "down", func(ctx context.Context, _ tracercontracts.TracerSpan) error {
		return m.delegate.Down(ctx, version)
	})
}

func (m *migrationTracer) Reset(ctx context.Context) error {
	return m.record(ctx, "reset", func(ctx context.Context, _ tracercontracts.TracerSpan) error {
		return m.delegate.Reset(ctx)
	})
}

func (m *migrationTracer) Status(ctx context.Context) error {
	return m.record(ctx, "status", func(ctx context.Context, _ tracercontracts.TracerSpan) error {
		return m.delegate.Status(ctx)
	})
}

func (m *migrationTracer) CurrentVersion(ctx context.Context) (int64, error) {
	var version int64

	err := m.record(ctx, "current_version", func(ctx context.Context, _ tracercontracts.TracerSpan) error {
		var err error
		version, err = m.delegate.CurrentVersion(ctx)
		return err
	})

	return version, err
}

func (m *migrationTracer) record(ctx context.Context, op string, fn func(context.Context, tracercontracts.TracerSpan) error) (err error) {
	ctx, span := m.tracer.Start(ctx, "migration."+m.name+"."+op)
	defer span.End()

	defer func() {
		if err != nil {
			span.RecordError(err)
		}

		if r := recover(); r != nil {
			span.RecordPanic(fmt.Errorf("panic: %v\n%s", r, debug.Stack()))
			panic(r)
		}
	}()

	return fn(ctx, span)
}

func (m *migrationTracer) Logger() loggercontracts.Logger {
	return m.delegate.Logger()
}

func (m *migrationTracer) Name() string {
	return m.delegate.Name()
}

func (m *migrationTracer) Type() migrationtype.MigrationType {
	return m.delegate.Type()
}
