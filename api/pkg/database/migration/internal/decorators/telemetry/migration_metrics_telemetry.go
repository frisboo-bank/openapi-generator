package telemetry

import (
	"context"
	"time"

	"frisboo-bank/openapi-generator-service/pkg/database/migration/contracts"
	"frisboo-bank/openapi-generator-service/pkg/database/migration/types/migrationtype"
	loggercontracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"
	metricscontracts "frisboo-bank/openapi-generator-service/pkg/telemetry/metrics/contracts"
	"frisboo-bank/openapi-generator-service/pkg/validation"
)

type migrationMetrics struct {
	delegate contracts.Migration
	name     string
	metrics  metricscontracts.Metrics
}

func decorateMigrationForMetrics(
	name string,
	delegate contracts.Migration,
	metrics metricscontracts.Metrics,
) contracts.Migration {
	validation.AssertNotEmpty("name", name)
	validation.AssertNotNil("delegate", delegate)
	validation.AssertNotNil("metrics", metrics)

	return &migrationMetrics{
		delegate: delegate,
		name:     name,
		metrics:  metrics,
	}
}

func (m *migrationMetrics) Up(ctx context.Context, version uint) error {
	return m.record(ctx, "up", func(ctx context.Context) error {
		return m.delegate.Up(ctx, version)
	})
}

func (m *migrationMetrics) Down(ctx context.Context, version uint) error {
	return m.record(ctx, "down", func(ctx context.Context) error {
		return m.delegate.Down(ctx, version)
	})
}

func (m *migrationMetrics) Reset(ctx context.Context) error {
	return m.record(ctx, "reset", func(ctx context.Context) error {
		return m.delegate.Reset(ctx)
	})
}

func (m *migrationMetrics) Status(ctx context.Context) error {
	return m.record(ctx, "status", func(ctx context.Context) error {
		return m.delegate.Status(ctx)
	})
}

func (m *migrationMetrics) CurrentVersion(ctx context.Context) (int64, error) {
	var version int64

	err := m.record(ctx, "current_version", func(ctx context.Context) error {
		var err error
		version, err = m.delegate.CurrentVersion(ctx)
		return err
	})

	return version, err
}

func (m *migrationMetrics) record(ctx context.Context, op string, fn func(context.Context) error) (err error) {
	start := time.Now()

	defer func() {
		status := "success"
		if err != nil {
			status = "error"
		}

		if r := recover(); r != nil {
			status = "panic"
			m.metrics.RecordDuration("migration.operation", time.Since(start),
				"name", m.name, "op", op, "status", status)
			panic(r)
		}

		m.metrics.RecordDuration("migration.operation", time.Since(start),
			"name", m.name, "op", op, "status", status)
	}()

	return fn(ctx)
}

func (m *migrationMetrics) Logger() loggercontracts.Logger {
	return m.delegate.Logger()
}

// Name implements [contracts.Migration].
func (m *migrationMetrics) Name() string {
	return m.delegate.Name()
}

func (m *migrationMetrics) Type() migrationtype.MigrationType {
	return m.delegate.Type()
}
