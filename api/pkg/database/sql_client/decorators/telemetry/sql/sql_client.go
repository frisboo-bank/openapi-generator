package sql

import (
	"context"
	"database/sql"
	"frisboo-bank/openapi-generator-service/pkg/database/sql_client/contracts"
	sqlclienttype "frisboo-bank/openapi-generator-service/pkg/database/sql_client/models/enums/sql_client_type"
	loggercontracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"
	metricscontracts "frisboo-bank/openapi-generator-service/pkg/telemetry/metrics/contracts"
	"time"

	tracercontracts "frisboo-bank/openapi-generator-service/pkg/telemetry/tracer/contracts"
)

var _ contracts.SQLClient = (*sqlClientTelemetry)(nil)

type sqlClientTelemetry struct {
	delegate contracts.SQLClient
	name     string
	tracer   tracercontracts.Tracer
	metrics  metricscontracts.Metrics
}

func WrapSQLClientForTelemetry(name string, delegate contracts.SQLClient, tracer tracercontracts.Tracer, metrics metricscontracts.Metrics) contracts.SQLClient {
	return &sqlClientTelemetry{
		delegate: delegate,
		metrics:  metrics,
		name:     name,
		tracer:   tracer,
	}
}

func (s *sqlClientTelemetry) BeginTx(ctx context.Context, opts *sql.TxOptions) (contracts.SQLTransaction, error) {
	start := time.Now()
	ctx, span := s.tracer.Start(ctx, "sql.begin_transaction")
	defer span.End()

	tx, err := s.delegate.BeginTx(ctx, opts)
	if err != nil {
		span.RecordError(err)
		s.metrics.RecordDuration("sql.operation", time.Since(start), "client", s.name, "op", "begin_transaction", "error", true)
		return nil, err
	}
	s.metrics.RecordDuration("sql.operation", time.Since(start), "client", s.name, "op", "begin_transaction", "error", false)
	return WrapSQLTransactionForTelemetry(s.name, tx, s.tracer, s.metrics), nil
}

func (s *sqlClientTelemetry) Close(ctx context.Context) error {
	start := time.Now()
	_, span := s.tracer.Start(context.Background(), "sql.close")
	defer span.End()

	err := s.delegate.Close(ctx)
	if err != nil {
		span.RecordError(err)
	}

	s.metrics.RecordDuration("sql.operation", time.Since(start), "client", s.name, "op", "close", "error", err != nil)
	return err
}

func (s *sqlClientTelemetry) Exec(ctx context.Context, query string, args ...any) (sql.Result, error) {
	start := time.Now()
	ctx, span := s.tracer.Start(ctx, "sql.exec")
	defer span.End()

	res, err := s.delegate.Exec(ctx, query, args...)
	if err != nil {
		span.RecordError(err)
	}

	s.metrics.RecordDuration("sql.operation", time.Since(start), "client", s.name, "op", "exec", "error", err != nil)
	return res, err
}

func (s *sqlClientTelemetry) Ping(ctx context.Context) error {
	start := time.Now()
	ctx, span := s.tracer.Start(ctx, "sql.ping")
	defer span.End()

	err := s.delegate.Ping(ctx)
	if err != nil {
		span.RecordError(err)
	}

	s.metrics.RecordDuration("sql.operation", time.Since(start), "client", s.name, "op", "ping", "error", err != nil)
	return err
}

func (s *sqlClientTelemetry) Query(ctx context.Context, query string, args ...any) (contracts.SQLRows, error) {
	start := time.Now()
	ctx, span := s.tracer.Start(ctx, "sql.query")

	rows, err := s.delegate.Query(ctx, query, args...)
	if err != nil {
		span.RecordError(err)
		s.metrics.RecordDuration("sql.operation", time.Since(start), "client", s.name, "op", "query", "error", true)
		span.End()
		return nil, err
	}

	s.metrics.RecordDuration("sql.operation", time.Since(start), "client", s.name, "op", "query", "error", false)
	return wrapSQLRowsForTelemetry(rows, span), nil
}

func (s *sqlClientTelemetry) QueryRow(ctx context.Context, query string, args ...any) contracts.SQLRow {
	ctx, span := s.tracer.Start(ctx, "sql.query_row")

	row := s.delegate.QueryRow(ctx, query, args...)
	return wrapSQLRowForTelemetry(row, span)
}

func (s *sqlClientTelemetry) Logger() loggercontracts.Logger    { return s.delegate.Logger() }
func (s *sqlClientTelemetry) Name() string                      { return s.delegate.Name() }
func (s *sqlClientTelemetry) Type() sqlclienttype.SqlClientType { return s.delegate.Type() }
