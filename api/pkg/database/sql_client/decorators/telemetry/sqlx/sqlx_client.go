package sqlx

import (
	"context"
	"database/sql"
	"time"

	"frisboo-bank/openapi-generator-service/pkg/database/sql_client/contracts"
	metricscontracts "frisboo-bank/openapi-generator-service/pkg/telemetry/metrics/contracts"
	"frisboo-bank/openapi-generator-service/pkg/validation"

	tracercontracts "frisboo-bank/openapi-generator-service/pkg/telemetry/tracer/contracts"
)

var _ contracts.SQLXClient = (*sqlxClientTelemetry)(nil)

type sqlxClientTelemetry struct {
	contracts.SQLXClient
	name    string
	tracer  tracercontracts.Tracer
	metrics metricscontracts.Metrics
}

func WrapSQLXClientForTelemetry(name string, delegate contracts.SQLXClient, tracer tracercontracts.Tracer, metrics metricscontracts.Metrics) contracts.SQLXClient {
	validation.AssertNotNil("delegate", delegate)
	validation.AssertNotNil("tracer", tracer)
	validation.AssertNotNil("metrics", metrics)

	return &sqlxClientTelemetry{
		SQLXClient: delegate,
		name:       name,
		tracer:     tracer,
		metrics:    metrics,
	}
}

func (s *sqlxClientTelemetry) BeginTransaction(ctx context.Context, opts *sql.TxOptions) (contracts.SQLXTransaction, error) {
	start := time.Now()
	ctx, span := s.tracer.Start(ctx, "sqlx.begin_transaction")
	defer span.End()

	tx, err := s.SQLXClient.BeginTransaction(ctx, opts)
	if err != nil {
		span.RecordError(err)
		s.metrics.RecordDuration("sql.operation", time.Since(start), "client", s.name, "op", "begin_transaction", "error", true)
		return nil, err
	}

	s.metrics.RecordDuration("sql.operation", time.Since(start), "client", s.name, "op", "begin_transaction", "error", false)
	return wrapSQLXTransactionForTelemetry(s.name, tx, s.tracer, s.metrics), nil
}

func (s *sqlxClientTelemetry) Close(ctx context.Context) error {
	start := time.Now()
	_, span := s.tracer.Start(context.Background(), "sqlx.close")
	defer span.End()

	err := s.SQLXClient.Close(ctx)
	if err != nil {
		span.RecordError(err)
	}

	s.metrics.RecordDuration("sql.operation", time.Since(start), "client", s.name, "op", "close", "error", err != nil)
	return err
}

func (s *sqlxClientTelemetry) NamedExec(ctx context.Context, query string, args map[string]any) (sql.Result, error) {
	start := time.Now()
	ctx, span := s.tracer.Start(ctx, "sqlx.named_exec")
	defer span.End()

	res, err := s.SQLXClient.NamedExec(ctx, query, args)
	if err != nil {
		span.RecordError(err)
	}

	s.metrics.RecordDuration("sql.operation", time.Since(start), "client", s.name, "op", "named_exec", "error", err != nil)
	return res, err
}

func (s *sqlxClientTelemetry) NamedGet(ctx context.Context, dest any, query string, args map[string]any) error {
	start := time.Now()
	ctx, span := s.tracer.Start(ctx, "sqlx.named_get")
	defer span.End()

	err := s.SQLXClient.NamedGet(ctx, dest, query, args)
	if err != nil {
		span.RecordError(err)
	}

	s.metrics.RecordDuration("sql.operation", time.Since(start), "client", s.name, "op", "named_get", "error", err != nil)
	return err
}

func (s *sqlxClientTelemetry) NamedQuery(ctx context.Context, query string, args map[string]any) (contracts.SQLXRows, error) {
	start := time.Now()
	ctx, span := s.tracer.Start(ctx, "sqlx.named_query")

	res, err := s.SQLXClient.NamedQuery(ctx, query, args)
	if err != nil {
		span.RecordError(err)
		s.metrics.RecordDuration("sql.operation", time.Since(start), "client", s.name, "op", "named_query", "error", true)
		span.End()
		return nil, err
	}

	s.metrics.RecordDuration("sql.operation", time.Since(start), "client", s.name, "op", "named_query", "error", false)
	return wrapSQLXRowsForTelemetry(res, span), nil
}

func (s *sqlxClientTelemetry) NamedSelect(ctx context.Context, dest any, query string, args map[string]any) error {
	validation.AssertNotNil("ctx", ctx)

	start := time.Now()
	ctx, span := s.tracer.Start(ctx, "sqlx.named_select")
	defer span.End()

	err := s.SQLXClient.NamedSelect(ctx, dest, query, args)
	if err != nil {
		span.RecordError(err)
	}

	s.metrics.RecordDuration("sql.operation", time.Since(start), "client", s.name, "op", "named_select", "error", err != nil)
	return err
}

func (s *sqlxClientTelemetry) Ping(ctx context.Context) error {
	start := time.Now()
	ctx, span := s.tracer.Start(ctx, "sqlx.ping")
	defer span.End()

	err := s.SQLXClient.Ping(ctx)
	if err != nil {
		span.RecordError(err)
	}

	s.metrics.RecordDuration("sql.operation", time.Since(start), "client", s.name, "op", "ping", "error", err != nil)
	return err
}
