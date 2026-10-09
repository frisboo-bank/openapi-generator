package telemetry

import (
	"context"
	"time"

	"frisboo-bank/openapi-generator-service/pkg/database/migration/contracts"
	metricscontracts "frisboo-bank/openapi-generator-service/pkg/telemetry/metrics/contracts"
	"frisboo-bank/openapi-generator-service/pkg/validation"
)

type migrationMetrics struct {
	contracts.Migration
	name    string
	metrics metricscontracts.Metrics
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
		Migration: delegate,
		name:      name,
		metrics:   metrics,
	}
}

func (m *migrationMetrics) Up(ctx context.Context, version uint) error {
	return m.record(ctx, "up", func(ctx context.Context) error {
		return m.Migration.Up(ctx, version)
	})
}

func (m *migrationMetrics) Down(ctx context.Context, version uint) error {
	return m.record(ctx, "down", func(ctx context.Context) error {
		return m.Migration.Down(ctx, version)
	})
}

func (m *migrationMetrics) Reset(ctx context.Context) error {
	return m.record(ctx, "reset", func(context.Context) error {
		return m.Migration.Reset(ctx)
	})
}

func (m *migrationMetrics) Status(ctx context.Context) error {
	return m.record(ctx, "status", func(context.Context) error {
		return m.Migration.Status(ctx)
	})
}

func (m *migrationMetrics) CurrentVersion(ctx context.Context) (int64, error) {
	var version int64

	err := m.record(ctx, "current_version", func(ctx context.Context) error {
		var err error
		version, err = m.Migration.CurrentVersion(ctx)
		return err
	})

	return version, err
}

func (m *migrationMetrics) record(ctx context.Context, op string, fn func(context.Context) error) error {
	start := time.Now()

	err := fn(ctx)

	m.metrics.RecordDuration(
		"migration.operation",
		time.Since(start),
		"client", m.name,
		"op", op,
		"error", err != nil,
	)

	return err
}
