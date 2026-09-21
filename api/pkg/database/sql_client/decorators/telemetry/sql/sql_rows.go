package sql

import (
	"database/sql"
	"frisboo-bank/openapi-generator-service/pkg/database/sql_client/contracts"
	tracercontracts "frisboo-bank/openapi-generator-service/pkg/telemetry/tracer/contracts"
	"sync"
)

var _ contracts.SQLRows = (*sqlRowsTelemetry)(nil)

type sqlRowsTelemetry struct {
	delegate contracts.SQLRows
	span     tracercontracts.TracerSpan
	once     sync.Once
}

func wrapSQLRowsForTelemetry(delegate contracts.SQLRows, span tracercontracts.TracerSpan) contracts.SQLRows {
	return &sqlRowsTelemetry{
		delegate: delegate,
		span:     span,
	}
}

func (s *sqlRowsTelemetry) Close() error {
	err := s.delegate.Close()
	s.end(err)
	return err
}

func (s *sqlRowsTelemetry) ColumnTypes() ([]*sql.ColumnType, error) {
	return s.delegate.ColumnTypes()
}

func (s *sqlRowsTelemetry) Columns() ([]string, error) {
	return s.delegate.Columns()
}

func (s *sqlRowsTelemetry) Err() error {
	return s.delegate.Err()
}

func (s *sqlRowsTelemetry) Next() bool {
	return s.delegate.Next()
}

func (s *sqlRowsTelemetry) NextResultSet() bool {
	return s.delegate.NextResultSet()
}

func (s *sqlRowsTelemetry) Scan(dest ...any) error {
	return s.delegate.Scan(dest...)
}

func (s *sqlRowsTelemetry) end(err error) {
	s.once.Do(func() {
		if err != nil {
			s.span.RecordError(err)
		}
		s.span.End()
	})
}
