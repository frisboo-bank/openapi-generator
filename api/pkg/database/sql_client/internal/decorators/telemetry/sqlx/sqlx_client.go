package sqlx

import (
	"context"
	"database/sql"
	"time"

	loggercontracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"

	"frisboo-bank/openapi-generator-service/pkg/database/sql_client/contracts"
	"frisboo-bank/openapi-generator-service/pkg/database/sql_client/types/sqlclienttype"
	metricscontracts "frisboo-bank/openapi-generator-service/pkg/telemetry/metrics/contracts"
	"frisboo-bank/openapi-generator-service/pkg/validation"

	tracercontracts "frisboo-bank/openapi-generator-service/pkg/telemetry/tracer/contracts"
)

var (
	_ contracts.SQLXClient   = (*sqlxClientTelemetry)(nil)
	_ contracts.WithDBGetter = (*sqlxClientTelemetry)(nil)
)

type sqlxClientTelemetry struct {
	delegate contracts.SQLXClient
	name     string
	tracer   tracercontracts.Tracer
	metrics  metricscontracts.Metrics
}

func WrapSQLXClientForTelemetry(
	name string,
	delegate contracts.SQLXClient,
	tracer tracercontracts.Tracer,
	metrics metricscontracts.Metrics,
) contracts.SQLXClient {
	validation.AssertNotEmpty("name", name)
	validation.AssertNotNil("delegate", delegate)
	validation.AssertNotNil("tracer", tracer)
	validation.AssertNotNil("metrics", metrics)

	return &sqlxClientTelemetry{
		name:     name,
		delegate: delegate,
		tracer:   tracer,
		metrics:  metrics,
	}
}

func (s *sqlxClientTelemetry) BeginTransaction(ctx context.Context, opts *sql.TxOptions) (contracts.SQLXTransaction, error) {
	start := time.Now()
	ctx, span := s.tracer.Start(ctx, "sqlx.begin_transaction")
	defer span.End()

	tx, err := s.delegate.BeginTransaction(ctx, opts)
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
	ctx, span := s.tracer.Start(ctx, "sqlx.close")
	defer span.End()

	err := s.delegate.Close(ctx)
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

	res, err := s.delegate.NamedExec(ctx, query, args)
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

	err := s.delegate.NamedGet(ctx, dest, query, args)
	if err != nil {
		span.RecordError(err)
	}

	s.metrics.RecordDuration("sql.operation", time.Since(start), "client", s.name, "op", "named_get", "error", err != nil)
	return err
}

func (s *sqlxClientTelemetry) NamedQuery(ctx context.Context, query string, args map[string]any) (contracts.SQLXRows, error) {
	start := time.Now()
	ctx, span := s.tracer.Start(ctx, "sqlx.named_query")

	res, err := s.delegate.NamedQuery(ctx, query, args)
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
	start := time.Now()
	ctx, span := s.tracer.Start(ctx, "sqlx.named_select")
	defer span.End()

	err := s.delegate.NamedSelect(ctx, dest, query, args)
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

	err := s.delegate.Ping(ctx)
	if err != nil {
		span.RecordError(err)
	}

	s.metrics.RecordDuration("sql.operation", time.Since(start), "client", s.name, "op", "ping", "error", err != nil)
	return err
}

func (s *sqlxClientTelemetry) DB() *sql.DB {
	if getter, ok := s.delegate.(contracts.WithDBGetter); ok {
		return getter.DB()
	}
	return nil
}

func (s *sqlxClientTelemetry) Logger() loggercontracts.Logger    { return s.delegate.Logger() }
func (s *sqlxClientTelemetry) Name() string                      { return s.delegate.Name() }
func (s *sqlxClientTelemetry) Type() sqlclienttype.SqlClientType { return s.delegate.Type() }
