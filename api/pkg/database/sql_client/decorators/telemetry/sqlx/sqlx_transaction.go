package sqlx

import (
	"context"
	"database/sql"
	"time"

	"frisboo-bank/openapi-generator-service/pkg/database/sql_client/contracts"
	metricscontracts "frisboo-bank/openapi-generator-service/pkg/telemetry/metrics/contracts"
	tracercontracts "frisboo-bank/openapi-generator-service/pkg/telemetry/tracer/contracts"
)

var _ contracts.SQLXTransaction = (*sqlxTransactionTelemetry)(nil)

type sqlxTransactionTelemetry struct {
	contracts.SQLXTransaction
	name    string
	tracer  tracercontracts.Tracer
	metrics metricscontracts.Metrics
}

func wrapSQLXTransactionForTelemetry(name string, delegate contracts.SQLXTransaction, tracer tracercontracts.Tracer, metrics metricscontracts.Metrics) contracts.SQLXTransaction {
	return &sqlxTransactionTelemetry{
		SQLXTransaction: delegate,
		name:            name,
		metrics:         metrics,
		tracer:          tracer,
	}
}

func (s *sqlxTransactionTelemetry) Commit(ctx context.Context) error {
	start := time.Now()
	ctx, span := s.tracer.Start(ctx, "sqlx.transaction.commit")
	defer span.End()

	err := s.SQLXTransaction.Commit(ctx)
	if err != nil {
		span.RecordError(err)
	}

	s.metrics.RecordDuration("sql.operation", time.Since(start), "client", s.name, "op", "tx.commit", "error", err != nil)
	return err
}

func (s *sqlxTransactionTelemetry) NamedExec(ctx context.Context, query string, args map[string]any) (sql.Result, error) {
	start := time.Now()
	ctx, span := s.tracer.Start(ctx, "sqlx.transaction.named_exec")
	defer span.End()

	res, err := s.SQLXTransaction.NamedExec(ctx, query, args)
	if err != nil {
		span.RecordError(err)
	}

	s.metrics.RecordDuration("sql.operation", time.Since(start), "client", s.name, "op", "tx.named_exec", "error", err != nil)
	return res, err
}

func (s *sqlxTransactionTelemetry) NamedGet(ctx context.Context, dest any, query string, args map[string]any) error {
	start := time.Now()
	ctx, span := s.tracer.Start(ctx, "sqlx.transaction.named_get")
	defer span.End()

	err := s.SQLXTransaction.NamedGet(ctx, dest, query, args)
	if err != nil {
		span.RecordError(err)
	}

	s.metrics.RecordDuration("sql.operation", time.Since(start), "client", s.name, "op", "tx.named_get", "error", err != nil)
	return err
}

func (s *sqlxTransactionTelemetry) NamedQuery(ctx context.Context, query string, args map[string]any) (contracts.SQLXRows, error) {
	start := time.Now()
	ctx, span := s.tracer.Start(ctx, "sqlx.transaction.named_query")

	rows, err := s.SQLXTransaction.NamedQuery(ctx, query, args)
	if err != nil {
		span.RecordError(err)
		s.metrics.RecordDuration("sql.operation", time.Since(start), "client", s.name, "op", "tx.named_query", "error", true)
		span.End()
		return nil, err
	}

	s.metrics.RecordDuration("sql.operation", time.Since(start), "client", s.name, "op", "tx.named_query", "error", false)
	return wrapSQLXRowsForTelemetry(rows, span), nil
}

func (s *sqlxTransactionTelemetry) NamedSelect(ctx context.Context, dest any, query string, args map[string]any) error {
	start := time.Now()
	ctx, span := s.tracer.Start(ctx, "sqlx.transaction.named_select")
	defer span.End()

	err := s.SQLXTransaction.NamedSelect(ctx, dest, query, args)
	if err != nil {
		span.RecordError(err)
	}

	s.metrics.RecordDuration("sql.operation", time.Since(start), "client", s.name, "op", "tx.named_select", "error", err != nil)
	return err
}

func (s *sqlxTransactionTelemetry) Rollback(ctx context.Context) error {
	start := time.Now()
	ctx, span := s.tracer.Start(ctx, "sqlx.transaction.rollback")
	defer span.End()

	err := s.SQLXTransaction.Rollback(ctx)
	if err != nil {
		span.RecordError(err)
	}

	s.metrics.RecordDuration("sql.operation", time.Since(start), "client", s.name, "op", "tx.rollback", "error", err != nil)
	return err
}
