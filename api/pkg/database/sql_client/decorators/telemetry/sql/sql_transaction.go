package sql

import (
	"context"
	"database/sql"
	"frisboo-bank/openapi-generator-service/pkg/database/sql_client/contracts"
	metricscontracts "frisboo-bank/openapi-generator-service/pkg/telemetry/metrics/contracts"
	tracercontracts "frisboo-bank/openapi-generator-service/pkg/telemetry/tracer/contracts"
	"time"
)

var _ contracts.SQLTransaction = (*sqlTransactionTelemetry)(nil)

type sqlTransactionTelemetry struct {
	name     string
	delegate contracts.SQLTransaction
	tracer   tracercontracts.Tracer
	metrics  metricscontracts.Metrics
}

func WrapSQLTransactionForTelemetry(name string, delegate contracts.SQLTransaction, tracer tracercontracts.Tracer, metrics metricscontracts.Metrics) contracts.SQLTransaction {
	return &sqlTransactionTelemetry{
		name:     name,
		delegate: delegate,
		metrics:  metrics,
		tracer:   tracer,
	}
}

func (s *sqlTransactionTelemetry) Commit(ctx context.Context) error {
	start := time.Now()
	ctx, span := s.tracer.Start(ctx, "sql.transaction.commit")
	defer span.End()

	err := s.delegate.Commit(ctx)
	if err != nil {
		span.RecordError(err)
	}

	s.metrics.RecordDuration("sql.operation", time.Since(start), "client", s.name, "op", "tx.commit", "error", err != nil)
	return err
}

func (s *sqlTransactionTelemetry) Exec(ctx context.Context, query string, args ...any) (sql.Result, error) {
	start := time.Now()
	ctx, span := s.tracer.Start(ctx, "sql.transaction.exec")
	defer span.End()

	res, err := s.delegate.Exec(ctx, query, args...)
	if err != nil {
		span.RecordError(err)
	}

	s.metrics.RecordDuration("sql.operation", time.Since(start), "client", s.name, "op", "tx.exec", "error", err != nil)
	return res, err
}

func (s *sqlTransactionTelemetry) Query(ctx context.Context, query string, args ...any) (contracts.SQLRows, error) {
	start := time.Now()
	ctx, span := s.tracer.Start(ctx, "sql.transaction.query")

	rows, err := s.delegate.Query(ctx, query, args...)
	if err != nil {
		span.RecordError(err)
		s.metrics.RecordDuration("sql.operation", time.Since(start), "client", s.name, "op", "tx.query", "error", true)
		span.End()
		return nil, err
	}

	s.metrics.RecordDuration("sql.operation", time.Since(start), "client", s.name, "op", "tx.query", "error", false)
	return wrapSQLRowsForTelemetry(rows, span), nil
}

func (s *sqlTransactionTelemetry) QueryRow(ctx context.Context, query string, args ...any) contracts.SQLRow {
	ctx, span := s.tracer.Start(ctx, "sql.transaction.query_row")

	row := s.delegate.QueryRow(ctx, query, args...)
	return wrapSQLRowForTelemetry(row, span)
}

func (s *sqlTransactionTelemetry) Rollback(ctx context.Context) error {
	start := time.Now()
	ctx, span := s.tracer.Start(ctx, "sql.transaction.rollback")

	err := s.delegate.Rollback(ctx)
	if err != nil {
		span.RecordError(err)
	}

	s.metrics.RecordDuration("sql.operation", time.Since(start), "client", s.name, "op", "tx.rollback", "error", err != nil)
	return err
}
