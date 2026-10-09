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
// metrics. A nil tracer or metrics is tolerated — the operation simply
// skips that instrument. Each instrument's logic lives in its own method
// (decorateMigrationForTracing / decorateMigrationForMetrics), so adding
// or changing one instrument never touches the other.
var _ contracts.MigrationAdapter = (*migrationAdapterTelemetry)(nil)

type migrationAdapterTelemetry struct {
	contracts.MigrationAdapter
	name    string
	tracer  tracercontracts.Tracer
	metrics metricscontracts.Metrics
}

// WrapMigrationAdapterForTelemetry wraps a MigrationAdapter so that every
// operation emits a tracing span and a duration metric. Either instrument
// may be nil — the operation degrades to the bare adapter for that
// instrument. The start time is captured once per operation and shared
// between the span and the metric, so their durations always agree.
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

// decorateMigrationForTracing wraps fn so that calling it emits a tracing
// span named "migration.<op>". When tracer is nil, fn is returned unchanged.
func (m *migrationAdapterTelemetry) decorateMigrationForTracing(op string, fn func(context.Context) error) func(context.Context) error {
	if m.tracer == nil {
		return fn
	}
	return func(ctx context.Context) error {
		ctx, span := m.tracer.Start(ctx, "migration."+op)
		defer span.End()
		err := fn(ctx)
		if err != nil {
			span.RecordError(err)
		}
		return err
	}
}

// decorateMigrationForMetrics wraps fn so that calling it records a duration
// metric named "migration.operation" tagged with the operation name. The
// start time is supplied by the caller so it can be shared with the tracing
// span — both instruments measure from the same reference point. When
// metrics is nil, fn is returned unchanged.
func (m *migrationAdapterTelemetry) decorateMigrationForMetrics(op string, fn func(context.Context) error, start time.Time) func(context.Context) error {
	if m.metrics == nil {
		return fn
	}
	return func(ctx context.Context) error {
		err := fn(ctx)
		m.metrics.RecordDuration("migration.operation", time.Since(start), "client", m.name, "op", op, "error", err != nil)
		return err
	}
}

func (m *migrationAdapterTelemetry) record(ctx context.Context, op string, fn func(context.Context) error) error {
	start := time.Now()
	// Tracing is the inner wrapper (starts the span, calls fn, ends the span).
	// Metrics is the outer wrapper (captures nothing — it uses the shared
	// start), so both instruments measure from the same reference point.
	decorated := m.decorateMigrationForMetrics(op, m.decorateMigrationForTracing(op, fn), start)
	return decorated(ctx)
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

	// Can't use record() because the return type differs, but the same
	// composition pattern: tracing inner, metrics outer, shared start.
	wrapped := m.MigrationAdapter.CurrentVersion
	if m.tracer != nil {
		prev := wrapped
		wrapped = func(ctx context.Context) (int64, error) {
			ctx, span := m.tracer.Start(ctx, "migration.current_version")
			defer span.End()
			v, err := prev(ctx)
			if err != nil {
				span.RecordError(err)
			}
			return v, err
		}
	}
	if m.metrics != nil {
		prev := wrapped
		wrapped = func(ctx context.Context) (int64, error) {
			v, err := prev(ctx)
			m.metrics.RecordDuration("migration.operation", time.Since(start), "client", m.name, "op", "current_version", "error", err != nil)
			return v, err
		}
	}
	return wrapped(ctx)
}
